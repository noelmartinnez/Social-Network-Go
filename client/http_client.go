package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
	"net/http"
)

// Función para crear un cliente HTTP con TLS
func newTLSClient(certFile string) *http.Client {
	cert, err := ioutil.ReadFile(certFile)
	if err != nil {
		log.Fatalf("Failed to read certificate file: %v", err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(cert)

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}
