package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"
)

// Config 结构体定义配置文件结构
type Config struct {
	Server struct {
		Listen  string `yaml:"listen"`
		Timeout int    `yaml:"timeout"`
	} `yaml:"server"`
	DoH struct {
		Servers []struct {
			URL  string `yaml:"url"`
			Name string `yaml:"name"`
		} `yaml:"servers"`
		Timeout  int  `yaml:"timeout"`
		UseHTTP2 bool `yaml:"use_http2"`
	} `yaml:"doh"`
	TLS struct {
		Enabled            bool     `yaml:"enabled"`
		PrintCertInfo      bool     `yaml:"print_cert_info"`
		PinCertIssuer      bool     `yaml:"pin_cert_issuer"`
		AllowedIssuers     []string `yaml:"allowed_issuers"`
		InsecureSkipVerify bool     `yaml:"insecure_skip_verify"`
	} `yaml:"tls"`
	Logging struct {
		Level    string `yaml:"level"`
		QueryLog struct {
			Enabled bool   `yaml:"enabled"`
			Target  string `yaml:"target"`
			File    struct {
				Format     string `yaml:"format"`
				Path       string `yaml:"path"`
				MaxSize    int    `yaml:"max_size"`
				MaxBackups int    `yaml:"max_backups"`
				MaxAge     int    `yaml:"max_age"`
			} `yaml:"file"`
			Database struct {
				Type   string `yaml:"type"`
				SQLite struct {
					Path string `yaml:"path"`
				} `yaml:"sqlite"`
				PostgreSQL struct {
					Host     string `yaml:"host"`
					Port     int    `yaml:"port"`
					Database string `yaml:"database"`
					User     string `yaml:"user"`
					Password string `yaml:"password"`
					SSLMode  string `yaml:"sslmode"`
				} `yaml:"postgresql"`
			} `yaml:"database"`
		} `yaml:"query_log"`
	} `yaml:"logging"`
}

var (
	config     Config
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yaml", "Path to config file")
}

// loadConfig 加载配置文件
func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil
}

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config = *cfg

	// 设置日志级别
	log.Printf("DNS to DoH converter starting...")
	log.Printf("Listen address: %s", config.Server.Listen)
	log.Printf("Configured %d DoH servers:", len(config.DoH.Servers))
	for i, server := range config.DoH.Servers {
		log.Printf("  [%d] %s - %s", i+1, server.Name, server.URL)
	}

	// 初始化 TLS 配置管理器
	tlsManager := NewTLSConfigManager(&config)
	if err := tlsManager.ValidateTLSConfig(); err != nil {
		log.Fatalf("TLS configuration validation failed: %v", err)
	}

	// 初始化查询日志记录器
	queryLogger, err := NewQueryLogger(&config)
	if err != nil {
		log.Fatalf("Failed to initialize query logger: %v", err)
	}
	defer queryLogger.Close()

	// 初始化 DoH 客户端
	dohClient := NewDoHClient(&config, tlsManager)

	// 启动 DNS 服务器
	dnsServer := NewDNSServer(&config, dohClient, queryLogger)
	if err := dnsServer.Start(); err != nil {
		log.Fatalf("Failed to start DNS server: %v", err)
	}

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down server...")
	dnsServer.Stop()
	log.Println("Server stopped")
}
