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
import base from './base.js'
import login from './login.js'
import welcome from './welcome.js'
import profile from './profile.js'
import overview from './overview.js'
import history from './history.js'
import servers from './servers.js'
import accesskey from './accesskey.js'
import rbac from './rbac.js'
import settings from './settings.js'
import quickaccess from './quickaccess.js'

export default {
  ...base,
  login,
  welcome,
  profile,
  overview,
  history,
  servers,
  accesskey,
  rbac,
  settings,
  quickaccess
}
