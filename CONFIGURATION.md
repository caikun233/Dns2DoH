# Configuration / 配置文档

[English](#english) | [中文](#中文)

---

## English

### Configuration File

The server is configured using a YAML file (default: `config.yaml`). You can use the `-config` flag to specify a custom configuration file path.

See `config.example.yaml` in the repository for a complete example with all available options.

### Basic Configuration Example

```yaml
# UDP DNS server configuration
server:
  # Listen address (default: 0.0.0.0:53)
  listen: "0.0.0.0:53"
  # Timeout in seconds
  timeout: 5

# DoH servers list
doh:
  # DoH server list (will try in order)
  servers:
    - url: "https://dns.alidns.com/dns-query"
      name: "AliDNS"
    - url: "https://cloudflare-dns.com/dns-query"
      name: "Cloudflare"
    - url: "https://dns.google/dns-query"
      name: "Google"
  
  # DoH request timeout in seconds
  timeout: 10
  
  # Enable HTTP/2
  use_http2: true

# TLS Configuration
tls:
  # Print detailed certificate information for each connection
  print_cert_info: false
  # Allow insecure connections (skip certificate verification) - NOT RECOMMENDED
  insecure_skip_verify: false
  # Certificate Pinning (optional) - Check if issuer matches
  # pin_issuer_common_name: "GlobalSign"

# Logging configuration
logging:
  level: "info"
  
  query_log:
    enabled: true
    # Log target: console, file, database
    target: "console"
    
    # File logging settings
    file:
      path: "logs/queries.log"
      format: "json" # json or csv
      max_size_mb: 10
      max_backups: 5
      max_age_days: 30
      compress: true
      
    # Database logging settings
    database:
      # sqlite or postgresql
      type: "sqlite"
      sqlite:
        path: "logs/queries.db"
      postgresql:
        host: "localhost"
        port: 5432
        user: "postgres"
        password: "password"
        database: "dnslogs"
        ssl_mode: "disable"
```

---

## 中文

### 配置文件

服务器使用 YAML 文件进行配置（默认为 `config.yaml`）。您可以使用 `-config` 参数指定自定义配置文件路径。

完整的配置选项请参考仓库中的 `config.example.yaml` 文件。

### 基础配置示例

```yaml
# UDP DNS 服务器配置
server:
  # 监听地址（默认：0.0.0.0:53）
  listen: "0.0.0.0:53"
  # 超时时间（秒）
  timeout: 5

# DoH 服务器配置
doh:
  # DoH 服务器列表（会按顺序尝试）
  servers:
    - url: "https://dns.alidns.com/dns-query"
      name: "AliDNS"
    - url: "https://cloudflare-dns.com/dns-query"
      name: "Cloudflare"
    - url: "https://dns.google/dns-query"
      name: "Google"
  
  # DoH 请求超时时间（秒）
  timeout: 10
  
  # 是否使用 HTTP/2
  use_http2: true

# TLS 配置
tls:
  # 打印详细的证书信息
  print_cert_info: false
  # 允许不安全的连接（跳过证书校验）- 不推荐
  insecure_skip_verify: false
  # 证书锁定（可选）- 校验签发者名称
  # pin_issuer_common_name: "GlobalSign"

# 日志配置
logging:
  # 日志级别: debug, info, warn, error
  level: "info"
  
  query_log:
    enabled: true
    # 日志目标: console (控制台), file (文件), database (数据库)
    target: "console"
    
    # 文件日志设置
    file:
      path: "logs/queries.log"
      format: "json" # json 或 csv
      max_size_mb: 10
      max_backups: 5
      max_age_days: 30
      compress: true
      
    # 数据库日志设置
    database:
      # sqlite 或 postgresql
      type: "sqlite"
      sqlite:
        path: "logs/queries.db"
      postgresql:
        host: "localhost"
        port: 5432
        user: "postgres"
        password: "password"
        database: "dnslogs"
        ssl_mode: "disable"
```
