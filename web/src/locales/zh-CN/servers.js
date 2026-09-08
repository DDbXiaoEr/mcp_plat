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
  title: 'MCP 服务器管理',
  add: '增加',
  department: '部门',
  backToEdit: '返回修改',
  confirmPush: '确认推送',
  notFilled: '未填写',
  status: {
    published: '已发布',
    maintenance: '维护中',
    unpublished: '未发布'
  },
  detail: {
    uriPath: 'MCP服务的URI',
    serviceAddress: 'MCP 服务地址',
    responsibleDepartment: '负责部门',
    protocolType: '协议类型',
    protocolVersion: '协议版本',
    tools: '工具列表',
    publishStatus: '发布状态',
    noDescription: '暂无描述',
    descriptionPlaceholder: '请输入描述信息'
  },
  form: {
    title: '新增 MCP 服务器',
    namePlaceholder: '请输入服务器名称',
    uriPathLabel: 'MCP服务器URI路径',
    uriPathHelp: '填写 MCP 提供服务的 URI，例如 {addr} 的话就填写 {path}。',
    uriPathPlaceholder: '{path}（仅路径，不含域名）',
    serviceAddressHelp: '填写 MCP 服务的实际 IP:端口，例如 {addr}，支持多个地址（用逗号或换行分隔）。',
    departmentPlaceholder: '请输入负责部门',
    descriptionPlaceholder: '请输入服务器描述信息',
    selectAddress: '选择地址',
    selectAddressPlaceholder: '请选择地址',
    customAddressOption: '自定义地址...',
    customAddressLabel: '自定义地址',
    customAddressPlaceholder: 'IP:端口，例如 {example}',
    connectionMethod: '连接方式',
    protocolVersionHelp: 'MCP 协议规范版本，获取工具列表时用于与服务器协商。默认使用最新版 {version}。',
    toolsPlaceholder: '多个工具以逗号或换行分隔',
    confirm: '确认新增'
  },
  fetch: {
    sectionTitle: '获取工具',
    fetching: '获取中...',
    button: '自动获取工具列表',
    requireAddress: '请先选择或输入 MCP 服务地址',
    failed: '获取工具列表失败'
  },
  shuttle: {
    available: '可选服务器',
    selected: '已选服务器',
    filterPlaceholder: '过滤名称/UUID',
    selectAll: '全选',
    emptyAvailable: '暂无可选服务器',
    emptySelected: '暂未选择'
  },
  publish: {
    button: '发布',
    title: '发布到 API 网关',
    desc: '选择要发布的 MCP 服务器：',
    gatewayConfigured: 'API 网关已配置{kong}',
    kongSuffix: '（Kong）',
    kongHint: 'Kong 仅发布上游与路由，不会下发认证/审计插件，以下开关在 Kong 下不生效。',
    optionAuth: '启用 Access Key 认证',
    authHint: '启用后请求本 MCP 服务时需要携带有效的 access key',
    optionAudit: '启用审计日志',
    auditHint: '启用后记录每次 MCP 服务调用的访问日志',
    confirm: '确认发布',
    confirmTitle: '确认发布配置',
    confirmDesc: '即将向 API 网关推送以下路由配置：',
    progress: '发布中...'
  },
  preview: {
    serviceName: '服务名称',
    gatewayPath: '网关路径',
    backendAddress: '后端地址',
    authStatus: '认证状态',
    authEnabled: '已启用 Access Key 认证',
    auditLog: '审计日志',
    enabled: '已启用',
    disabled: '未启用'
  },
  maintenance: {
    button: '维护',
    title: '设置服务器维护',
    desc: '选择要进入维护的服务器（移入右侧），或将维护中的服务器移回左侧取消维护：',
    underMaintenance: '维护中',
    emptyUnderMaintenance: '暂无维护中的服务器',
    confirm: '确认维护',
    confirmCancel: '确认取消维护',
    confirmBoth: '确认维护配置',
    confirmDesc: '即将向 API 网关推送以下维护配置，维护中的路由将直接返回 503：',
    sectionEnter: '进入维护',
    sectionRestore: '取消维护',
    submitting: '提交中...'
  },
  alert: {
    deleteConfirm: '确认删除',
    deleteBody: '确定要删除 MCP 服务器「{name}」吗？',
    deleteWarn: '注意：此操作仅从数据库中移除记录，不会自动从 API 网关下线服务器，请手动到 API 网关删除相关路由和上游规则。',
    deleting: '删除中...',
    deleteSuccess: '删除成功。请手动到 API 网关下线路由和上游规则。',
    publishSuccess: '发布成功',
    publishFailed: '发布失败',
    maintenanceSuccess: '维护设置成功',
    operationFailed: '操作失败'
  }
}
