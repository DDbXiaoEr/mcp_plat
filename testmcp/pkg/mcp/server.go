package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

const ProtocolVersion = "2026-07-28"

var defaultProtocolVersions = []string{ProtocolVersion}

type ServerOption func(*Server)

func WithProtocolVersions(versions []string) ServerOption {
	return func(s *Server) {
		if len(versions) > 0 {
			s.protocolVersions = append([]string(nil), versions...)
		}
	}
}

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      json.RawMessage `json:"id,omitempty"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

type InputSchema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

type Property struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type CallToolResult struct {
	Content []Content `json:"content"`
	IsError bool      `json:"isError,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type InitializeParams struct {
	ProtocolVersion string          `json:"protocolVersion"`
	Capabilities    json.RawMessage `json:"capabilities,omitempty"`
	ClientInfo      json.RawMessage `json:"clientInfo,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
	Instructions    string       `json:"instructions,omitempty"`
}

type Capabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

type ToolHandler func(ctx context.Context, args map[string]interface{}) (*CallToolResult, error)

type Server struct {
	name             string
	version          string
	protocolVersions []string
	tools            []Tool
	handlers         map[string]ToolHandler
	mu               sync.RWMutex
}

func NewServer(name, version string, opts ...ServerOption) *Server {
	s := &Server{
		name:             name,
		version:          version,
		protocolVersions: append([]string(nil), defaultProtocolVersions...),
		handlers:         make(map[string]ToolHandler),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *Server) Name() string    { return s.name }
func (s *Server) Version() string { return s.version }

func (s *Server) SetProtocolVersions(versions ...string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(versions) > 0 {
		s.protocolVersions = append([]string(nil), versions...)
	}
}

func (s *Server) ProtocolVersions() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]string(nil), s.protocolVersions...)
}

func (s *Server) DefaultProtocolVersion() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.protocolVersions) == 0 {
		return ProtocolVersion
	}
	return s.protocolVersions[len(s.protocolVersions)-1]
}

func (s *Server) supportsProtocolVersion(v string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, pv := range s.protocolVersions {
		if pv == v {
			return true
		}
	}
	return false
}

func (s *Server) RegisterTool(tool Tool, handler ToolHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tools = append(s.tools, tool)
	s.handlers[tool.Name] = handler
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/mcp", s.serve)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok", "service": s.name})
	})
	return mux
}

func (s *Server) ListenAndServe(addr string) error {
	log.Printf("[%s] MCP Server (protocol %s) starting on %s", s.name, s.DefaultProtocolVersion(), addr)
	return http.ListenAndServe(addr, s.Handler())
}

func (s *Server) ListenAndServeMulti(addrs []string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, len(addrs))

	for _, addr := range addrs {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			log.Printf("[%s] MCP Server (protocol %s) starting on %s", s.name, s.DefaultProtocolVersion(), addr)
			if err := http.ListenAndServe(addr, s.Handler()); err != nil {
				errCh <- err
			}
		}(addr)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		return err
	}
	return nil
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, MCP-Protocol-Version, Mcp-Method, Mcp-Name")
	w.Header().Set("Access-Control-Expose-Headers", "X-Accel-Buffering")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{
			"jsonrpc": "2.0",
			"error":   `{"code": -32601, "message": "Method Not Allowed"}`,
		})
		return
	}

	if r.Method != http.MethodPost {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "only POST is allowed"})
		return
	}

	s.handlePost(w, r)
}

func (s *Server) handlePost(w http.ResponseWriter, r *http.Request) {
	protoVersion := r.Header.Get("MCP-Protocol-Version")
	if protoVersion != "" && !s.supportsProtocolVersion(protoVersion) {
		s.writeHTTPError(w, http.StatusBadRequest, &JSONRPCError{
			Code:    -32602,
			Message: "Unsupported protocol version",
			Data: map[string]interface{}{
				"supported": s.ProtocolVersions(),
				"requested": protoVersion,
			},
		})
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		s.writeHTTPError(w, http.StatusBadRequest, &JSONRPCError{
			Code:    -32700,
			Message: "Parse error",
			Data:    "failed to read request body",
		})
		return
	}

	if len(body) == 0 {
		s.writeHTTPError(w, http.StatusBadRequest, &JSONRPCError{
			Code:    -32700,
			Message: "Parse error",
			Data:    "empty request body",
		})
		return
	}

	var req JSONRPCRequest
	if err := json.Unmarshal(body, &req); err != nil {
		s.writeHTTPError(w, http.StatusBadRequest, &JSONRPCError{
			Code:    -32700,
			Message: "Parse error",
			Data:    err.Error(),
		})
		return
	}

	if req.JSONRPC != "2.0" {
		s.writeHTTPError(w, http.StatusBadRequest, &JSONRPCError{
			Code:    -32600,
			Message: "Invalid Request",
			Data:    "jsonrpc must be 2.0",
		})
		return
	}

	hasID := req.ID != nil && len(req.ID) > 0 && string(req.ID) != "null"

	if !hasID {
		switch req.Method {
		case "notifications/initialized":
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusAccepted)
		}
		return
	}

	accept := r.Header.Get("Accept")
	s.routeRequest(w, req, accept)
}

func (s *Server) routeRequest(w http.ResponseWriter, req JSONRPCRequest, accept string) {
	switch req.Method {
	case "initialize":
		s.handleInitialize(w, req, accept)
	case "tools/list":
		s.handleListTools(w, req, accept)
	case "tools/call":
		s.handleCallTool(w, req, accept)
	case "ping":
		s.writeJSONRPCResponse(w, req.ID, map[string]interface{}{}, accept)
	default:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		resp := JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error:   &JSONRPCError{Code: -32601, Message: "Method not found", Data: fmt.Sprintf("unknown method: %s", req.Method)},
		}
		json.NewEncoder(w).Encode(resp)
	}
}

func (s *Server) handleInitialize(w http.ResponseWriter, req JSONRPCRequest, accept string) {
	var params InitializeParams
	if len(req.Params) > 0 {
		_ = json.Unmarshal(req.Params, &params)
	}

	negotiated := s.DefaultProtocolVersion()
	if params.ProtocolVersion != "" {
		if !s.supportsProtocolVersion(params.ProtocolVersion) {
			s.writeJSONRPCError(w, req.ID, -32602, "Unsupported protocol version", map[string]interface{}{
				"supported": s.ProtocolVersions(),
				"requested": params.ProtocolVersion,
			})
			return
		}
		negotiated = params.ProtocolVersion
	}

	result := InitializeResult{
		ProtocolVersion: negotiated,
		Capabilities: Capabilities{
			Tools: &ToolsCapability{ListChanged: false},
		},
		ServerInfo:   ServerInfo{Name: s.name, Version: s.version},
		Instructions: fmt.Sprintf("Mock campus %s - supports %d tools for student services", s.name, len(s.tools)),
	}
	s.writeJSONRPCResponse(w, req.ID, result, accept)
}

func (s *Server) handleListTools(w http.ResponseWriter, req JSONRPCRequest, accept string) {
	s.mu.RLock()
	tools := make([]Tool, len(s.tools))
	copy(tools, s.tools)
	s.mu.RUnlock()

	result := ListToolsResult{Tools: tools}
	s.writeJSONRPCResponse(w, req.ID, result, accept)
}

func (s *Server) handleCallTool(w http.ResponseWriter, req JSONRPCRequest, accept string) {
	var params CallToolParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		s.writeJSONRPCError(w, req.ID, -32602, "Invalid params", fmt.Sprintf("failed to parse params: %v", err))
		return
	}

	s.mu.RLock()
	handler, ok := s.handlers[params.Name]
	s.mu.RUnlock()

	if !ok {
		s.writeJSONRPCError(w, req.ID, -32602, "Invalid params", fmt.Sprintf("unknown tool: %s", params.Name))
		return
	}

	result, err := handler(context.Background(), params.Arguments)
	if err != nil {
		result = &CallToolResult{
			Content: []Content{{Type: "text", Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		}
	}

	s.writeJSONRPCResponse(w, req.ID, result, accept)
}

func (s *Server) writeJSONRPCResponse(w http.ResponseWriter, id json.RawMessage, result interface{}, accept string) {
	if shouldSSE(accept) {
		s.writeSSE(w, id, result)
	} else {
		s.writeJSON(w, id, result)
	}
}

func (s *Server) writeJSON(w http.ResponseWriter, id json.RawMessage, result interface{}) {
	resp := JSONRPCResponse{JSONRPC: "2.0", ID: id, Result: result}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) writeSSE(w http.ResponseWriter, id json.RawMessage, result interface{}) {
	resp := JSONRPCResponse{JSONRPC: "2.0", ID: id, Result: result}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "event: message\ndata: %s\n\n", data)
	if flusher, ok := w.(http.Flusher); ok {
		flusher.Flush()
	}
}

func (s *Server) writeJSONRPCError(w http.ResponseWriter, id json.RawMessage, code int, message string, data interface{}) {
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &JSONRPCError{Code: code, Message: message, Data: data},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) writeHTTPError(w http.ResponseWriter, status int, err *JSONRPCError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := JSONRPCResponse{
		JSONRPC: "2.0",
		Error:   err,
	}
	json.NewEncoder(w).Encode(resp)
}

func shouldSSE(accept string) bool {
	return accept == "text/event-stream" || contains(accept, "text/event-stream")
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
