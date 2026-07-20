# APISIX 路由配置示例

## 1. 创建路由，绑定 accesskey_verify 插件

curl http://127.0.0.1:9180/apisix/admin/routes/1 \
-H 'X-API-KEY: edd1c9f034335f136f87ad84b625c8f1' -X PUT -d '{
  "uri": "/api/*",
  "plugins": {
    "ext-plugin-pre-req": {
      "conf": [
        {
          "name": "accesskey_verify",
          "value": "{\"header_name\":\"X-Access-Key\",\"grpc_addr\":\":9090\"}"
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
# grpc_addr:   gRPC 校验服务地址（必填，对应 grpc_server.yaml 中 grpc_addr）

## 3. 测试请求

# 正常请求（携带有效 access key）
curl http://127.0.0.1:9080/api/test \
  -H "X-Access-Key: ak-<your_jwt_token>"

# 缺少 access key，应返回 401
curl http://127.0.0.1:9080/api/test
