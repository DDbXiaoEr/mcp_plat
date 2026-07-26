package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"mcp_plat-console/database"
	"mcp_plat-console/model"

	"github.com/google/uuid"
)

type NetworkSecuritySetting struct {
	Allowlist []string `json:"allowlist"`
}

var privateNetworks = []*net.IPNet{
	{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},     // loopback IPv4
	{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},      // private A
	{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},   // private B
	{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},  // private C
	{IP: net.IPv4(169, 254, 0, 0), Mask: net.CIDRMask(16, 32)},  // link-local
	{IP: net.IPv4(224, 0, 0, 0), Mask: net.CIDRMask(4, 32)},     // multicast
	{IP: net.IPv4(0, 0, 0, 0), Mask: net.CIDRMask(8, 32)},       // current network
}

func loadAllowedCIDRs() []*net.IPNet {
	var ns NetworkSecuritySetting
	if err := GetSetting("network_security", &ns); err != nil {
		return nil
	}
	var cidrs []*net.IPNet
	for _, cidrStr := range ns.Allowlist {
		cidrStr = strings.TrimSpace(cidrStr)
		if cidrStr == "" {
			continue
		}
		_, cidr, err := net.ParseCIDR(cidrStr)
		if err != nil {
			continue
		}
		cidrs = append(cidrs, cidr)
	}
	return cidrs
}

func isSafeIP(ip net.IP, allowedCIDRs []*net.IPNet) bool {
	if ip == nil {
		return false
	}
	for _, cidr := range allowedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	if ip.IsLoopback() {
		return true
	}
	if ip.IsPrivate() || ip.IsLinkLocalMulticast() ||
		ip.IsLinkLocalUnicast() || ip.IsMulticast() || ip.IsUnspecified() {
		return false
	}
	for _, n := range privateNetworks {
		if n.Contains(ip) {
			return false
		}
	}
	return true
}

func safeDialContext(dialer *net.Dialer, allowedCIDRs []*net.IPNet) func(ctx context.Context, network, addr string) (net.Conn, error) {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			host = addr
			port = ""
		}

		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
		if err != nil {
			return nil, fmt.Errorf("DNS 解析失败: %w", err)
		}
		for _, ip := range ips {
			if !isSafeIP(ip, allowedCIDRs) {
				return nil, fmt.Errorf("拒绝连接到内网地址: %s", ip.String())
			}
		}

		if port != "" {
			addr = net.JoinHostPort(host, port)
		}
		return dialer.DialContext(ctx, network, addr)
	}
}

func safeCheckRedirect(client *http.Client, allowedCIDRs []*net.IPNet) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return errors.New("重定向次数过多")
		}

		host := req.URL.Hostname()
		ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip", host)
		if err != nil {
			return fmt.Errorf("重定向目标 DNS 解析失败")
		}
		for _, ip := range ips {
			if !isSafeIP(ip, allowedCIDRs) {
				return fmt.Errorf("拒绝重定向到内网地址")
			}
		}
		return nil
	}
}

func safeHTTPClient(timeout time.Duration, allowedCIDRs []*net.IPNet) *http.Client {
	dialer := &net.Dialer{Timeout: timeout}
	c := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			DialContext:           safeDialContext(dialer, allowedCIDRs),
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			DisableCompression:    false,
			ResponseHeaderTimeout: timeout,
			ExpectContinueTimeout: 2 * time.Second,
		},
	}
	c.CheckRedirect = safeCheckRedirect(c, allowedCIDRs)
	return c
}

func validateMCPAddress(raw string, allowedCIDRs []*net.IPNet) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", fmt.Errorf("URL 解析失败: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return "", errors.New("仅支持 http/https 协议")
	}

	host := u.Hostname()
	ips, err := net.DefaultResolver.LookupIP(context.Background(), "ip", host)
	if err != nil {
		return "", fmt.Errorf("DNS 解析失败: %w", err)
	}
	for _, ip := range ips {
		if !isSafeIP(ip, allowedCIDRs) {
			return "", fmt.Errorf("不允许连接到内网地址")
		}
	}
	return raw, nil
}

type CreateServerInput struct {
	Name           string `json:"name" binding:"required"`
	Address        string `json:"address" binding:"required"`
	ServiceAddress string `json:"service_address"`
	Department     string `json:"department"`
	Protocol       string `json:"protocol"`
	Tools          string `json:"tools"`
}

type UpdateServerInput struct {
	Name           string `json:"name"`
	Address        string `json:"address"`
	ServiceAddress string `json:"service_address"`
	Department     string `json:"department"`
	Protocol       string `json:"protocol"`
	Tools          string `json:"tools"`
}

func ListServers() ([]model.MCPServer, error) {
	var servers []model.MCPServer
	err := database.DB.Order("created_at desc").Find(&servers).Error
	return servers, err
}

func CreateServer(input CreateServerInput) (*model.MCPServer, error) {
	protocol := input.Protocol
	if protocol == "" {
		protocol = "streamable http"
	}
	server := model.MCPServer{
		ID:             uuid.New().String(),
		Name:           input.Name,
		Address:        input.Address,
		ServiceAddress: input.ServiceAddress,
		Department:     input.Department,
		Protocol:       protocol,
		Tools:          input.Tools,
	}
	if err := database.DB.Create(&server).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

func UpdateServer(id string, input UpdateServerInput) error {
	var server model.MCPServer
	if err := database.DB.First(&server, id).Error; err != nil {
		return errors.New("MCP 服务器不存在")
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Address != "" {
		updates["address"] = input.Address
	}
	if input.ServiceAddress != "" {
		updates["service_address"] = input.ServiceAddress
	}
	if input.Department != "" {
		updates["department"] = input.Department
	}
	if input.Protocol != "" {
		updates["protocol"] = input.Protocol
	}
	if input.Tools != "" {
		updates["tools"] = input.Tools
	}

	return database.DB.Model(&server).Updates(updates).Error
}

func DeleteServer(id string) error {
	database.DB.Where("server_id = ?", id).Delete(&model.RoleServer{})
	result := database.DB.Where("id = ?", id).Delete(&model.MCPServer{})
	if result.RowsAffected == 0 {
		return errors.New("MCP 服务器不存在")
	}
	return result.Error
}

type FetchToolsInput struct {
	Address  string `json:"address" binding:"required"`
	Protocol string `json:"protocol"`
}

type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolsListResult struct {
	Tools []toolInfo `json:"tools"`
}

type toolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func FetchTools(input FetchToolsInput) ([]toolInfo, error) {
	protocol := input.Protocol
	if protocol == "" {
		protocol = "Streamable HTTP"
	}

	allowedCIDRs := loadAllowedCIDRs()

	address := input.Address
	if !strings.Contains(address, "://") {
		if strings.HasPrefix(address, "/") {
			var gw ApiGatewaySetting
			if err := GetSetting("api_gateway", &gw); err != nil || gw.DefaultPublishDomain == "" {
				return nil, errors.New("请先在系统设置中配置 API 网关默认域名")
			}
			address = strings.TrimRight(gw.DefaultPublishDomain, "/") + address
		} else {
			address = "http://" + address
		}
	}

	validatedAddr, err := validateMCPAddress(address, allowedCIDRs)
	if err != nil {
		return nil, err
	}

	switch protocol {
	case "Streamable HTTP":
		return fetchToolsStreamableHTTP(validatedAddr, allowedCIDRs)
	case "SSE":
		return fetchToolsSSE(validatedAddr, allowedCIDRs)
	case "stdio":
		return nil, errors.New("stdio 协议暂不支持远程获取工具列表")
	default:
		return nil, errors.New("不支持的协议类型: " + protocol)
	}
}

func fetchToolsStreamableHTTP(address string, allowedCIDRs []*net.IPNet) ([]toolInfo, error) {
	client := safeHTTPClient(15*time.Second, allowedCIDRs)

	initReq := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
		Params: map[string]interface{}{
			"protocolVersion": "2025-03-26",
			"capabilities":    map[string]interface{}{},
			"clientInfo": map[string]interface{}{
				"name":    "mcp_plat-console",
				"version": "1.0.0",
			},
		},
	}

	initResp, sessionID, err := postJSONRPC(client, address, "", initReq)
	if err != nil {
		return nil, fmt.Errorf("初始化 MCP 会话失败: %w", err)
	}
	if initResp.Error != nil {
		return nil, fmt.Errorf("MCP 初始化返回错误: %s (code: %d)", initResp.Error.Message, initResp.Error.Code)
	}

	notifyBody, _ := json.Marshal(map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "notifications/initialized",
	})
	notifyReq, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(notifyBody))
	if err == nil {
		notifyReq.Header.Set("Content-Type", "application/json")
		notifyReq.Header.Set("Accept", "application/json, text/event-stream")
		if sessionID != "" {
			notifyReq.Header.Set("Mcp-Session-Id", sessionID)
		}
		if notifyResp, err := client.Do(notifyReq); err == nil {
			io.Copy(io.Discard, notifyResp.Body)
			notifyResp.Body.Close()
		}
	}

	listReq := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
		Params:  map[string]interface{}{},
	}

	result, _, err := postJSONRPC(client, address, sessionID, listReq)
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, fmt.Errorf("MCP 服务器返回错误: %s (code: %d)", result.Error.Message, result.Error.Code)
	}

	var toolsResult toolsListResult
	if err := json.Unmarshal(result.Result, &toolsResult); err != nil {
		return nil, fmt.Errorf("解析工具列表失败: %w", err)
	}

	return toolsResult.Tools, nil
}

func postJSONRPC(client *http.Client, address, sessionID string, reqBody jsonRPCRequest) (*jsonRPCResponse, string, error) {
	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, "", fmt.Errorf("构建请求失败: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, address, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, "", fmt.Errorf("构建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")
	if sessionID != "" {
		req.Header.Set("Mcp-Session-Id", sessionID)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("连接 MCP 服务器失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		io.Copy(io.Discard, resp.Body)
		return nil, "", fmt.Errorf("MCP 服务器返回错误状态 %d", resp.StatusCode)
	}

	newSessionID := resp.Header.Get("Mcp-Session-Id")
	if newSessionID == "" {
		newSessionID = sessionID
	}

	var result jsonRPCResponse
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/event-stream") {
		data, err := readSSEData(resp.Body)
		if err != nil {
			return nil, "", err
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, "", fmt.Errorf("解析响应失败: %w", err)
		}
	} else {
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, "", fmt.Errorf("解析响应失败: %w", err)
		}
	}

	return &result, newSessionID, nil
}

func readSSEData(body io.Reader) ([]byte, error) {
	scanner := bufio.NewScanner(body)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			return []byte(strings.TrimPrefix(line, "data: ")), nil
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取 SSE 响应失败: %w", err)
	}
	return nil, errors.New("未从 SSE 响应中读取到数据")
}

func fetchToolsSSE(address string, allowedCIDRs []*net.IPNet) ([]toolInfo, error) {
	client := safeHTTPClient(15*time.Second, allowedCIDRs)

	resp, err := client.Get(address)
	if err != nil {
		return nil, fmt.Errorf("连接 SSE 服务器失败: %w", err)
	}
	defer resp.Body.Close()

	var messageEndpoint string

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			eventType := strings.TrimPrefix(line, "event: ")
			if strings.EqualFold(eventType, "endpoint") {
				if scanner.Scan() {
					dataLine := scanner.Text()
					if strings.HasPrefix(dataLine, "data: ") {
						messageEndpoint = strings.TrimPrefix(dataLine, "data: ")
						messageEndpoint = strings.TrimSpace(messageEndpoint)

						parsed, err := url.Parse(messageEndpoint)
						if err != nil || parsed.IsAbs() {
							return nil, errors.New("SSE 端点不允许使用绝对 URL")
						}
						base, err := url.Parse(address)
						if err != nil {
							return nil, errors.New("解析 SSE 端点地址失败")
						}
						endpointURL, err := url.Parse(messageEndpoint)
						if err != nil {
							return nil, errors.New("解析 SSE 端点相对路径失败")
						}
						messageEndpoint = base.ResolveReference(endpointURL).String()
					}
				}
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取 SSE 流失败: %w", err)
	}

	if messageEndpoint == "" {
		return nil, errors.New("未能获取 SSE 消息端点")
	}

	tools, err := fetchToolsStreamableHTTP(messageEndpoint, allowedCIDRs)
	if err != nil {
		return nil, fmt.Errorf("通过 SSE 端点获取工具列表失败: %w", err)
	}

	return tools, nil
}

type PublishInput struct {
	ServerIDs  []string `json:"server_ids" binding:"required"`
	EnableAuth bool     `json:"enable_auth"`
}

type ApiGatewaySetting struct {
	Provider             string `json:"provider"`
	AdminURL             string `json:"adminUrl"`
	AdminKey             string `json:"adminKey"`
	DefaultPublishDomain string `json:"defaultPublishDomain"`
	AuthGrpcAddr         string `json:"authGrpcAddr"`
}

func PublishServers(input PublishInput) error {
	var gw ApiGatewaySetting
	if err := GetSetting("api_gateway", &gw); err != nil {
		return errors.New("请先在系统设置中配置 API 网关参数")
	}
	if gw.AdminURL == "" || gw.AdminKey == "" {
		return errors.New("API 网关 Admin API 地址和 Key 未配置")
	}

	var servers []model.MCPServer
	if err := database.DB.Where("id IN ?", input.ServerIDs).Find(&servers).Error; err != nil {
		return fmt.Errorf("查询服务器失败: %w", err)
	}
	if len(servers) == 0 {
		return errors.New("未找到指定的服务器")
	}

	client := &http.Client{Timeout: 10 * time.Second}
	var errs []string

	for _, srv := range servers {
		if srv.ServiceAddress == "" {
			errs = append(errs, fmt.Sprintf("%s: 未配置后端服务地址", srv.Name))
			continue
		}

		upstreamBody := map[string]interface{}{
			"name":  srv.Name,
			"type":  "roundrobin",
			"nodes": map[string]int{srv.ServiceAddress: 1},
		}
		if err := putAPISIXAdmin(client, gw.AdminURL, gw.AdminKey, "/apisix/admin/upstreams/"+srv.ID, upstreamBody); err != nil {
			errs = append(errs, fmt.Sprintf("%s: 创建上游失败: %v", srv.Name, err))
			continue
		}

		routeBody := map[string]interface{}{
			"name":        srv.Name,
			"uris":        []string{"/" + srv.ID, "/" + srv.ID + "/*"},
			"upstream_id": srv.ID,
		}
		if gw.DefaultPublishDomain != "" {
			routeBody["host"] = gw.DefaultPublishDomain
		}

		plugins := map[string]interface{}{}
		if srv.Address != "" && srv.Address != "/" {
			plugins["proxy-rewrite"] = map[string]interface{}{
				"regex_uri": []string{
					"^/" + srv.ID + "(.*)",
					srv.Address + "$1",
				},
			}
		}
		if input.EnableAuth {
			if gw.AuthGrpcAddr == "" {
				errs = append(errs, fmt.Sprintf("%s: 启用认证但未配置 gRPC 校验地址", srv.Name))
				continue
			}
			authValue := fmt.Sprintf(`{"header_name":"X-Access-Key","grpc_addr":"%s","server_id":"%s"}`,
				gw.AuthGrpcAddr, srv.ID)
			plugins["ext-plugin-pre-req"] = map[string]interface{}{
				"conf": []map[string]interface{}{
					{
						"name":  "accesskey_verify",
						"value": authValue,
					},
				},
			}
		}
		if len(plugins) > 0 {
			routeBody["plugins"] = plugins
		}

		if err := putAPISIXAdmin(client, gw.AdminURL, gw.AdminKey, "/apisix/admin/routes/"+srv.ID, routeBody); err != nil {
			errs = append(errs, fmt.Sprintf("%s: 创建路由失败: %v", srv.Name, err))
			continue
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("部分发布失败: %s", strings.Join(errs, "; "))
	}

	return nil
}

func putAPISIXAdmin(client *http.Client, adminURL, adminKey, apiPath string, body interface{}) error {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}

	urlStr := strings.TrimRight(adminURL, "/") + apiPath
	req, err := http.NewRequest(http.MethodPut, urlStr, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("构建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", adminKey)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求 APISIX Admin API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("APISIX 返回错误 (status=%d)", resp.StatusCode)
	}

	return nil
}
