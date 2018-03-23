package main

import (
	"net/http"
	"os"

	"github.com/cloudfoundry/secure-credentials-broker/broker"
	"github.com/cloudfoundry/secure-credentials-broker/brokerconfig"

	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
	"github.com/cloudfoundry-incubator/credhub-cli/util"
	"github.com/pivotal-cf/brokerapi"
)

func main() {
	brokerLogger := lager.NewLogger("secure-credentials-broker")
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))
	brokerLogger.RegisterSink(lager.NewWriterSink(os.Stderr, lager.ERROR))
	brokerLogger.Info("the secure credentials broker is starting up...")

	config, err := brokerconfig.LoadAndValidateConfig()
	if err != nil {
		brokerLogger.Error("Broker is not configured correctly", err)
		os.Exit(1)
	}

	credHubClient, err := authenticate(config)
	if err != nil {
		brokerLogger.Error("CredHub client is failed to create", err)
		os.Exit(2)
	}

	serviceBroker := &broker.CredhubServiceBroker{CredHubClient: credHubClient, Logger: brokerLogger}

	brokerCredentials := brokerapi.BrokerCredentials{
		Username: config.BrokerConfiguration.Username,
		Password: config.BrokerConfiguration.Password,
	}

	brokerAPI := brokerapi.New(serviceBroker, brokerLogger, brokerCredentials)

	http.Handle("/", brokerAPI)

	brokerLogger.Info("the secure credentials broker is up and listening on port: " + config.BrokerConfiguration.Port)
	brokerLogger.Fatal("http-listen", http.ListenAndServe(":"+config.BrokerConfiguration.Port, nil))
}

func authenticate(config *brokerconfig.Config) (*credhub.CredHub, error) {
	ch, err := credhub.New(
		util.AddDefaultSchemeIfNecessary(config.CredHubConfiguration.ServerURL),
		credhub.SkipTLSValidation(config.CredHubConfiguration.SkipTLSValidation),
		credhub.Auth(auth.UaaClientCredentials(config.CredHubConfiguration.UAAClient, config.CredHubConfiguration.UAASecret)),
	)

	return ch, err
}
