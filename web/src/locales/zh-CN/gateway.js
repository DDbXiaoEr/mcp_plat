/*
 * Copyright (C) 2026 Zhaoquan Wang
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Author: deepseek-v4-flash / opencode
export default {
  title: 'API 网关服务器管理',
  desc: '管理平台使用的 API 网关实例，发布 MCP 服务器时使用默认网关。',
  add: '新增网关',
  refresh: '刷新状态',
  loading: '加载中…',
  refreshing: '刷新中…',
  empty: '暂无网关服务器',
  emptyHint: '添加 API 网关实例后，即可向网关发布 MCP 服务器。',
  defaultTag: '默认',
  provider: {
    apisix: 'Apache APISIX',
    kong: 'Kong'
  },
  status: {
    online: '在线',
    offline: '离线',
    unknown: '未检测'
  },
  card: {
    adminUrl: 'Admin API',
    accessKeyHeader: 'Access Key Header',
    authGrpcAddrs: '认证 gRPC 地址',
    noGrpc: '未配置',
    latency: '响应 {ms} ms',
    noDescription: '暂无说明'
  },
  actions: {
    test: '测试连接',
    testing: '测试中…',
    setDefault: '设为默认',
    edit: '编辑'
  },
  form: {
    createTitle: '新增网关服务器',
    editTitle: '编辑网关服务器',
    name: '网关名称',
    namePlaceholder: '例如：主网关',
    provider: '网关类型',
    adminUrl: 'Admin API 地址',
    adminUrlPlaceholder: 'http://127.0.0.1:9180',
    adminKey: 'Admin API Key',
    adminKeyPlaceholder: '请输入 Admin API Key',
    adminKeyKeep: '留空则保持现有 Key 不变',
    accessKeyHeader: 'Access Key Header',
    accessKeyHeaderPlaceholder: 'X-Access-Key',
    authGrpcAddrs: '认证 gRPC 地址',
    grpcPlaceholder: '127.0.0.1:9090',
    addAddress: '添加地址',
    grpcHint: '对应 accesskey-auth-server 的 gRPC 监听地址，每行一个，发布时下发到 accesskey_verify 插件。',
    description: '说明',
    descriptionPlaceholder: '可选，备注该网关的用途',
    providerHintApisix: '填写 APISIX Admin API 的基础地址，例如 http://127.0.0.1:9180（不要带 /apisix/admin）。',
    providerHintKong: '填写 Kong Admin API 的基础地址，例如 http://127.0.0.1:8001；未启用 RBAC 时 Admin Key 可留空。',
    test: '测试连通性',
    testing: '测试中…',
    testOk: '连接成功',
    testFail: '连接失败'
  },
  delete: {
    title: '删除网关服务器',
    desc: '确定删除网关「{name}」吗？删除后不可恢复，已发布的路由不会自动从该网关下线。',
    confirm: '删除',
    deleting: '删除中…'
  },
  alert: {
    createSuccess: '网关服务器创建成功',
    updateSuccess: '网关服务器更新成功',
    deleteSuccess: '网关服务器已删除',
    setDefaultSuccess: '已设为默认网关',
    loadFailed: '加载网关列表失败',
    saveFailed: '保存失败',
    deleteFailed: '删除失败'
  }
}
