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
  title: 'API Gateway Servers',
  desc: 'Manage the API gateway instances used by the platform. MCP servers are published to the default gateway.',
  add: 'Add gateway',
  refresh: 'Refresh status',
  loading: 'Loading…',
  refreshing: 'Refreshing…',
  empty: 'No gateway servers yet',
  emptyHint: 'Add an API gateway instance to start publishing MCP servers to it.',
  defaultTag: 'Default',
  provider: {
    apisix: 'Apache APISIX',
    kong: 'Kong'
  },
  status: {
    online: 'Online',
    offline: 'Offline',
    unknown: 'Not checked'
  },
  card: {
    adminUrl: 'Admin API',
    accessKeyHeader: 'Access Key Header',
    authGrpcAddrs: 'Auth gRPC addresses',
    auditGrpcAddrs: 'Audit gRPC addresses',
    noGrpc: 'Not configured',
    latency: '{ms} ms',
    noDescription: 'No description'
  },
  actions: {
    test: 'Test connection',
    testing: 'Testing…',
    setDefault: 'Set as default',
    edit: 'Edit'
  },
  form: {
    createTitle: 'Add Gateway Server',
    editTitle: 'Edit Gateway Server',
    name: 'Gateway name',
    namePlaceholder: 'e.g. Primary gateway',
    provider: 'Gateway type',
    adminUrl: 'Admin API URL',
    adminUrlPlaceholder: 'http://127.0.0.1:9180',
    adminKey: 'Admin API Key',
    adminKeyPlaceholder: 'Enter the Admin API Key',
    adminKeyKeep: 'Leave blank to keep the current key',
    accessKeyHeader: 'Access Key Header',
    accessKeyHeaderPlaceholder: 'X-Access-Key',
    authGrpcAddrs: 'Auth gRPC addresses',
    grpcPlaceholder: '127.0.0.1:9090',
    addAddress: 'Add address',
    grpcHint: 'The gRPC listen address of accesskey-auth-server, one per line; it is pushed to the accesskey_verify plugin on publish.',
    auditGrpcAddrs: 'Audit log gRPC addresses',
    auditGrpcPlaceholder: '127.0.0.1:9091',
    auditGrpcHint: 'The gRPC listen address of audit-log-server, one per line; it is pushed to the audit_log plugin on publish.',
    description: 'Description',
    descriptionPlaceholder: 'Optional notes about this gateway',
    providerHintApisix: 'Enter the base URL of the APISIX Admin API, e.g. http://127.0.0.1:9180 (without /apisix/admin).',
    providerHintKong: 'Enter the base URL of the Kong Admin API, e.g. http://127.0.0.1:8001; the Admin Key can be left blank when RBAC is not enabled.',
    test: 'Test connection',
    testing: 'Testing…',
    testOk: 'Connected successfully',
    testFail: 'Connection failed'
  },
  delete: {
    title: 'Delete Gateway Server',
    desc: 'Delete gateway "{name}"? This cannot be undone, and published routes will not be removed from the gateway automatically.',
    confirm: 'Delete',
    deleting: 'Deleting…'
  },
  alert: {
    createSuccess: 'Gateway server created',
    updateSuccess: 'Gateway server updated',
    deleteSuccess: 'Gateway server deleted',
    setDefaultSuccess: 'Default gateway updated',
    loadFailed: 'Failed to load gateways',
    saveFailed: 'Failed to save',
    deleteFailed: 'Failed to delete'
  }
}
