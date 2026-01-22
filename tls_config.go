package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"
)

// TLSConfigManager manages TLS configuration and certificate validation
type TLSConfigManager struct {
	config *Config
}

// NewTLSConfigManager creates a new TLS config manager
func NewTLSConfigManager(config *Config) *TLSConfigManager {
	return &TLSConfigManager{
		config: config,
	}
}

// GetTLSConfig returns a configured tls.Config
func (m *TLSConfigManager) GetTLSConfig() *tls.Config {
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// If TLS advanced options are not enabled, return default config
	if !m.config.TLS.Enabled {
		return tlsConfig
	}

	// Set InsecureSkipVerify if configured
	if m.config.TLS.InsecureSkipVerify {
		log.Println("[WARNING] TLS verification is disabled - this is insecure!")
		tlsConfig.InsecureSkipVerify = true
		return tlsConfig
	}

	// Set up custom verification function
	if m.config.TLS.PinCertIssuer || m.config.TLS.PrintCertInfo {
		tlsConfig.VerifyPeerCertificate = m.verifyPeerCertificate
	}

	return tlsConfig
}

// verifyPeerCertificate is a custom certificate verification function
func (m *TLSConfigManager) verifyPeerCertificate(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	if len(verifiedChains) == 0 || len(verifiedChains[0]) == 0 {
		return fmt.Errorf("no verified certificate chains")
	}

	// Get the leaf certificate (server certificate)
	cert := verifiedChains[0][0]

	// Print certificate information if enabled
	if m.config.TLS.PrintCertInfo {
		m.printCertificateInfo(cert, verifiedChains[0])
	}

	// Verify certificate issuer if pinning is enabled
	if m.config.TLS.PinCertIssuer {
		if err := m.verifyCertificateIssuer(verifiedChains[0]); err != nil {
			return err
		}
	}

	return nil
}

// printCertificateInfo prints detailed certificate information
func (m *TLSConfigManager) printCertificateInfo(cert *x509.Certificate, chain []*x509.Certificate) {
	log.Println("=== TLS Certificate Information ===")
	log.Printf("Subject: %s", cert.Subject.CommonName)
	log.Printf("Issuer: %s", cert.Issuer.CommonName)
	log.Printf("Serial Number: %s", cert.SerialNumber.String())
	log.Printf("Valid From: %s", cert.NotBefore.Format("2006-01-02 15:04:05"))
	log.Printf("Valid Until: %s", cert.NotAfter.Format("2006-01-02 15:04:05"))

	// Print SANs (Subject Alternative Names)
	if len(cert.DNSNames) > 0 {
		log.Printf("DNS Names: %s", strings.Join(cert.DNSNames, ", "))
	}

	// Print certificate chain
	if len(chain) > 1 {
		log.Printf("Certificate Chain (%d certificates):", len(chain))
		for i, c := range chain {
			log.Printf("  [%d] %s (Issuer: %s)", i, c.Subject.CommonName, c.Issuer.CommonName)
		}
	}

	log.Println("===================================")
}

// verifyCertificateIssuer verifies that the certificate is issued by an allowed issuer
func (m *TLSConfigManager) verifyCertificateIssuer(chain []*x509.Certificate) error {
	if len(m.config.TLS.AllowedIssuers) == 0 {
		log.Println("[WARNING] Certificate issuer pinning is enabled but no allowed issuers configured")
		return nil
	}

	// Check each certificate in the chain
	for _, cert := range chain {
		issuerCN := cert.Issuer.CommonName

		// Check if issuer is in the allowed list
		for _, allowedIssuer := range m.config.TLS.AllowedIssuers {
			if strings.Contains(issuerCN, allowedIssuer) || strings.Contains(allowedIssuer, issuerCN) {
				log.Printf("[TLS] Certificate issuer verified: %s", issuerCN)
				return nil
			}
		}
	}

	// No matching issuer found
	allowedList := strings.Join(m.config.TLS.AllowedIssuers, ", ")
	return fmt.Errorf("certificate issuer not in allowed list. Allowed: [%s]", allowedList)
}

// ValidateTLSConfig validates the TLS configuration
func (m *TLSConfigManager) ValidateTLSConfig() error {
	if !m.config.TLS.Enabled {
		return nil
	}

	// Warn about insecure configurations
	if m.config.TLS.InsecureSkipVerify {
		log.Println("[WARNING] TLS certificate verification is disabled")
	}

	// Validate certificate pinning configuration
	if m.config.TLS.PinCertIssuer {
		if len(m.config.TLS.AllowedIssuers) == 0 {
			return fmt.Errorf("certificate issuer pinning is enabled but no allowed issuers specified")
		}
		log.Printf("[TLS] Certificate issuer pinning enabled with %d allowed issuers", len(m.config.TLS.AllowedIssuers))
	}

	if m.config.TLS.PrintCertInfo {
		log.Println("[TLS] Certificate information printing enabled")
	}

	return nil
}
