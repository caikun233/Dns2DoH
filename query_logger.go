package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/lib/pq" // PostgreSQL driver
	"gopkg.in/natefinch/lumberjack.v2"
	_ "modernc.org/sqlite" // Pure Go SQLite driver (no CGO required)
)

// QueryLogEntry represents a DNS query log entry
type QueryLogEntry struct {
	Timestamp    time.Time     `json:"timestamp"`
	ClientIP     string        `json:"client_ip"`
	Domain       string        `json:"domain"`
	QueryType    string        `json:"query_type"`
	ResponseCode string        `json:"response_code"`
	AnswerCount  int           `json:"answer_count"`
	Answers      []AnswerEntry `json:"answers,omitempty"`
	Duration     int64         `json:"duration_ms"`
	DoHServer    string        `json:"doh_server"`
}

// AnswerEntry represents a single DNS answer record
type AnswerEntry struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	TTL   uint32 `json:"ttl"`
	Value string `json:"value"`
}

// QueryLogger interface for different logging backends
type QueryLogger interface {
	Log(entry QueryLogEntry) error
	Close() error
}

// ConsoleLogger logs to console
type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Log(entry QueryLogEntry) error {
	log.Printf("Query received: %s (type: %s) from: %s", entry.Domain, entry.QueryType, entry.ClientIP)
	if entry.AnswerCount > 0 {
		log.Printf("Query successful: %s -> %d answers (elapsed: %dms)", entry.Domain, entry.AnswerCount, entry.Duration)
	}
	return nil
}

func (l *ConsoleLogger) Close() error {
	return nil
}

// FileLogger logs to file (JSON or CSV)
type FileLogger struct {
	format     string
	file       *os.File
	jsonWriter *json.Encoder
	csvWriter  *csv.Writer
	logger     *lumberjack.Logger
	mu         sync.Mutex
}

func NewFileLogger(config *Config) (*FileLogger, error) {
	logDir := filepath.Dir(config.Logging.QueryLog.File.Path)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %v", err)
	}

	logger := &lumberjack.Logger{
		Filename:   config.Logging.QueryLog.File.Path,
		MaxSize:    config.Logging.QueryLog.File.MaxSize,
		MaxBackups: config.Logging.QueryLog.File.MaxBackups,
		MaxAge:     config.Logging.QueryLog.File.MaxAge,
		Compress:   true,
	}

	fl := &FileLogger{
		format: config.Logging.QueryLog.File.Format,
		logger: logger,
	}

	if fl.format == "json" {
		fl.jsonWriter = json.NewEncoder(logger)
	} else if fl.format == "csv" {
		fl.csvWriter = csv.NewWriter(logger)
		// Write CSV header
		fl.csvWriter.Write([]string{"Timestamp", "ClientIP", "Domain", "QueryType", "ResponseCode", "AnswerCount", "Answers", "DurationMs", "DoHServer"})
		fl.csvWriter.Flush()
	}

	return fl, nil
}

func (l *FileLogger) Log(entry QueryLogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.format == "json" {
		return l.jsonWriter.Encode(entry)
	} else if l.format == "csv" {
		// Serialize answers to JSON string for CSV
		answersJSON := ""
		if len(entry.Answers) > 0 {
			answersBytes, _ := json.Marshal(entry.Answers)
			answersJSON = string(answersBytes)
		}

		record := []string{
			entry.Timestamp.Format(time.RFC3339),
			entry.ClientIP,
			entry.Domain,
			entry.QueryType,
			entry.ResponseCode,
			fmt.Sprintf("%d", entry.AnswerCount),
			answersJSON,
			fmt.Sprintf("%d", entry.Duration),
			entry.DoHServer,
		}
		if err := l.csvWriter.Write(record); err != nil {
			return err
		}
		l.csvWriter.Flush()
		return l.csvWriter.Error()
	}

	return fmt.Errorf("unsupported file format: %s", l.format)
}

func (l *FileLogger) Close() error {
	if l.csvWriter != nil {
		l.csvWriter.Flush()
	}
	return l.logger.Close()
}

// DatabaseLogger logs to database (SQLite or PostgreSQL)
type DatabaseLogger struct {
	db     *sql.DB
	dbType string
	mu     sync.Mutex
}

func NewDatabaseLogger(config *Config) (*DatabaseLogger, error) {
	dbType := config.Logging.QueryLog.Database.Type
	var db *sql.DB
	var err error

	switch dbType {
	case "sqlite":
		dbPath := config.Logging.QueryLog.Database.SQLite.Path
		dbDir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %v", err)
		}

		db, err = sql.Open("sqlite", dbPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open SQLite database: %v", err)
		}

	case "postgresql":
		pgConfig := config.Logging.QueryLog.Database.PostgreSQL
		connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			pgConfig.Host, pgConfig.Port, pgConfig.User, pgConfig.Password, pgConfig.Database, pgConfig.SSLMode)

		db, err = sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to open PostgreSQL database: %v", err)
		}

	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	dl := &DatabaseLogger{
		db:     db,
		dbType: dbType,
	}

	// Create table
	if err := dl.createTable(); err != nil {
		return nil, err
	}

	return dl, nil
}

func (l *DatabaseLogger) createTable() error {
	var createTableSQL string

	if l.dbType == "sqlite" {
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS query_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			timestamp DATETIME NOT NULL,
			client_ip TEXT NOT NULL,
			domain TEXT NOT NULL,
			query_type TEXT NOT NULL,
			response_code TEXT NOT NULL,
			answer_count INTEGER NOT NULL,
			answers TEXT,
			duration_ms INTEGER NOT NULL,
			doh_server TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_timestamp ON query_logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_domain ON query_logs(domain);
		CREATE INDEX IF NOT EXISTS idx_client_ip ON query_logs(client_ip);
		`
	} else {
		createTableSQL = `
		CREATE TABLE IF NOT EXISTS query_logs (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP NOT NULL,
			client_ip VARCHAR(45) NOT NULL,
			domain VARCHAR(255) NOT NULL,
			query_type VARCHAR(10) NOT NULL,
			response_code VARCHAR(20) NOT NULL,
			answer_count INTEGER NOT NULL,
			answers TEXT,
			duration_ms INTEGER NOT NULL,
			doh_server VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		CREATE INDEX IF NOT EXISTS idx_timestamp ON query_logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_domain ON query_logs(domain);
		CREATE INDEX IF NOT EXISTS idx_client_ip ON query_logs(client_ip);
		`
	}

	_, err := l.db.Exec(createTableSQL)
	return err
}

func (l *DatabaseLogger) Log(entry QueryLogEntry) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	// Serialize answers to JSON string
	answersJSON := ""
	if len(entry.Answers) > 0 {
		answersBytes, _ := json.Marshal(entry.Answers)
		answersJSON = string(answersBytes)
	}

	query := `INSERT INTO query_logs 
		(timestamp, client_ip, domain, query_type, response_code, answer_count, answers, duration_ms, doh_server)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	if l.dbType == "sqlite" {
		query = `INSERT INTO query_logs 
			(timestamp, client_ip, domain, query_type, response_code, answer_count, answers, duration_ms, doh_server)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	}

	_, err := l.db.Exec(query,
		entry.Timestamp,
		entry.ClientIP,
		entry.Domain,
		entry.QueryType,
		entry.ResponseCode,
		entry.AnswerCount,
		answersJSON,
		entry.Duration,
		entry.DoHServer,
	)

	return err
}

func (l *DatabaseLogger) Close() error {
	return l.db.Close()
}

// NewQueryLogger creates a query logger based on configuration
func NewQueryLogger(config *Config) (QueryLogger, error) {
	if !config.Logging.QueryLog.Enabled {
		return NewConsoleLogger(), nil
	}

	switch config.Logging.QueryLog.Target {
	case "console":
		log.Println("Query logging: Console")
		return NewConsoleLogger(), nil

	case "file":
		log.Printf("Query logging: File (%s format) -> %s",
			config.Logging.QueryLog.File.Format,
			config.Logging.QueryLog.File.Path)
		return NewFileLogger(config)

	case "database":
		log.Printf("Query logging: Database (%s)",
			config.Logging.QueryLog.Database.Type)
		return NewDatabaseLogger(config)

	default:
		return nil, fmt.Errorf("unsupported query log target: %s", config.Logging.QueryLog.Target)
	}
}
