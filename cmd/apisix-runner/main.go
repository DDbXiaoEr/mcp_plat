// Copyright (C) 2026 Zhaoquan Wang
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

// Author: deepseek-v4-pro / opencode

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
	fmt.Printf("registered plugins: [accesskey_verify, audit_log]\n")
	fmt.Println("================================")

	runner.Run(runner.RunnerConfig{})
}
