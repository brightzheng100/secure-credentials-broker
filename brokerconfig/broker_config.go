package brokerconfig

import (
	"errors"
	"os"
	"strconv"
)

// Config contains the broker's primary configuration
type Config struct {
	BrokerConfiguration
	CredHubConfiguration
}

// BrokerConfiguration contains broker's configuration info
type BrokerConfiguration struct {
	Port     string
	Username string
	Password string
}

// CredHubConfiguration contains CredHub's configuration info
type CredHubConfiguration struct {
	ServerURL         string
	SkipTLSValidation bool
	UAAClient         string
	UAASecret         string
}

// LoadConfig loads environment variables into Config
func LoadAndValidateConfig() (*Config, error) {
	config := &Config{}

	//Broker Config
	port := "8080"
	if len(os.Getenv("PORT")) > 0 {
		if _, err := strconv.Atoi(os.Getenv("PORT")); err == nil {
			port = os.Getenv("PORT")
		}
	}
	config.BrokerConfiguration.Port = port

	username := os.Getenv("BROKER_AUTH_USERNAME")
	if len(username) == 0 {
		return config, errors.New("BROKER_AUTH_USERNAME is NOT set")
	}
	config.BrokerConfiguration.Username = username

	password := os.Getenv("BROKER_AUTH_PASSWORD")
	if len(password) == 0 {
		return config, errors.New("BROKER_AUTH_PASSWORD is NOT set")
	}
	config.BrokerConfiguration.Password = password

	//CredHub Config
	credhubServer := os.Getenv("CREDHUB_SERVER")
	if len(credhubServer) == 0 {
		return config, errors.New("CREDHUB_SERVER is NOT set")
	}
	config.CredHubConfiguration.ServerURL = credhubServer

	skipTLSValidation := false
	if skipTLS := os.Getenv("SKIP_TLS_VALIDATION"); skipTLS == "true" {
		skipTLSValidation = true
	}
	config.CredHubConfiguration.SkipTLSValidation = skipTLSValidation

	uaaCient := os.Getenv("CREDHUB_CLIENT")
	if len(uaaCient) == 0 {
		return config, errors.New("CREDHUB_CLIENT is NOT set")
	}
	config.CredHubConfiguration.UAAClient = uaaCient

	uaaSecret := os.Getenv("CREDHUB_SECRET")
	if len(uaaSecret) == 0 {
		return config, errors.New("CREDHUB_SECRET is NOT set")
	}
	config.CredHubConfiguration.UAASecret = uaaSecret

	return config, nil
}
