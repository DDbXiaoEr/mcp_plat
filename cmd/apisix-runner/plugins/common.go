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

package plugins

import (
	"encoding/json"
	"strings"

	runnerHttp "github.com/apache/apisix-go-plugin-runner/pkg/http"
)

type jsonRpcRequest struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type toolCallParams struct {
	Name string `json:"name"`
}

func extractAccessKey(r runnerHttp.Request, headerName string) string {
	authHeader := r.Header().Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return r.Header().Get(headerName)
}

func extractToolName(r runnerHttp.Request) string {
	if r.Method() != "POST" {
		return ""
	}
	body, err := r.Body()
	if err != nil || len(body) == 0 {
		return ""
	}
	var jrpc jsonRpcRequest
	if json.Unmarshal(body, &jrpc) != nil || jrpc.Method != "tools/call" {
		return ""
	}
	var params toolCallParams
	if json.Unmarshal(jrpc.Params, &params) != nil {
		return ""
	}
	return params.Name
}
