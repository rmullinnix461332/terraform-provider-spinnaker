// Copyright (c) 2018, Google, Inc.
// Copyright (c) 2019, Noel Cower.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package gateclient

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/pkg/errors"

	spin509 "github.com/spinnaker/spin/config/auth/x509"
	gate "github.com/spinnaker/spin/gateapi"
	"github.com/spinnaker/spin/version"
)

const (
	// defaultConfigFileMode is the default file mode used for config files. This corresponds to
	// the Unix file permissions u=rw,g=,o= so that config files with cached tokens, at least by
	// default, are only readable by the user that owns the config file.
	defaultConfigFileMode os.FileMode = 0600 // u=rw,g=,o=
)

// GatewayClient is the wrapper with authentication
type GatewayClient struct {
	// The exported fields below should be set by anyone using a command
	// with an GatewayClient field. These are expected to be set externally
	// (not from within the command itself).

	// Generate Gate Api client.
	*gate.APIClient

	Context context.Context
	// Spin CLI configuration.
	client_x509 *spin509.Config

	// This is the set of flags global to the command parser.
	gateEndpoint     string
	ignoreCertErrors bool
	ignoreRedirects  bool

	// Raw Http Client to do OAuth2 login.
	httpClient *http.Client

	// Maximum time to wait (when polling) for a task to become completed.
	retryTimeout int
}

func (m *GatewayClient) GateEndpoint() string {
	if m.gateEndpoint == "" {
		return "http://localhost:8085"
	}
	return m.gateEndpoint
}

func (m *GatewayClient) RetryTimeout() int {
	return m.retryTimeout
}

// Create new spinnaker gateway client with flag
func NewGateClient(gateEndpoint, defaultHeaders, x509_cert string, x509_key string, ignoreCertErrors bool) (*GatewayClient, error) {
	gateClient := &GatewayClient{
		gateEndpoint:     gateEndpoint,
		ignoreCertErrors: ignoreCertErrors,
		ignoreRedirects:  false,
		retryTimeout:     60,
		Context:          context.Background(),
	}

	X509 := new(spin509.Config)
	X509.Cert = x509_cert
	X509.Key = x509_key

	gateClient.client_x509 = X509

	// Api client initialization.
	err := gateClient.InitializeHTTPClient()
	if err != nil {
		return nil, errors.New("Could not initialize http client, failing.")
	}

	m := make(map[string]string)

	if defaultHeaders != "" {
		headers := strings.Split(defaultHeaders, ",")
		for _, element := range headers {
			header := strings.SplitN(element, "=", 2)
			if len(header) != 2 {
				return nil, fmt.Errorf("Bad default-header value, use key=value form: %s", element)
			}
			m[strings.TrimSpace(header[0])] = strings.TrimSpace(header[1])
		}
	}

	cfg := &gate.Configuration{
		BasePath:      gateClient.GateEndpoint(),
		DefaultHeader: m,
		UserAgent:     fmt.Sprintf("%s/%s", version.UserAgent, version.String()),
		HTTPClient:    gateClient.httpClient,
	}

	gateClient.APIClient = gate.NewAPIClient(cfg)

	// TODO: Verify version compatibility between Spin CLI and Gate.
	_, _, err = gateClient.VersionControllerApi.GetVersionUsingGET(gateClient.Context)
	if err != nil {
		return nil, errors.New("Could not reach Gate, please ensure it is running. Failing.")
	}

	return gateClient, nil
}

// InitializeHTTPClient will return an *http.Client with TLS
func (m *GatewayClient) InitializeHTTPClient() error {
	cookieJar, _ := cookiejar.New(nil)
	client := http.Client{
		Jar:       cookieJar,
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	client.Transport.(*http.Transport).TLSClientConfig = &tls.Config{
		InsecureSkipVerify: m.ignoreCertErrors,
	}

	if !m.client_x509.IsValid() {
		// Misconfigured.
		return errors.New("Incorrect x509 auth configuration.\nMust specify certPath/keyPath or cert/key pair.")
	}

	if m.client_x509.Cert != "" && m.client_x509.Key != "" {
		certBytes := []byte(m.client_x509.Cert)
		keyBytes := []byte(m.client_x509.Key)
		cert, err := tls.X509KeyPair(certBytes, keyBytes)

		if err != nil {
			return err
		}

		clientCertPool := x509.NewCertPool()
		clientCertPool.AppendCertsFromPEM(certBytes)

		client.Transport.(*http.Transport).TLSClientConfig.MinVersion = tls.VersionTLS12
		client.Transport.(*http.Transport).TLSClientConfig.PreferServerCipherSuites = true
		client.Transport.(*http.Transport).TLSClientConfig.Certificates = []tls.Certificate{cert}

		m.httpClient = &client

		return nil
	}

	// Misconfigured.
	return errors.New("Incorrect x509 auth configuration.\nMust specify certPath/keyPath or cert/key pair.")
}
