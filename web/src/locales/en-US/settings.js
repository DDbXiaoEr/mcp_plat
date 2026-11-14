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
    ops: 'Operations',
    operation: 'Management'
  },
  general: {
    saved: 'Saved',
    saveFailed: 'Save failed',
    remove: 'Remove',
    gotIt: 'Got it',
    addAddress: '+ Add address'
  },
  log: {
    title: 'Log Settings',
    enableSyslog: 'Enable Syslog',
    offNotice: 'When Syslog is disabled, logs are written to both standard output and a local file (with rotation).',
    onNotice: 'When Syslog is enabled, standard output and local file logging are disabled and logs are only sent to the Syslog server.',
    pathLabel: 'Log directory',
    levelLabel: 'Log level',
    fileNameLabel: 'Log file name',
    fileNameHelp: 'Log files are named \\{prefix\\}.log; the default is mcp_plat.log',
    rotation: 'Log rotation',
    maxSizeLabel: 'Max single file size (MB)',
    maxBackupsLabel: 'Number of history files kept',
    maxAgeLabel: 'Days to keep',
    compressLabel: 'Compress history logs',
    compressNone: 'No compression',
    compressGzip: 'gzip compression',
    syslogHostLabel: 'Syslog server address',
    syslogHostPlaceholder: 'Enter the Syslog server address',
    portLabel: 'Port',
    protocolLabel: 'Protocol'
  },
  mail: {
    title: 'SMTP Settings',
    enable: 'Enable email notifications',
    toggleHint: 'The toggle controls whether notification emails are actually sent; the fields below are always editable for testing and pre-configuration',
    hostLabel: 'SMTP server address',
    portLabel: 'Port',
    encryptionLabel: 'Encryption',
    encryptionNone: 'None',
    usernameLabel: 'Username',
    usernamePlaceholder: 'Enter the SMTP username',
    passwordLabel: 'Password',
    passwordPlaceholder: 'Enter the SMTP password',
    fromAddressLabel: 'From address',
    fromNameLabel: 'From name',
    fromNamePlaceholder: 'MCP Service Platform',
    sendTest: 'Send test email',
    sending: 'Sending…',
    testToPlaceholder: 'Enter a recipient email to verify the SMTP configuration',
    testOk: 'Test email sent successfully. Please check the inbox.',
    sendFailed: 'Failed to send'
  },
  network: {
    title: 'Network Security',
    allowLabel: 'Allowed intranet CIDR',
    tip1: 'One CIDR block per line, e.g. ', 
    tip2: '.',
    tip3: 'Once configured, connections to these intranet addresses will be allowed when fetching the tool list.',
    tip4: 'If you need to reach an MCP server\'s intranet IP, add the matching CIDR block here.'
  },
  platform: {
    title: 'Platform Settings',
    nameLabel: 'Platform name',
    namePlaceholder: 'Enter the platform name, e.g. the university name',
    logoLabel: 'Logo image URL',
    logoPlaceholder: 'Enter the URL of the logo image',
    loginBgLabel: 'Login page background image',
    loginBgPlaceholder: 'Enter the URL of the login page background image',
    loginBgHint: 'Leave blank to use the default gradient background; when filled, the login page fills and displays this image responsively',
    siteUrlLabel: 'Site link',
    siteUrlPlaceholder: 'URL to open when the Logo is clicked, e.g. https://www.example.edu.cn'
  },
  publishDomain: {
    title: 'Publish Domains',
    domainsLabel: 'MCP publish domains',
    domainPlaceholder: 'e.g. mcp.xauat.edu.cn',
    hint: 'Multiple publish domains are supported and can be selected when publishing an MCP service; enter domains without protocol or path.'
  },
  auth: {
    title: 'User Authentication',
    methodLabel: 'Authentication method',
    cas: {
      serverUrlLabel: 'CAS server URL',
      serviceUrlLabel: 'Service callback URL',
      versionLabel: 'Protocol version',
      mappingTitle: 'Attribute mapping (platform field → CAS attribute)',
      selectPlatformField: 'Select a platform field',
      attrPlaceholder: 'CAS attribute name',
      addMapping: '+ Add mapping'
    },
    ldap: {
      enabledNotice: 'LDAP authentication is enabled. For security reasons, the currently configured values are not shown; leaving a field blank keeps the existing configuration unchanged.',
      hostLabel: 'Server address',
      portLabel: 'Port',
      bindPasswordLabel: 'Bind password',
      bindPasswordPlaceholder: 'Enter the Bind password',
      userFilterLabel: 'User filter',
      mappingTitle: 'Attribute mapping (platform field → LDAP attribute)',
      selectPlatformField: 'Select a platform field',
      ldapAttrPlaceholder: 'LDAP attribute name',
      addMapping: '+ Add mapping',
      testTitle: 'Test mapping',
      testing: 'Testing…',
      testPlaceholder: 'Enter a username to test the mapping',
      testFailed: 'Test failed',
      allAttributes: 'All LDAP attributes ({count})',
      currentMapping: 'Current mapping result',
      noMapping: 'No mapping configured',
      saveHint1: 'Make sure you have clicked the "Save" button above to save the authentication settings. After saving, users need to',
      saveHint2: 'log in again',
      saveHint3: ' for their LDAP attributes to be synced to their profile.'
    },
    oauth: {
      authorizeUrlLabel: 'Authorization URL',
      tokenUrlLabel: 'Token URL',
      userinfoUrlLabel: 'User info URL',
      clientIdPlaceholder: 'Enter the Client ID',
      clientSecretPlaceholder: 'Enter the Client Secret',
      redirectUrlLabel: 'Callback URL'
    }
  },
  userOps: {
    title: 'User Management',
    maxAccessKeysLabel: 'Max AccessKeys per user',
    accessKeyCronLabel: 'Expiration scan cron expression',
    accessKeyCronHint: 'Periodically scans for expired AccessKeys and disables them; standard crontab format is supported',
    roleRuleTitle: 'Auto role assignment rules for new users',
    filterAttrLabel: 'Filter attribute',
    filterAttrHint: 'All role conditions share the same filter attribute; a user is assigned the role of the first rule their attribute value matches',
    conditionIndex: 'Condition {n}',
    satisfies: 'matches',
    patternPlaceholder: 'Regular expression, e.g. ^\\d{7}$',
    assignRole: 'assign role',
    selectRole: 'Select a role',
    removeCondition: 'Remove condition',
    addCondition: '+ Add condition'
  },
  fields: {
    username: 'Username',
    uid: 'UID',
    name: 'Name',
    email: 'Email',
    phone: 'Phone',
    organization: 'Organization'
  },
  quickAccess: {
    supportedClients: 'Supported clients',
    supportedClientsHint: 'Controls which clients are available on the user\'s "Quick Access" page; checked clients are exposed to users.',
    schemeLabel: 'Access protocol',
    schemeHint: 'Determines the protocol (HTTP or HTTPS) used to generate access URLs.'
  }
}
