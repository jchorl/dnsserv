package common

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"
)

// UpdateRequest is the body for a request to update an ip address
type UpdateRequest struct {
	Domain string
}

// LoadTLSConfigOrPanic loads tls config or panics
func LoadTLSConfigOrPanic(caPath, certPath, keyPath string) *tls.Config {
	rootPEM, err := ioutil.ReadFile(caPath)
	if err != nil {
		log.Fatalf("Unable to load ca cert: %s\n", err)
	}
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM(rootPEM)
	if !ok {
		log.Fatal("Failed to parse root certificate")
	}

	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("Unable to load cert: %s\n", err)
	}

	tlsConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      roots,
		ClientCAs:    roots,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	return &tlsConfig
}
