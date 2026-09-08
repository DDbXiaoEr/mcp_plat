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
  title: 'MCP Server Management',
  add: 'Add',
  department: 'Department',
  backToEdit: 'Back to edit',
  confirmPush: 'Confirm & push',
  notFilled: 'Not provided',
  status: {
    published: 'Published',
    maintenance: 'In maintenance',
    unpublished: 'Unpublished'
  },
  detail: {
    uriPath: 'MCP service URI',
    serviceAddress: 'MCP service address',
    responsibleDepartment: 'Responsible department',
    protocolType: 'Protocol type',
    protocolVersion: 'Protocol version',
    tools: 'Tools',
    publishStatus: 'Publish status',
    noDescription: 'No description',
    descriptionPlaceholder: 'Please enter a description'
  },
  form: {
    title: 'Add MCP server',
    namePlaceholder: 'Please enter the server name',
    uriPathLabel: 'MCP server URI path',
    uriPathHelp: 'Enter the URI where the MCP service is served. For example, if it is available at {addr}, fill in {path}.',
    uriPathPlaceholder: '{path} (path only, no domain)',
    serviceAddressHelp: 'Enter the actual IP:port of the MCP service, e.g. {addr}. Multiple addresses are supported (separate with commas or line breaks).',
    departmentPlaceholder: 'Please enter the responsible department',
    descriptionPlaceholder: 'Please enter the server description',
    selectAddress: 'Select address',
    selectAddressPlaceholder: 'Select an address',
    customAddressOption: 'Custom address...',
    customAddressLabel: 'Custom address',
    customAddressPlaceholder: 'IP:port, e.g. {example}',
    connectionMethod: 'Connection method',
    protocolVersionHelp: 'The MCP protocol specification version used to negotiate with the server when fetching tools. Defaults to the latest {version}.',
    toolsPlaceholder: 'Separate multiple tools with commas or line breaks',
    confirm: 'Confirm add'
  },
  fetch: {
    sectionTitle: 'Fetch tools',
    fetching: 'Fetching...',
    button: 'Auto-fetch tool list',
    requireAddress: 'Please select or enter an MCP service address first',
    failed: 'Failed to fetch the tool list'
  },
  shuttle: {
    available: 'Available servers',
    selected: 'Selected servers',
    filterPlaceholder: 'Filter by name/UUID',
    selectAll: 'Select all',
    emptyAvailable: 'No servers available',
    emptySelected: 'Nothing selected yet'
  },
  publish: {
    button: 'Publish',
    title: 'Publish to API gateway',
    desc: 'Select the MCP servers to publish:',
    gatewayConfigured: 'API gateway is configured{kong}',
    kongSuffix: ' (Kong)',
    kongHint: 'Kong only publishes upstreams and routes; it does not install authentication/audit plugins, so the switches below have no effect under Kong.',
    optionAuth: 'Enable Access Key authentication',
    authHint: 'When enabled, requests to this MCP service must carry a valid access key',
    optionAudit: 'Enable audit logging',
    auditHint: 'When enabled, an access log is recorded for every MCP service call',
    confirm: 'Confirm publish',
    confirmTitle: 'Confirm publish settings',
    confirmDesc: 'The following route configuration will be pushed to the API gateway:',
    progress: 'Publishing...'
  },
  preview: {
    serviceName: 'Service name',
    gatewayPath: 'Gateway path',
    backendAddress: 'Backend address',
    authStatus: 'Auth status',
    authEnabled: 'Access Key authentication enabled',
    auditLog: 'Audit logging',
    enabled: 'Enabled',
    disabled: 'Not enabled'
  },
  maintenance: {
    button: 'Maintenance',
    title: 'Set server maintenance',
    desc: 'Select the servers to put into maintenance (move to the right), or move servers currently in maintenance back to the left to cancel maintenance:',
    underMaintenance: 'Under maintenance',
    emptyUnderMaintenance: 'No servers under maintenance',
    confirm: 'Confirm maintenance',
    confirmCancel: 'Confirm cancel maintenance',
    confirmBoth: 'Confirm maintenance configuration',
    confirmDesc: 'The following maintenance configuration will be pushed to the API gateway; routes under maintenance will directly return 503:',
    sectionEnter: 'Enter maintenance',
    sectionRestore: 'Cancel maintenance',
    submitting: 'Submitting...'
  },
  alert: {
    deleteConfirm: 'Confirm delete',
    deleteBody: 'Are you sure you want to delete the MCP server "{name}"?',
    deleteWarn: 'Note: this only removes the record from the database and does not take the server offline from the API gateway automatically. Please manually delete the related route and upstream rules in the API gateway.',
    deleting: 'Deleting...',
    deleteSuccess: 'Server deleted. Please manually remove the route and upstream rules from the API gateway.',
    publishSuccess: 'Published successfully',
    publishFailed: 'Publish failed',
    maintenanceSuccess: 'Maintenance settings applied',
    operationFailed: 'Operation failed'
  }
}
