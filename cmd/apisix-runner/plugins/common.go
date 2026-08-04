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
