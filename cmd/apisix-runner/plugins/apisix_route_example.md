# APISIX 路由配置示例

## 1. 创建路由，绑定 accesskey_verify + audit_log 插件

audit_log 需要同时在 `ext-plugin-pre-req`（请求阶段捕获元数据）和
`ext-plugin-post-resp`（响应阶段按真实结果上报成功/失败）两个阶段挂载，
缺少 post-resp 阶段时请求会由 60s 兜底任务统一记为失败，无法反映真实结果。

curl http://127.0.0.1:9180/apisix/admin/routes/1 \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' -X PUT -d '{
  "uri": "/mcp/*",
  "plugins": {
    "ext-plugin-pre-req": {
      "conf": [
        {
          "name": "accesskey_verify",
          "value": "{\"header_name\":\"X-Access-Key\",\"grpc_addrs\":[\"127.0.0.1:9090\",\"127.0.0.1:9091\"]}"
        },
        {
          "name": "audit_log",
          "value": "{\"header_name\":\"X-Access-Key\",\"server_id\":\"<server-uuid>\",\"grpc_addrs\":[\"127.0.0.1:9091\"]}"
        }
      ]
    },
    "ext-plugin-post-resp": {
      "conf": [
        {
          "name": "audit_log",
          "value": "{\"header_name\":\"X-Access-Key\",\"server_id\":\"<server-uuid>\",\"grpc_addrs\":[\"127.0.0.1:9091\"]}"
        }
      ]
    }
  },
  "upstream": {
    "type": "roundrobin",
    "nodes": {
      "127.0.0.1:8080": 1
    }
  }
}'

## 2. 插件配置字段说明

# header_name: 从哪个 HTTP Header 提取 access key（默认 "X-Access-Key"）
# server_id:   该路由对应的 MCP 服务器 ID，写入审计日志用
# grpc_addrs:  gRPC 校验服务地址列表（必填，支持多个后端，插件按连接数做最少连接负载均衡）
# grpc_addr:   兼容旧配置，单个 gRPC 校验服务地址；配置了 grpc_addrs 时优先使用 grpc_addrs

## 3. 成功判定口径

audit_log 在响应阶段记录每条请求并判定是否成功（Success 字段）：
- HTTP 非 2xx → 失败，Message 为 "HTTP <code> <status>"
- HTTP 2xx 但响应体为 JSON-RPC 且含顶层 error 成员 → 失败（rpc error code=...）
- HTTP 2xx 且 result.isError=true → 失败（tool returned isError）
- 其余（含工具正常返回 / 非 JSON 响应体）→ 成功
- 超过 60s 仍未进入响应阶段（前置拦截、上游无响应）→ 兜底记为失败

## 4. 测试请求

# 正常请求（携带有效 access key）
curl http://127.0.0.1:9080/mcp/<server-uuid> \
  -H "X-Access-Key: ak-<your_jwt_token>"

# 缺少 access key，应返回 401
curl http://127.0.0.1:9080/mcp/<server-uuid>
