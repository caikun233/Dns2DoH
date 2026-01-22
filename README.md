# Dns2DoH

[English](#english) | [中文](#中文)

---

## English

### Overview

Dns2DoH is a lightweight DNS to DNS-over-HTTPS (DoH) converter that acts as a bridge between traditional UDP DNS clients and modern DoH servers. It accepts standard UDP DNS queries on port 53 and forwards them to DoH servers over HTTPS, providing enhanced privacy and security for DNS resolution.

### Features

- **UDP DNS Server**: Accepts traditional DNS queries over UDP protocol
- **DoH Client**: Forwards queries to DNS-over-HTTPS servers
- **Transparent Conversion**: Seamlessly converts between UDP DNS and DoH protocols
- **Privacy & Security**: Encrypts DNS queries using HTTPS
- **Easy Configuration**: Simple setup with customizable DoH server endpoints
- **Lightweight**: Minimal resource usage and fast response times
- **Cross-Platform**: Works on Linux, Windows, and macOS

### How It Works

```
Client (UDP DNS Query) → Dns2DoH → DoH Server (HTTPS) → DoH Server Response → Dns2DoH → Client (UDP DNS Response)
```

1. Client sends a standard UDP DNS query to Dns2DoH (usually on port 53)
2. Dns2DoH receives the query and converts it to DoH format
3. The query is sent to a configured DoH server over HTTPS
4. DoH server processes the query and returns the response
5. Dns2DoH converts the response back to UDP DNS format
6. Client receives the DNS response as if from a traditional DNS server

### Installation

#### Prerequisites

- [To be added based on implementation language]

#### From Source

```bash
git clone https://github.com/caikun233/Dns2DoH.git
cd Dns2DoH
# Build instructions to be added
```

### Usage

#### Basic Usage

```bash
# Start the DNS to DoH converter
# Command to be added based on implementation
```

#### Configuration

The converter can be configured using a configuration file or command-line arguments:

- **Listen Address**: The local address and port to listen for DNS queries (default: `0.0.0.0:53`)
- **DoH Server**: The DNS-over-HTTPS server endpoint (e.g., `https://dns.google/dns-query`)
- **Timeout**: Request timeout in seconds
- **Cache**: Enable/disable DNS response caching

#### Example Configuration

```yaml
# Configuration example to be added
listen: "0.0.0.0:53"
doh_server: "https://dns.google/dns-query"
timeout: 5
cache_enabled: true
```

### Popular DoH Servers

- **Google Public DNS**: `https://dns.google/dns-query`
- **Cloudflare DNS**: `https://cloudflare-dns.com/dns-query`
- **Quad9**: `https://dns.quad9.net/dns-query`
- **AdGuard DNS**: `https://dns.adguard.com/dns-query`
- **OpenDNS**: `https://doh.opendns.com/dns-query`

### Use Cases

- **Privacy Enhancement**: Encrypt DNS queries to prevent ISP monitoring
- **Bypass DNS Filtering**: Access blocked content by using alternative DoH servers
- **Legacy System Support**: Enable DoH for devices that only support traditional DNS
- **Network Testing**: Test and compare different DoH providers
- **Development**: Integrate DoH support into applications without modifying them

### Performance Considerations

- DNS responses may be cached to improve performance
- HTTPS connection pooling reduces latency for subsequent queries
- Concurrent query handling for better throughput

### Security

- All DNS queries are encrypted using HTTPS/TLS
- Server certificate validation ensures connection security
- No logging of DNS queries (depending on DoH server policy)

### Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

### Acknowledgments

- DNS-over-HTTPS (DoH) specification: [RFC 8484](https://tools.ietf.org/html/rfc8484)
- Inspired by the need for privacy-preserving DNS resolution

---

## 中文

### 概述

Dns2DoH 是一个轻量级的 DNS 到 DNS-over-HTTPS (DoH) 转换器，充当传统 UDP DNS 客户端和现代 DoH 服务器之间的桥梁。它在 53 端口接受标准的 UDP DNS 查询，并通过 HTTPS 将它们转发到 DoH 服务器，为 DNS 解析提供增强的隐私和安全性。

### 功能特性

- **UDP DNS 服务器**：接受通过 UDP 协议的传统 DNS 查询
- **DoH 客户端**：将查询转发到 DNS-over-HTTPS 服务器
- **透明转换**：在 UDP DNS 和 DoH 协议之间无缝转换
- **隐私与安全**：使用 HTTPS 加密 DNS 查询
- **简易配置**：可自定义 DoH 服务器端点的简单设置
- **轻量级**：资源占用少，响应速度快
- **跨平台**：支持 Linux、Windows 和 macOS

### 工作原理

```
客户端 (UDP DNS 查询) → Dns2DoH → DoH 服务器 (HTTPS) → DoH 服务器响应 → Dns2DoH → 客户端 (UDP DNS 响应)
```

1. 客户端向 Dns2DoH 发送标准 UDP DNS 查询（通常在 53 端口）
2. Dns2DoH 接收查询并将其转换为 DoH 格式
3. 通过 HTTPS 将查询发送到配置的 DoH 服务器
4. DoH 服务器处理查询并返回响应
5. Dns2DoH 将响应转换回 UDP DNS 格式
6. 客户端接收 DNS 响应，就像来自传统 DNS 服务器一样

### 安装

#### 前置要求

- [根据实现语言添加]

#### 从源码安装

```bash
git clone https://github.com/caikun233/Dns2DoH.git
cd Dns2DoH
# 构建说明待添加
```

### 使用方法

#### 基本使用

```bash
# 启动 DNS 到 DoH 转换器
# 根据实现添加命令
```

#### 配置

转换器可以通过配置文件或命令行参数进行配置：

- **监听地址**：监听 DNS 查询的本地地址和端口（默认：`0.0.0.0:53`）
- **DoH 服务器**：DNS-over-HTTPS 服务器端点（例如：`https://dns.google/dns-query`）
- **超时时间**：请求超时时间（秒）
- **缓存**：启用/禁用 DNS 响应缓存

#### 配置示例

```yaml
# 配置示例待添加
listen: "0.0.0.0:53"
doh_server: "https://dns.google/dns-query"
timeout: 5
cache_enabled: true
```

### 常用 DoH 服务器

- **Google Public DNS**: `https://dns.google/dns-query`
- **Cloudflare DNS**: `https://cloudflare-dns.com/dns-query`
- **Quad9**: `https://dns.quad9.net/dns-query`
- **AdGuard DNS**: `https://dns.adguard.com/dns-query`
- **OpenDNS**: `https://doh.opendns.com/dns-query`

### 使用场景

- **隐私增强**：加密 DNS 查询以防止 ISP 监控
- **绕过 DNS 过滤**：通过使用替代 DoH 服务器访问被屏蔽的内容
- **传统系统支持**：为仅支持传统 DNS 的设备启用 DoH
- **网络测试**：测试和比较不同的 DoH 提供商
- **开发**：在不修改应用程序的情况下集成 DoH 支持

### 性能考虑

- DNS 响应可以被缓存以提高性能
- HTTPS 连接池减少后续查询的延迟
- 并发查询处理以获得更好的吞吐量

### 安全性

- 所有 DNS 查询都使用 HTTPS/TLS 加密
- 服务器证书验证确保连接安全
- 不记录 DNS 查询（取决于 DoH 服务器政策）

### 贡献

欢迎贡献！请随时提交 Pull Request。

### 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

### 致谢

- DNS-over-HTTPS (DoH) 规范：[RFC 8484](https://tools.ietf.org/html/rfc8484)
- 受隐私保护 DNS 解析需求的启发
