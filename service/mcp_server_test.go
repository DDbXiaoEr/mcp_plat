package service

// Author: deepseek-v4-pro / opencode

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const testSessionID = "mcp-session-test-0001"

func newMockStreamableHTTPServer(t *testing.T, sseResponse bool) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			JSONRPC string `json:"jsonrpc"`
			ID      int    `json:"id"`
			Method  string `json:"method"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		switch req.Method {
		case "initialize":
			w.Header().Set("Mcp-Session-Id", testSessionID)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{"protocolVersion":"2025-03-26","capabilities":{"tools":{}},"serverInfo":{"name":"mock","version":"1.0.0"}}}`, req.ID)
		case "notifications/initialized":
			if r.Header.Get("Mcp-Session-Id") != testSessionID {
				http.Error(w, "Invalid session ID", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusAccepted)
		case "tools/list":
			if r.Header.Get("Mcp-Session-Id") != testSessionID {
				http.Error(w, "Invalid session ID", http.StatusNotFound)
				return
			}
			body := fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":{"tools":[{"name":"tool_a","description":"a"},{"name":"tool_b","description":"b"}]}}`, req.ID)
			if sseResponse {
				w.Header().Set("Content-Type", "text/event-stream")
				fmt.Fprintf(w, "event: message\ndata: %s\n\n", body)
			} else {
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, body)
			}
		default:
			http.Error(w, "unknown method", http.StatusBadRequest)
		}
	}))
}

func TestFetchToolsStreamableHTTP(t *testing.T) {
	srv := newMockStreamableHTTPServer(t, false)
	defer srv.Close()

	tools, err := FetchTools(FetchToolsInput{Address: srv.URL, Protocol: "Streamable HTTP"})
	if err != nil {
		t.Fatalf("FetchTools 失败: %v", err)
	}
	if len(tools) != 2 || tools[0].Name != "tool_a" || tools[1].Name != "tool_b" {
		t.Fatalf("工具列表不符合预期: %v", tools)
	}
}

func TestFetchToolsStreamableHTTPWithSSEResponse(t *testing.T) {
	srv := newMockStreamableHTTPServer(t, true)
	defer srv.Close()

	tools, err := FetchTools(FetchToolsInput{Address: srv.URL, Protocol: "Streamable HTTP"})
	if err != nil {
		t.Fatalf("FetchTools 失败: %v", err)
	}
	if len(tools) != 2 || tools[0].Name != "tool_a" || tools[1].Name != "tool_b" {
		t.Fatalf("工具列表不符合预期: %v", tools)
	}
}

func TestFetchToolsMissingSessionID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     int    `json:"id"`
			Method string `json:"method"`
		}
		json.NewDecoder(r.Body).Decode(&req)
		if req.Method == "initialize" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"jsonrpc":"2.0","id":%d,"result":{}}`, req.ID)
			return
		}
		http.Error(w, "Invalid session ID", http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := FetchTools(FetchToolsInput{Address: srv.URL, Protocol: "Streamable HTTP"})
	if err == nil {
		t.Fatal("预期返回错误，实际为 nil")
	}
	if !strings.Contains(err.Error(), "404") {
		t.Fatalf("预期错误包含 404，实际: %v", err)
	}
}

func TestFetchToolsUnsupportedProtocol(t *testing.T) {
	_, err := FetchTools(FetchToolsInput{Address: "http://example.com", Protocol: "stdio"})
	if err == nil {
		t.Fatal("预期 stdio 协议返回错误，实际为 nil")
	}
}

func TestFetchToolsLive(t *testing.T) {
	address := os.Getenv("MCP_TEST_ADDRESS")
	if address == "" {
		t.Skip("未设置 MCP_TEST_ADDRESS，跳过真实服务器测试")
	}

	tools, err := FetchTools(FetchToolsInput{Address: address, Protocol: "Streamable HTTP"})
	if err != nil {
		t.Fatalf("FetchTools 失败: %v", err)
	}
	t.Logf("获取到工具列表: %v", tools)
}
