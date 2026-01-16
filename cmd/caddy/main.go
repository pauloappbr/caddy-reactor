package main

import (
	caddycmd "github.com/caddyserver/caddy/v2/cmd"

	_ "github.com/caddyserver/caddy/v2/modules/standard"

	_ "github.com/pauloappbr/caddy-wazero-reactor"
)

func main() {
	caddycmd.Main()
}
