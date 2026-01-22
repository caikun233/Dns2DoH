# DNS to DoH Converter

[English](#english) | [中文](#中文)

---

## English

A lightweight DNS to DoH (DNS over HTTPS) converter built with Golang. It accepts traditional UDP DNS queries and forwards them to DoH servers via HTTPS.

### Features

- ✅ Accept UDP DNS queries
- ✅ Forward queries via DoH (DNS over HTTPS) protocol
- ✅ Multiple DoH servers support (automatic failover)
- ✅ YAML configuration file
- ✅ Customizable listen address and port
- ✅ HTTP/2 support
- ✅ Detailed query logging (Console, File, SQLite, PostgreSQL)
- ✅ TLS certification verification control
- ✅ Support all DNS record types (A, AAAA, CNAME, MX, TXT, etc.)

### Installation

#### Prerequisites

- Go 1.21 or higher

#### Build from Source

```bash
# Clone the repository
git clone https://github.com/yourusername/Dns2DoH.git
cd Dns2DoH

# Download dependencies
go mod download

# Build (CGO not required for SQLite)
go build -o dns2doh.exe
```

### Configuration

Please refer to [CONFIGURATION.md](CONFIGURATION.md) for detailed configuration instructions.

### Usage

#### Start the Server

```bash
# Use default config file (config.yaml)
./dns2doh.exe

# Use custom config file
./dns2doh.exe -config /path/to/config.yaml
```

#### Test DNS Queries

Test using `nslookup` or `dig`:

```bash
# Windows (PowerShell)
nslookup google.com 127.0.0.1

# Linux/Mac
dig @127.0.0.1 google.com
```

#### Windows Service

To run without admin privileges on port 53, change the port to a value greater than 1024 (e.g., 5353):

```yaml
server:
  listen: "0.0.0.0:5353"
```

### Popular DoH Servers

Here are some available public DoH servers:

| Provider | DoH URL |
|----------|---------|
| Cloudflare | https://cloudflare-dns.com/dns-query |
| Google | https://dns.google/dns-query |
| Quad9 | https://dns.quad9.net/dns-query |
| Alibaba DNS | https://dns.alidns.com/dns-query |
| DNSPod | https://doh.pub/dns-query |
| 360 DNS | https://doh.360.cn/dns-query |

### Project Structure

```
Dns2DoH/
├── main.go          # Main program entry
├── dns_handler.go   # DNS request handling
├── doh_client.go    # DoH client implementation
├── config.yaml      # Configuration file
├── go.mod           # Go module dependencies
└── README.md        # Project documentation
```

### Log Example

```
2026/01/22 10:30:00 DNS to DoH converter starting...
2026/01/22 10:30:00 Listen address: 0.0.0.0:53
2026/01/22 10:30:00 Configured 3 DoH servers:
2026/01/22 10:30:00   [1] Cloudflare - https://cloudflare-dns.com/dns-query
2026/01/22 10:30:00   [2] Google - https://dns.google/dns-query
2026/01/22 10:30:00   [3] AliDNS - https://dns.alidns.com/dns-query
2026/01/22 10:30:00 UDP DNS server listening on 0.0.0.0:53
2026/01/22 10:30:15 Query received: google.com. (type: A) from: 192.168.1.100:54321
2026/01/22 10:30:15 Query successful: google.com. -> 1 answers (elapsed: 45ms)
2026/01/22 10:30:15   A record: google.com. -> 142.250.185.46 (TTL: 300)
```

### Tech Stack

- **Language**: Go 1.21+
- **DNS Library**: github.com/miekg/dns
- **YAML Parser**: gopkg.in/yaml.v3
- **HTTP/2**: golang.org/x/net/http2

### License

See [LICENSE](LICENSE) file for details.

### Contributing

Issues and Pull Requests are welcome!

### FAQ

#### Q: How to run with admin privileges on Windows?

A: Right-click the program and select "Run as administrator", or change the listen port to a value greater than 1024.

#### Q: Does it support IPv6?

A: Yes, the program fully supports IPv6 DNS queries (AAAA records).

#### Q: How to add custom DoH servers?

A: Add new server entries in the `doh.servers` section of `config.yaml`.

#### Q: What are the resource requirements?

A: Very lightweight, typically 10-20MB memory usage with minimal CPU usage.

---

## 中文

一个基于 Golang 的 DNS 到 DoH (DNS over HTTPS) 转换器，可以接收传统的 UDP DNS 查询请求，并通过 DoH 协议转发到上游 DNS 服务器。

### 功能特性

- ✅ 接收 UDP DNS 查询请求
- ✅ 通过 DoH (DNS over HTTPS) 协议转发查询
- ✅ 支持多个 DoH 服务器（自动故障转移）
- ✅ YAML 配置文件支持
- ✅ 可自定义监听地址和端口
- ✅ HTTP/2 支持
- ✅ 详细的查询日志（支持控制台、文件、SQLite、PostgreSQL）
- ✅ TLS 证书校验控制
- ✅ 支持所有 DNS 记录类型（A, AAAA, CNAME, MX, TXT 等）

### 安装

#### 前置要求

- Go 1.21 或更高版本

#### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/yourusername/Dns2DoH.git
cd Dns2DoH

# 安装依赖
go mod download

# 编译 (SQLite 无需 CGO 环境)
go build -o dns2doh.exe
```

### 配置

请参考 [CONFIGURATION.md](CONFIGURATION.md) 查看详细的配置说明。

### 使用方法

#### 启动服务器

```bash
# 使用默认配置文件（config.yaml）
./dns2doh.exe

# 使用自定义配置文件
./dns2doh.exe -config /path/to/config.yaml
```

#### 测试 DNS 查询

使用 `nslookup` 或 `dig` 工具测试：

```bash
# Windows (PowerShell)
nslookup google.com 127.0.0.1

# Linux/Mac
dig @127.0.0.1 google.com
```

#### Windows 服务运行

如果需要在非管理员权限下监听 53 端口，可以修改配置文件中的端口为大于 1024 的端口（如 5353）：

```yaml
server:
  listen: "0.0.0.0:5353"
```

### 常见 DoH 服务器

以下是一些可用的公共 DoH 服务器：

| 提供商 | DoH URL |
|--------|---------|
| Cloudflare | https://cloudflare-dns.com/dns-query |
| Google | https://dns.google/dns-query |
| Quad9 | https://dns.quad9.net/dns-query |
| 阿里 DNS | https://dns.alidns.com/dns-query |
| DNSPod | https://doh.pub/dns-query |
| 360 DNS | https://doh.360.cn/dns-query |

### 项目结构

```
Dns2DoH/
├── main.go          # 主程序入口
├── dns_handler.go   # DNS 请求处理
├── doh_client.go    # DoH 客户端实现
├── config.yaml      # 配置文件
├── go.mod           # Go 模块依赖
└── README.md        # 项目说明
```

### 日志示例

```
2026/01/22 10:30:00 DNS to DoH 转换器启动中...
2026/01/22 10:30:00 监听地址: 0.0.0.0:53
2026/01/22 10:30:00 配置了 3 个 DoH 服务器:
2026/01/22 10:30:00   [1] Cloudflare - https://cloudflare-dns.com/dns-query
2026/01/22 10:30:00   [2] Google - https://dns.google/dns-query
2026/01/22 10:30:00   [3] AliDNS - https://dns.alidns.com/dns-query
2026/01/22 10:30:00 UDP DNS 服务器正在监听 0.0.0.0:53
2026/01/22 10:30:15 收到查询: google.com. (类型: A) 来自: 192.168.1.100:54321
2026/01/22 10:30:15 查询成功: google.com. -> 1 条应答 (耗时: 45ms)
2026/01/22 10:30:15   A 记录: google.com. -> 142.250.185.46 (TTL: 300)
```

### 技术栈

- **语言**: Go 1.21+
- **DNS 库**: github.com/miekg/dns
- **YAML 解析**: gopkg.in/yaml.v3
- **HTTP/2**: golang.org/x/net/http2

### 许可证

查看 [LICENSE](LICENSE) 文件了解详情。

### 贡献

欢迎提交 Issue 和 Pull Request！

### 常见问题

#### Q: 如何在 Windows 上以管理员权限运行？

A: 右键点击程序，选择"以管理员身份运行"，或者修改监听端口为大于 1024 的端口。

#### Q: 支持 IPv6 吗？

A: 是的，程序完全支持 IPv6 DNS 查询（AAAA 记录）。

#### Q: 如何添加自定义 DoH 服务器？

A: 在 `config.yaml` 的 `doh.servers` 部分添加新的服务器条目即可。

#### Q: 程序占用多少资源？

A: 非常轻量，内存占用通常在 10-20MB 左右，CPU 占用极低。
