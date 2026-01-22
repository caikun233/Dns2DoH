package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/miekg/dns"
	"golang.org/x/net/http2"
)

// DoHClient DoH 客户端结构体
type DoHClient struct {
	config     *Config
	httpClient *http.Client
	tlsManager *TLSConfigManager
}

// NewDoHClient 创建新的 DoH 客户端
func NewDoHClient(config *Config, tlsManager *TLSConfigManager) *DoHClient {
	// 创建 HTTP 客户端
	transport := &http.Transport{
		TLSClientConfig:     tlsManager.GetTLSConfig(),
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}

	// 如果配置启用 HTTP/2
	if config.DoH.UseHTTP2 {
		if err := http2.ConfigureTransport(transport); err != nil {
			log.Printf("Failed to configure HTTP/2: %v, using HTTP/1.1", err)
		}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.DoH.Timeout) * time.Second,
	}

	return &DoHClient{
		config:     config,
		httpClient: httpClient,
		tlsManager: tlsManager,
	}
}

// Query 通过 DoH 查询 DNS
func (c *DoHClient) Query(req *dns.Msg) (*dns.Msg, error) {
	resp, _, err := c.QueryWithServer(req)
	return resp, err
}

// QueryWithServer 通过 DoH 查询 DNS 并返回使用的服务器
func (c *DoHClient) QueryWithServer(req *dns.Msg) (*dns.Msg, string, error) {
	// 将 DNS 消息打包为字节
	packed, err := req.Pack()
	if err != nil {
		return nil, "", fmt.Errorf("failed to pack DNS message: %v", err)
	}

	// 尝试每个配置的 DoH 服务器
	var lastErr error
	for _, server := range c.config.DoH.Servers {
		resp, err := c.queryServer(server.URL, packed)
		if err != nil {
			lastErr = err
			if c.config.Logging.Level == "debug" {
				log.Printf("DoH server %s (%s) query failed: %v", server.Name, server.URL, err)
			}
			continue
		}
		return resp, server.Name, nil
	}

	// 所有服务器都失败
	if lastErr != nil {
		return nil, "", fmt.Errorf("all DoH servers failed, last error: %v", lastErr)
	}
	return nil, "", fmt.Errorf("no available DoH servers")
}

// queryServer 向指定的 DoH 服务器发送查询
func (c *DoHClient) queryServer(serverURL string, packed []byte) (*dns.Msg, error) {
	// 创建 HTTP POST 请求
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.DoH.Timeout)*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, "POST", serverURL, bytes.NewReader(packed))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// 设置 DoH 请求头
	httpReq.Header.Set("Content-Type", "application/dns-message")
	httpReq.Header.Set("Accept", "application/dns-message")
	httpReq.Header.Set("User-Agent", "Dns2DoH/1.0")

	// 发送请求
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send DoH request: %v", err)
	}
	defer httpResp.Body.Close()

	// 检查 HTTP 状态码
	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DoH server returned error status code: %d", httpResp.StatusCode)
	}

	// 检查内容类型
	contentType := httpResp.Header.Get("Content-Type")
	if contentType != "application/dns-message" {
		return nil, fmt.Errorf("DoH server returned invalid content type: %s", contentType)
	}

	// 读取响应体
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read DoH response: %v", err)
	}

	// 解析 DNS 响应
	resp := new(dns.Msg)
	if err := resp.Unpack(body); err != nil {
		return nil, fmt.Errorf("failed to parse DNS response: %v", err)
	}

	return resp, nil
}

// Close 关闭 DoH 客户端
func (c *DoHClient) Close() {
	c.httpClient.CloseIdleConnections()
}
