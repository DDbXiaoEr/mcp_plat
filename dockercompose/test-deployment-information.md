# 测试部署账户信息

## 管理员账户

| 服务 | 用户名 | 密码 |
|------|--------|------|
| MCP平台管理员 | admin | admin123 |

## 数据库账户

### PostgreSQL (主数据库)
| 配置项 | 值 |
|--------|-----|
| 主机 | mainpostgres |
| 端口 | 5432 |
| 用户名 | root |
| 密码 | 123456 |
| 数据库名 | mcp_plat_db |

### PostgreSQL (审计日志数据库)
| 配置项 | 值 |
|--------|-----|
| 主机 | auditpostgres |
| 端口 | 5432 |
| 用户名 | root |
| 密码 | 123456 |
| 数据库名 | mcp_history_db |

### ClickHouse
| 配置项 | 值 |
|--------|-----|
| 地址 | clickhouse:9000 |
| 用户名 | default |
| 密码 | (空) |
| 数据库名 | audit_logs |

## LDAP账户

### LDAP管理员
| 配置项 | 值 |
|--------|-----|
| 组织 | XXX University |
| 域名 | xxx.edu.cn |
| 管理员密码 | admin@123456 |
| 配置密码 | admin@123456 |
| Base DN | dc=xxx,dc=edu,dc=cn |

### 测试用户账户

#### 学生账户 (student01 - student20)
| 用户名 | 邮箱 | 密码 (SSHA) |
|--------|------|-------------|
| student01 | student01@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student02 | student02@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student03 | student03@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student04 | student04@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student05 | student05@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student06 | student06@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student07 | student07@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student08 | student08@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student09 | student09@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student10 | student10@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student11 | student11@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student12 | student12@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student13 | student13@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student14 | student14@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student15 | student15@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student16 | student16@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student17 | student17@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student18 | student18@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student19 | student19@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| student20 | student20@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |

#### 教职工账户 (teacher01 - teacher20)
| 用户名 | 邮箱 | 密码 (SSHA) |
|--------|------|-------------|
| teacher01 | teacher01@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher02 | teacher02@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher03 | teacher03@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher04 | teacher04@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher05 | teacher05@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher06 | teacher06@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher07 | teacher07@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher08 | teacher08@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher09 | teacher09@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher10 | teacher10@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher11 | teacher11@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher12 | teacher12@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher13 | teacher13@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher14 | teacher14@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher15 | teacher15@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher16 | teacher16@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher17 | teacher17@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher18 | teacher18@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher19 | teacher19@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |
| teacher20 | teacher20@xxx.edu.cn | {SSHA}+N64LmEZzCpa6JbdI19pR/dS128JnO5p |

## 密钥配置

| 密钥类型 | 值 |
|----------|-----|
| JWT Secret | 48554e69f60343595de5aae5430d0c3f90bd38f5cd513d88b16d972a7b23ad74 |
| Access Key Secret | 093e4f7d88cf7c8fa98244574dac873399ca3cf94f7128c1fb8f8f2789a8b132 |

## 服务端口

| 服务 | 端口 |
|------|------|
| MCP平台主服务 | 8080 |
| APISIX网关 | 80, 443, 9180 |
| AccessKey认证服务 | 9090 |
| 审计日志服务 | 9091 |
| PostgreSQL (主) | 5432 |
| PostgreSQL (审计) | 5433 |
| ClickHouse HTTP | 8123 |
| ClickHouse TCP | 9000 |
| ClickHouse 其他 | 9009 |
| etcd | 2379, 2380 |

## 镜像版本

| 环境 | TAG |
|------|-----|
| ClickHouse环境 | latest |
| PostgreSQL环境 | 64c736b |
| 根目录 | a04ec93 |
