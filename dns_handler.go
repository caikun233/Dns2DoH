package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

// DNSServer DNS 服务器结构体
type DNSServer struct {
	config      *Config
	dohClient   *DoHClient
	server      *dns.Server
	queryLogger QueryLogger
}

// NewDNSServer 创建新的 DNS 服务器实例
func NewDNSServer(config *Config, dohClient *DoHClient, queryLogger QueryLogger) *DNSServer {
	return &DNSServer{
		config:      config,
		dohClient:   dohClient,
		queryLogger: queryLogger,
	}
}

// Start 启动 DNS 服务器
func (s *DNSServer) Start() error {
	// 创建 UDP 服务器
	s.server = &dns.Server{
		Addr: s.config.Server.Listen,
		Net:  "udp",
	}

	// 设置 DNS 查询处理器
	dns.HandleFunc(".", s.handleDNSRequest)

	// 启动服务器
	log.Printf("UDP DNS server listening on %s", s.config.Server.Listen)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Fatalf("DNS server error: %v", err)
		}
	}()

	return nil
}

// Stop 停止 DNS 服务器
func (s *DNSServer) Stop() error {
	if s.server != nil {
		return s.server.Shutdown()
	}
	return nil
}

// handleDNSRequest 处理 DNS 查询请求
func (s *DNSServer) handleDNSRequest(w dns.ResponseWriter, req *dns.Msg) {
	startTime := time.Now()
	clientAddr := w.RemoteAddr().String()

	var domain string
	var queryType string
	if len(req.Question) > 0 {
		domain = req.Question[0].Name
		queryType = dns.TypeToString[req.Question[0].Qtype]
	}

	// 创建响应消息
	resp := new(dns.Msg)
	resp.SetReply(req)
	resp.Compress = false

	// 检查是否有查询问题
	if len(req.Question) == 0 {
		resp.SetRcode(req, dns.RcodeFormatError)
		w.WriteMsg(resp)
		return
	}

	// 通过 DoH 查询 DNS
	dohResp, dohServer, err := s.dohClient.QueryWithServer(req)
	queryDuration := time.Since(startTime)

	if err != nil {
		log.Printf("DoH query failed: %v", err)
		resp.SetRcode(req, dns.RcodeServerFailure)
		w.WriteMsg(resp)

		// Log failed query
		s.queryLogger.Log(QueryLogEntry{
			Timestamp:    startTime,
			ClientIP:     clientAddr,
			Domain:       domain,
			QueryType:    queryType,
			ResponseCode: "SERVFAIL",
			AnswerCount:  0,
			Duration:     queryDuration.Milliseconds(),
			DoHServer:    dohServer,
		})
		return
	}

	// Extract answer details
	answers := make([]AnswerEntry, 0, len(dohResp.Answer))
	for _, ans := range dohResp.Answer {
		answer := AnswerEntry{
			Name: ans.Header().Name,
			Type: dns.TypeToString[ans.Header().Rrtype],
			TTL:  ans.Header().Ttl,
		}

		switch rr := ans.(type) {
		case *dns.A:
			answer.Value = rr.A.String()
		case *dns.AAAA:
			answer.Value = rr.AAAA.String()
		case *dns.CNAME:
			answer.Value = rr.Target
		case *dns.MX:
			answer.Value = fmt.Sprintf("%s (priority: %d)", rr.Mx, rr.Preference)
		case *dns.TXT:
			answer.Value = fmt.Sprintf("%v", rr.Txt)
		case *dns.NS:
			answer.Value = rr.Ns
		case *dns.PTR:
			answer.Value = rr.Ptr
		case *dns.SOA:
			answer.Value = fmt.Sprintf("ns: %s, mbox: %s", rr.Ns, rr.Mbox)
		default:
			answer.Value = ans.String()
		}

		answers = append(answers, answer)
	}

	// Log successful query with answers
	s.queryLogger.Log(QueryLogEntry{
		Timestamp:    startTime,
		ClientIP:     clientAddr,
		Domain:       domain,
		QueryType:    queryType,
		ResponseCode: dns.RcodeToString[dohResp.Rcode],
		AnswerCount:  len(dohResp.Answer),
		Answers:      answers,
		Duration:     queryDuration.Milliseconds(),
		DoHServer:    dohServer,
	})

	// Print detailed answer records if enabled
	if s.config.Logging.QueryLog.Enabled && s.config.Logging.Level == "debug" {
		for _, ans := range dohResp.Answer {
			switch rr := ans.(type) {
			case *dns.A:
				log.Printf("  A record: %s -> %s (TTL: %d)", rr.Hdr.Name, rr.A.String(), rr.Hdr.Ttl)
			case *dns.AAAA:
				log.Printf("  AAAA record: %s -> %s (TTL: %d)", rr.Hdr.Name, rr.AAAA.String(), rr.Hdr.Ttl)
			case *dns.CNAME:
				log.Printf("  CNAME record: %s -> %s (TTL: %d)", rr.Hdr.Name, rr.Target, rr.Hdr.Ttl)
			case *dns.MX:
				log.Printf("  MX record: %s -> %s (priority: %d, TTL: %d)", rr.Hdr.Name, rr.Mx, rr.Preference, rr.Hdr.Ttl)
			case *dns.TXT:
				log.Printf("  TXT record: %s -> %v (TTL: %d)", rr.Hdr.Name, rr.Txt, rr.Hdr.Ttl)
			default:
				log.Printf("  Other record: %s (type: %s)", ans.Header().Name, dns.TypeToString[ans.Header().Rrtype])
			}
		}
	}

	// 发送响应
	if err := w.WriteMsg(dohResp); err != nil {
		log.Printf("Failed to send response: %v", err)
	}
}

// validateIPAddress 验证 IP 地址格式
func validateIPAddress(addr string) error {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("invalid address format: %v", err)
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return fmt.Errorf("invalid IP address: %s", host)
	}

	return nil
}
