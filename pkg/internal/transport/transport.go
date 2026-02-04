// Copyright © 2023-2025 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// This software contains the intellectual property of Dell Inc.
// or is licensed to Dell Inc. from third parties. Use of this software
// and the intellectual property contained therein is expressly limited to the
// terms and conditions of the License Agreement under which it is provided by or
// on behalf of Dell Inc. or its subsidiaries.

// Package transport implements transport for HTTP client
// which is used further in custom client from goobjectscale.
package transport

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/dell/csmlog"

	"github.com/dell/cosi/pkg/config"
)

var (
	// ErrClientCertMissing indicates that one of client-cert or client-key is missing.
	ErrClientCertMissing = errors.New("client-cert or client-key missing")

	// ErrRootCAMissing indicates that root CA (certificate authority) is missing.
	ErrRootCAMissing = errors.New("root certificate authority is missing")
	log              = csmlog.GetLogger()
)

// New creates new HTTP or HTTPS transport based on provided config.
func New(cfg config.Tls) (*http.Transport, error) {
	var tlsConfig *tls.Config
	if cfg.Insecure {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
			CipherSuites:       getSecuredCipherSuites(),
		}

		log.Debug("Insecure connection applied")
	} else {
		cert, err := clientCert(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("unable to create TLS Certificate: %w", err)
		}

		if cfg.RootCas == nil || *cfg.RootCas == "" {
			return nil, ErrRootCAMissing
		}

		b, err := base64.StdEncoding.DecodeString(*cfg.RootCas)
		if err != nil {
			return nil, fmt.Errorf("unable to decode RootCas: %w", err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(b)

		tlsConfig = &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
			Certificates:       cert,
			RootCAs:            caCertPool,
			CipherSuites:       getSecuredCipherSuites(),
		}

		log.Debug("Secure connection applied")
	}

	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}, nil
}

func clientCert(certData, keyData *string) ([]tls.Certificate, error) {
	if certData == nil && keyData == nil {
		log.Debug("Default certificate created")
		// no certificates and key is a valid option
		return []tls.Certificate{}, nil
	} else if certData == nil || keyData == nil {
		// only one of those two is missing, it is not a valid option
		log.Error("client-cert or client-key missing")
		return nil, ErrClientCertMissing
	}

	if *certData == "" && *keyData == "" {
		// both certificate and key are empty, this is also a valid option
		log.Debug("Default certificate created")
		return []tls.Certificate{}, nil
	} else if *certData == "" || *keyData == "" {
		// only one of those two is empty, it is not a valid option
		log.Error("client-cert or client-key missing")
		return nil, ErrClientCertMissing
	}

	// decode both certificate and key
	cert, err := base64.StdEncoding.DecodeString(*certData)
	if err != nil {
		return nil, fmt.Errorf("unable to decode client-cert: %w", err)
	}

	key, err := base64.StdEncoding.DecodeString(*keyData)
	if err != nil {
		return nil, fmt.Errorf("unable to decode client-key: %w", err)
	}

	// parse a public/private key pair from a pair of PEM encoded data
	x509, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse a public/private key pair: %w", err)
	}
	return []tls.Certificate{x509}, nil
}

// getSecuredCipherSuites returns a slice of secured cipher suites.
// It iterates over the tls.CipherSuites() and appends the ID of each cipher suite to the suites slice.
// The function returns the suites slice.
func getSecuredCipherSuites() (suites []uint16) {
	securedSuite := tls.CipherSuites()
	for _, v := range securedSuite {
		suites = append(suites, v.ID)
	}

	return suites
}
