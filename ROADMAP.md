# 🗺️ Project Roadmap

This roadmap is impact-driven, avoiding hype to focus on solving real architectural constraints found in cloud-hosted serverless solutions.

## 🟢 Phase 0: Foundation (Current)
*Establish Reactor as a stable, deterministic, production-grade runtime.*

- [x] **In-process WASM execution** via `wazero`.
- [x] **Stdin / Stdout model** for request-response.
- [x] **Zero idle memory** footprint.
- [x]Crash-safe execution** (panic isolation).
- [x]*Explicit execution timeouts** configuration.
- [x]*Memory limits** per execution.
- [x] **Documentation:** Architecture deep-dive & Transparent benchmarks.

## 🟡 Phase 1: Workers Parity
*Allow conceptual migration from Cloudflare Workers without copying edge abstractions.*

- [x] **Structured HTTP Context:** Pass Method, Headers, Path, and Query via structured JSON/Proto to Stdin.
- [x] **Header Mutation:** Allow WASM to modify response headers.
- [x] **Status Code Control:** Allow WASM to set HTTP 4xx/5xx codes.
- [x] **Environment Variables:** Configuration per function.
- [x] **Per-route Binding:** Map specific `.wasm` files to specific Caddy routes.

## 🔵 Phase 2: What Workers Cannot Do
*Surpass Cloud Workers by exploiting local, in-process execution properties.*

- [ ] **CPU Budgeting:** Strict CPU cycle limits per execution.
- [ ] **Hard Memory Ceiling:** Deterministic memory caps per sandbox.
- [ ] **Observability:** Native OpenTelemetry hooks for function execution time.
- [ ] **Multi-Tenancy:** Isolation logic for running 1,000+ distinct functions on one server.

## 🟣 Phase 3: Reusable Compute
*Turn Reactor into a compute primitive.*

- [ ] **Local Function Registry.**
- [ ] **Signed WASM verification.**
- [ ] **Built-in Modules:** JWT Validation, Rate Limiting, Webhook Verification.

## 🔴 Phase 4: Distributed, Without Lock-in
*Scale without becoming a cloud provider.*

- [ ] **Universal Runtime:** Same binary runs on Bare Metal, VPS, and Edge.
- [ ] **Replication:** Optional function replication strategies.

---

### ❌ Explicit Non-Goals
To keep the project focused, we will **NOT** build:
* Proprietary APIs or Vendor-specific abstractions.
* Forced control planes or billing layers.
* "Magic" networking or hidden overlays.