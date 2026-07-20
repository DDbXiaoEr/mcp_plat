package main

import (
	_ "mcp_plat-console/cmd/apisix-runner/plugins"

	"github.com/apache/apisix-go-plugin-runner/pkg/runner"
)

func main() {
	runner.Run(runner.RunnerConfig{})
}
