package main

import (
	"fmt"
	"os"

	_ "mcp_plat-console/cmd/apisix-runner/plugins"

	"github.com/apache/apisix-go-plugin-runner/pkg/runner"
)

func main() {
	listenAddr := os.Getenv("APISIX_LISTEN_ADDRESS")
	confTTL := os.Getenv("APISIX_CONF_EXPIRE_TIME")

	fmt.Println("=== apisix-go-runner startup ===")
	fmt.Printf("APISIX_LISTEN_ADDRESS: %s\n", listenAddr)
	fmt.Printf("APISIX_CONF_EXPIRE_TIME: %s\n", confTTL)
	fmt.Printf("registered plugins: [accesskey_verify]\n")
	fmt.Println("================================")

	runner.Run(runner.RunnerConfig{})
}
