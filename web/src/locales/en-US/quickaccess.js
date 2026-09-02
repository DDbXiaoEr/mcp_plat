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
  title: 'Quick Access',
  description: 'Choose a client and bind a server with an AccessKey to generate a config you can import directly.',
  gatewayNotConfigured: 'The API gateway has not been configured yet, so an access address cannot be generated. Please contact an administrator to complete the setup.',
  noClients: 'No quick-access clients are enabled yet. Please try again later.',
  noServers: 'There are no published MCP servers to connect to yet. Please wait until an administrator publishes one.',
  noAccessKey: "You don't have a usable AccessKey yet,",
  createKeyAction: 'Go create',
  selectClient: '① Select client',
  selectServerAndKey: '② Select server and AccessKey',
  serverLabel: 'MCP server',
  accessKeyLabel: 'AccessKey',
  openServerNote: 'This server has key authentication disabled, so no AccessKey is needed.',
  endpointLabel: 'Endpoint',
  authLabel: 'Authentication',
  authNone: 'No AccessKey required',
  generateConfig: '③ Generate config',
  noConfigWarn: 'Unable to generate the access config: the gateway has no default publish domain (defaultPublishDomain) configured, and this server has no full access address. Please contact an administrator to complete the setup.',
  importConfig: 'Import config for {client}',
  openInClient: 'Open in {client}',
  openInClientTitle: 'Open and import in {client}',
  setupSteps: 'Setup steps ({client})',
  clients: {
    cherrystudio: {
      description: 'A multi-model AI client that supports importing MCP servers via JSON.',
      step1: 'Click the “Open in {client}” button above to launch the app and confirm the import.',
      step2: 'If your browser asks whether to open {client}, choose Allow.',
      step3: 'In the popup, confirm importing this MCP server.',
      step4: 'After binding this MCP server in Agent editing, you can start using it.'
    }
  }
}
