package config

import (
	"crypto/tls"
	"fmt"
	"os"
)

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableTLS       bool
	TLSCertFile     string
	TLSKeyFile      string
	MinTLSVersion   uint16
	MaxTLSVersion   uint16
	CipherSuites    []uint16
	RequireClientCert bool
}

// LoadSecurityConfig 載入安全配置
func LoadSecurityConfig() *SecurityConfig {
	config := &SecurityConfig{
		EnableTLS:       getEnvBool("ENABLE_TLS", false),
		TLSCertFile:     getEnv("TLS_CERT_FILE", ""),
		TLSKeyFile:      getEnv("TLS_KEY_FILE", ""),
		MinTLSVersion:   tls.VersionTLS12, // 最低 TLS 1.2
		MaxTLSVersion:   tls.VersionTLS13, // 最高 TLS 1.3
		RequireClientCert: getEnvBool("REQUIRE_CLIENT_CERT", false),
		CipherSuites: []uint16{
			// 推薦的安全加密套件
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}
	
	return config
}

// GetTLSConfig 獲取 TLS 配置
func (sc *SecurityConfig) GetTLSConfig() (*tls.Config, error) {
	if !sc.EnableTLS {
		return nil, nil
	}
	
	if sc.TLSCertFile == "" || sc.TLSKeyFile == "" {
		return nil, fmt.Errorf("TLS cert file and key file are required when TLS is enabled")
	}
	
	cert, err := tls.LoadX509KeyPair(sc.TLSCertFile, sc.TLSKeyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate: %w", err)
	}
	
	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   sc.MinTLSVersion,
		MaxVersion:   sc.MaxTLSVersion,
		CipherSuites: sc.CipherSuites,
	}
	
	if sc.RequireClientCert {
		config.ClientAuth = tls.RequireAndVerifyClientCert
	}
	
	return config, nil
}

// ValidateSecurityConfig 驗證安全配置
func (sc *SecurityConfig) ValidateSecurityConfig() error {
	if sc.EnableTLS {
		if sc.TLSCertFile == "" {
			return fmt.Errorf("TLS_CERT_FILE is required when TLS is enabled")
		}
		
		if sc.TLSKeyFile == "" {
			return fmt.Errorf("TLS_KEY_FILE is required when TLS is enabled")
		}
		
		// 檢查檔案是否存在
		if _, err := os.Stat(sc.TLSCertFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS cert file not found: %s", sc.TLSCertFile)
		}
		
		if _, err := os.Stat(sc.TLSKeyFile); os.IsNotExist(err) {
			return fmt.Errorf("TLS key file not found: %s", sc.TLSKeyFile)
		}
	}
	
	return nil
}
