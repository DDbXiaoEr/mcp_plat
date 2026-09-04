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

export default {
  tabs: {
    ops: '运维设置',
    operation: '运营设置'
  },
  general: {
    saved: '已保存',
    saveFailed: '保存失败',
    remove: '移除',
    gotIt: '知道了',
    addAddress: '+ 添加地址'
  },
  log: {
    title: '日志设置',
    enableSyslog: '启用 Syslog',
    offNotice: '未启用 Syslog 时，日志同时输出到标准输出与本地文件（支持轮转）。',
    onNotice: '启用 Syslog 后，标准输出与本地文件日志均失效，仅输出到 Syslog 服务器。',
    pathLabel: '日志目录',
    levelLabel: '日志级别',
    fileNameLabel: '日志文件名',
    fileNameHelp: '日志文件为 \\{前缀\\}.log，默认 mcp_plat.log',
    rotation: '日志轮转',
    maxSizeLabel: '单文件大小上限（MB）',
    maxBackupsLabel: '保留历史文件数',
    maxAgeLabel: '保留天数',
    compressLabel: '压缩历史日志',
    compressNone: '不压缩',
    compressGzip: 'gzip 压缩',
    syslogHostLabel: 'Syslog 服务器地址',
    syslogHostPlaceholder: '请输入 Syslog 服务器地址',
    portLabel: '端口',
    protocolLabel: '协议'
  },
  mail: {
    title: 'SMTP 设置',
    enable: '启用邮件通知',
    toggleHint: '开关控制是否实际发送通知邮件；下方配置项始终可编辑，用于测试与预先配置',
    hostLabel: 'SMTP 服务器地址',
    portLabel: '端口',
    encryptionLabel: '加密方式',
    encryptionNone: '无',
    usernameLabel: '用户名',
    usernamePlaceholder: '请输入 SMTP 用户名',
    passwordLabel: '密码',
    passwordPlaceholder: '请输入 SMTP 密码',
    fromAddressLabel: '发件人地址',
    fromNameLabel: '发件人名称',
    fromNamePlaceholder: 'MCP 服务平台',
    sendTest: '发送测试邮件',
    sending: '发送中…',
    testToPlaceholder: '输入收件人邮箱，验证 SMTP 配置',
    testOk: '测试邮件发送成功，请检查收件箱',
    sendFailed: '发送失败'
  },
  apiGw: {
    title: 'API 网关设置',
    configuredWarning: '已检测到当前系统配置了 API 网关，修改保存后会覆盖当前配置。',
    providerLabel: 'API 网关',
    kongHelpPre: 'Kong 仅发布上游与服务路由，不下发插件；本地 Docker 测试时 Admin API 默认 ', 
    kongHelpPost: '，未启用 RBAC 时 Admin Key 可留空。',
    adminUrlLabel: '网关 Admin API 地址',
    adminUrlTooltipKongPre: '填写 Kong Admin API 的基础地址，不要附带路径，例如 ', 
    adminUrlTooltipKongPost: '。',
    adminUrlTooltipDefaultPre: '填写网关 Admin API 的基础地址，不要附带路径，例如 ', 
    adminUrlTooltipDefaultPost: '（而非 .../apisix/admin）。',
    defaultPublishLabel: '默认发布域名',
    defaultPublishPlaceholder: '例如 mcp.xauat.edu.cn',
    adminKeyLabel: 'Admin API Key',
    adminKeyPlaceholder: '请输入管理员 API Key',
    adminKeyKongHint: 'Kong 未启用 RBAC 时可不填；启用后填写 RBAC Admin 令牌。',
    authGrpcLabel: 'Access Key 认证 gRPC 地址',
    authGrpcTip1: '对应 accesskey-auth-server 的 gRPC 监听地址，例如 ', 
    authGrpcTip2: ' 或 ',
    authGrpcTip3: '。支持配置多个地址，发布时会下发到 accesskey_verify 插件，并按连接数（最少连接）负载均衡。每行一个地址。',
    headerLabel: 'Access Key Header 名称',
    headerTip1: '指定请求中携带 Access Key 的 HTTP Header 名称。',
    headerTip2: '发布路由时若已配置则使用自定义 Header，未配置则默认使用 ', 
    headerTip3: '。'
  },
  kongDialog: {
    title: 'Kong 网关接入说明',
    intro1: '当前 API 网关已切换为 ',
    intro2: '，发布 MCP 服务时请注意以下限制：',
    item1a: '仅发布「上游（Upstream）」「Service」「路由（Route）」，',
    item1b: '不会下发任何插件',
    item1c: '。',
    item2: 'Access Key 认证、审计日志、路径重写等能力在 Kong 下不生效。',
    item3: '发布后的网关路径使用 MCP 服务器配置的 URI 路径（address），而非服务器 UUID。',
    dockerPre: '本地 Docker 测试：Admin API 默认 ', 
    dockerPost: '，未启用 RBAC 时 Admin API Key 可留空。'
  },
  network: {
    title: '网络安全',
    allowLabel: '允许连接的内网 CIDR',
    tip1: '每行一个 CIDR 地址段，例如 ', 
    tip2: '。',
    tip3: '配置后，获取工具列表时将允许连接这些内网地址。',
    tip4: '如需连接 MCP 服务器内网 IP，在此添加对应的地址段即可。'
  },
  audit: {
    title: '审计日志',
    enable: '启用审计日志',
    grpcAddrLabel: '审计日志 gRPC 地址',
    tip1: '对应 audit-log-server 的 gRPC 监听地址，例如',
    tip2: ' 或 ',
    tip3: '。支持配置多个地址，发布时会下发到 audit_log 插件，并按连接数（最少连接）负载均衡。每行一个地址。'
  },
  platform: {
    title: '平台设置',
    nameLabel: '平台名称',
    namePlaceholder: '请输入平台名称，如：某某大学',
    logoLabel: 'Logo 图片地址',
    logoPlaceholder: '请输入 Logo 图片的 URL 地址',
    loginBgLabel: '登录页背景图',
    loginBgPlaceholder: '请输入登录页背景图片的 URL 地址',
    loginBgHint: '留空则使用默认渐变背景；填写后登录页自适应填充显示该图片',
    siteUrlLabel: '跳转链接',
    siteUrlPlaceholder: '点击 Logo 时跳转到的网址，如：https://www.example.edu.cn'
  },
  auth: {
    title: '用户认证',
    methodLabel: '认证方式',
    cas: {
      serverUrlLabel: 'CAS 服务器地址',
      serviceUrlLabel: '服务回调地址',
      versionLabel: '协议版本',
      mappingTitle: '属性映射（平台字段 → CAS 属性）',
      selectPlatformField: '选择平台字段',
      attrPlaceholder: 'CAS 属性名',
      addMapping: '+ 添加映射'
    },
    ldap: {
      enabledNotice: 'LDAP 认证已启用。出于安全考虑，不显示当前已配置的值；输入框留空表示保持现有配置不变。',
      hostLabel: '服务器地址',
      portLabel: '端口',
      bindPasswordLabel: 'Bind 密码',
      bindPasswordPlaceholder: '请输入 Bind 密码',
      userFilterLabel: '用户过滤器',
      mappingTitle: '属性映射（平台字段 → LDAP 属性）',
      selectPlatformField: '选择平台字段',
      ldapAttrPlaceholder: 'LDAP 属性名',
      addMapping: '+ 添加映射',
      testTitle: '测试映射',
      testing: '测试中…',
      testPlaceholder: '输入用户名测试映射效果',
      testFailed: '测试失败',
      allAttributes: 'LDAP 全部属性（{count} 个）',
      currentMapping: '当前映射结果',
      noMapping: '暂无映射配置',
      saveHint1: '请确保已点击上方「保存」按钮保存认证设置，保存后用户需',
      saveHint2: '重新登录',
      saveHint3: '才能同步 LDAP 属性到个人信息。'
    },
    oauth: {
      authorizeUrlLabel: '授权地址',
      tokenUrlLabel: 'Token 地址',
      userinfoUrlLabel: '用户信息地址',
      clientIdPlaceholder: '请输入 Client ID',
      clientSecretPlaceholder: '请输入 Client Secret',
      redirectUrlLabel: '回调地址'
    }
  },
  userOps: {
    title: '用户运营',
    maxAccessKeysLabel: '每用户最大 AccessKey 数量',
    accessKeyCronLabel: '过期扫描 cron 表达式',
    accessKeyCronHint: '每隔一段时间自动扫描过期 AccessKey 并禁用，支持标准 crontab 格式',
    roleRuleTitle: '新用户角色自动分配规则',
    filterAttrLabel: '过滤属性',
    filterAttrHint: '所有角色条件共用同一个过滤属性，用户该属性匹配到哪条规则就分配对应角色',
    conditionIndex: '条件 {n}',
    satisfies: '满足',
    patternPlaceholder: '正则表达式，如 ^\\d{7}$',
    assignRole: '自动分配角色',
    selectRole: '选择角色',
    removeCondition: '移除条件',
    addCondition: '+ 添加条件'
  },
  fields: {
    username: '用户名',
    uid: 'UID',
    name: '姓名',
    email: '邮箱',
    phone: '电话',
    organization: '组织'
  },
  quickAccess: {
    supportedClients: '支持的客户端',
    supportedClientsHint: '控制用户「快速接入」页可选用的客户端，勾选后该客户端对用户开放。',
    schemeLabel: '接入协议',
    schemeHint: '决定生成接入地址使用的协议（HTTP 或 HTTPS）。'
  }
}
