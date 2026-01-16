package reactor

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"

	"github.com/dustin/go-humanize"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

func init() {
	caddy.RegisterModule(Reactor{})
	httpcaddyfile.RegisterHandlerDirective("reactor", parseCaddyfile)
}

type Reactor struct {
	Path        string            `json:"path,omitempty"`
	Args        []string          `json:"args,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Timeout     caddy.Duration    `json:"timeout,omitempty"`
	MemoryLimit string            `json:"memory_limit,omitempty"`

	logger *zap.Logger
	code   wazero.CompiledModule
	engine wazero.Runtime
}

func (Reactor) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.reactor",
		New: func() caddy.Module { return &Reactor{} },
	}
}

func (r *Reactor) Provision(ctx caddy.Context) error {
	r.logger = ctx.Logger()

	if r.Path == "" {
		return fmt.Errorf("wasm file path is required")
	}

	// 1. Configuração do Runtime (Engine)
	ctxWazero := context.Background()
	rConfig := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)

	// ✅ CORREÇÃO: Aplicamos o limite de memória AQUI, na criação do Runtime
	if r.MemoryLimit != "" {
		bytes, err := humanize.ParseBytes(r.MemoryLimit)
		if err != nil {
			return fmt.Errorf("invalid memory_limit: %v", err)
		}

		if bytes > 0 {
			const wasmPageSize = 65536
			pages := uint32(bytes / wasmPageSize)
			if bytes%wasmPageSize != 0 {
				pages++
			}
			rConfig = rConfig.WithMemoryLimitPages(pages)
			r.logger.Info("memory limit configured", zap.String("limit", r.MemoryLimit), zap.Uint32("pages", pages))
		}
	}

	// Cria o runtime com as configurações (Timeout kill + Memory limit)
	r.engine = wazero.NewRuntimeWithConfig(ctxWazero, rConfig)

	// 2. Setup do WASI
	wasi_snapshot_preview1.MustInstantiate(ctxWazero, r.engine)

	// 3. Compilação do Módulo
	wasmBytes, err := os.ReadFile(r.Path)
	if err != nil {
		return fmt.Errorf("failed to read wasm file: %w", err)
	}

	r.code, err = r.engine.CompileModule(ctxWazero, wasmBytes)
	if err != nil {
		return fmt.Errorf("failed to compile wasm binary: %w", err)
	}

	// Default Timeout
	if r.Timeout == 0 {
		r.Timeout = caddy.Duration(60 * time.Second)
	}

	return nil
}

func (r *Reactor) Cleanup() error {
	if r.engine != nil {
		return r.engine.Close(context.Background())
	}
	return nil
}

func (r *Reactor) ServeHTTP(rw http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	ctx, cancel := context.WithTimeout(req.Context(), time.Duration(r.Timeout))
	defer cancel()

	// Configuração da Execução (Instância)
	config := wazero.NewModuleConfig().
		WithStdout(rw).
		WithStderr(os.Stderr).
		WithStdin(req.Body).
		WithArgs(r.Args...)

	for k, v := range r.Env {
		config = config.WithEnv(k, v)
	}

	// Instancia e executa (O limite de memória já está imposto pelo r.engine)
	instance, err := r.engine.InstantiateModule(ctx, r.code, config)

	if err != nil {
		// Timeout
		if ctx.Err() == context.DeadlineExceeded {
			r.logger.Error("wasm execution timed out",
				zap.Duration("limit", time.Duration(r.Timeout)),
			)
			return caddyhttp.Error(http.StatusGatewayTimeout, fmt.Errorf("execution time limit exceeded"))
		}

		// Out of Memory ou Panic
		r.logger.Error("wasm execution failed", zap.Error(err))
		return caddyhttp.Error(http.StatusInternalServerError, err)
	}

	defer instance.Close(ctx)
	return nil
}

func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var m Reactor
	m.Env = make(map[string]string)

	for h.Next() {
		args := h.RemainingArgs()
		if len(args) > 0 {
			m.Path = args[0]
		}

		for h.NextBlock(0) {
			switch h.Val() {
			case "env":
				if h.NextArg() {
					key := h.Val()
					if h.NextArg() {
						m.Env[key] = h.Val()
					}
				}
			case "args":
				m.Args = h.RemainingArgs()
			case "timeout":
				if h.NextArg() {
					val, err := caddy.ParseDuration(h.Val())
					if err != nil {
						return nil, h.Errf("invalid duration: %v", err)
					}
					m.Timeout = caddy.Duration(val)
				}
			case "memory_limit":
				if h.NextArg() {
					m.MemoryLimit = h.Val()
				}
			}
		}
	}
	return &m, nil
}
