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

// Author: deepseek-v4-pro / opencode
export default {
  title: '快速接入',
  description: '选择客户端并绑定服务器与 AccessKey，即可生成可直接导入的配置。',
  gatewayNotConfigured: '平台尚未配置 API 网关，无法生成接入地址，请先联系管理员完成配置。',
  noClients: '管理员暂未开放任何快速接入客户端，请稍后再试。',
  noServers: '暂无可接入的已发布 MCP 服务器，请等待管理员发布后再试。',
  noAccessKey: '您还没有可用的 AccessKey，',
  createKeyAction: '去创建',
  selectClient: '① 选择客户端',
  selectServerAndKey: '② 选择服务器与 AccessKey',
  serverLabel: 'MCP 服务器',
  accessKeyLabel: 'AccessKey',
  openServerNote: '该服务器未开启 Key 验证，无需选择 AccessKey',
  endpointLabel: '接入地址',
  authLabel: '鉴权',
  authNone: '无需 AccessKey',
  generateConfig: '③ 生成配置',
  noConfigWarn: '无法生成接入配置：网关未配置默认发布域名（defaultPublishDomain），且该服务器未填写完整接入地址，请先联系管理员配置。',
  importConfig: '{client} 导入配置',
  openInClient: '在 {client} 中打开',
  openInClientTitle: '在 {client} 中打开并导入',
  setupSteps: '接入步骤（{client}）',
  clients: {
    cherrystudio: {
      description: '多模型 AI 客户端，支持通过 JSON 导入 MCP 服务器',
      step1: '点击上方「在 CherryStudio 中打开」按钮，一键唤起应用并弹出导入确认',
      step2: '若浏览器提示是否打开 CherryStudio，请选择允许',
      step3: '在弹窗中确认导入该 MCP 服务器',
      step4: '在 Agent 编辑中绑定该 MCP 服务器后即可使用'
    }
  }
}
