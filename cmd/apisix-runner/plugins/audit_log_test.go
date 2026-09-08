package plugins

import "testing"

func TestResolveResult(t *testing.T) {
	cases := []struct {
		name    string
		status  int
		body    string
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "http ok plain json",
			status:  200,
			body:    `{"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"hi"}]}}`,
			wantOK:  true,
			wantMsg: "ok",
		},
		{
			name:    "http 500",
			status:  500,
			body:    `oops`,
			wantOK:  false,
			wantMsg: "HTTP 500 Internal Server Error",
		},
		{
			name:    "jsonrpc top-level error even with 200",
			status:  200,
			body:    `{"jsonrpc":"2.0","id":1,"error":{"code":-32602,"message":"unknown tool: foo"}}`,
			wantOK:  false,
			wantMsg: "rpc error code=-32602 unknown tool: foo",
		},
		{
			name:    "jsonrpc error without message",
			status:  200,
			body:    `{"jsonrpc":"2.0","id":1,"error":{"code":-32000}}`,
			wantOK:  false,
			wantMsg: "rpc error code=-32000",
		},
		{
			name:    "tool returned isError inside result",
			status:  200,
			body:    `{"jsonrpc":"2.0","id":1,"result":{"content":[{"type":"text","text":"Error: boom"}],"isError":true}}`,
			wantOK:  false,
			wantMsg: "tool returned isError",
		},
		{
			name:    "empty body with 200",
			status:  200,
			body:    ``,
			wantOK:  true,
			wantMsg: "ok",
		},
		{
			name:    "non-json body with 200",
			status:  200,
			body:    `not-json`,
			wantOK:  true,
			wantMsg: "ok",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gotOK, gotMsg := resolveResult(c.status, []byte(c.body))
			if gotOK != c.wantOK || gotMsg != c.wantMsg {
				t.Fatalf("resolveResult(%d, %q) = (%v, %q), want (%v, %q)",
					c.status, c.body, gotOK, gotMsg, c.wantOK, c.wantMsg)
			}
		})
	}
}
