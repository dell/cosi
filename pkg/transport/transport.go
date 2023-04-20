// Copyright © 2023 Dell Inc. or its subsidiaries. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package transport

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/dell/cosi-driver/pkg/config"
	"github.com/dell/cosi-driver/util"
)

var (
	// ErrClientCertMissing indicates that one of client-cert or client-key is missing.
	ErrClientCertMissing = errors.New("client-cert or client-key missing")

	// ErrRootCAMissing indicates that root CA (certificate authority) is missing.
	ErrRootCAMissing = errors.New("root certificate authority is missing")

	// ErrRootCAMissing indicates that root CA (certificate authority) is missing.
	ErrCertCreating = errors.New("unable to create TLS Certificate")

	// ErrRootCAMissing indicates that root CA (certificate authority) is missing.
	ErrRootCasDecoding = errors.New("unable to decode RootCas")
)

// New creates new HTTP or HTTPS transport based on provided config.
func New(cfg config.Tls) (*http.Transport, error) {
	var tlsConfig *tls.Config
	if cfg.Insecure {
		/* #nosec */
		tlsConfig = &tls.Config{InsecureSkipVerify: true}
	} else {
		cert, err := clientCert(cfg.ClientCert, cfg.ClientKey)
		if err != nil {
			log.WithFields(log.Fields{
				"error_msg": err,
			}).Error("unable to create TLS Certificate")
			return nil, ErrCertCreating
		}

		if cfg.RootCas == nil || *cfg.RootCas == "" {
			log.WithFields(log.Fields{
				"error_msg": err,
			}).Error("root certificate authority is missing")
			return nil, ErrRootCAMissing
		}

		b, err := base64.StdEncoding.DecodeString(*cfg.RootCas)
		if err != nil {
			log.WithFields(log.Fields{
				"error_msg": err,
			}).Error("unable to decode RootCas")
			return nil, ErrRootCasDecoding
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(b)

		tlsConfig = &tls.Config{
			InsecureSkipVerify: false,
			MinVersion:         tls.VersionTLS12,
			Certificates:       cert,
			RootCAs:            caCertPool,
		}
	}

	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}, nil
}

func clientCert(certData, keyData *string) ([]tls.Certificate, error) {
	if certData == nil && keyData == nil {
		// no certificates and key is a valid option
		return []tls.Certificate{}, nil
	} else if certData == nil || keyData == nil {
		// only one of those two is missing, it is not a valid option
		log.WithFields(log.Fields{
			"cert_data": certData,
			"key_data":  keyData,
		}).Error("client-cert or client-key missing")
		return nil, ErrClientCertMissing
	}

	if *certData == "" && *keyData == "" {
		// both certificate and key are empty, this is also a valid option
		return []tls.Certificate{}, nil
	} else if *certData == "" || *keyData == "" {
		// only one of those two is empty, it is not a valid option
		log.WithFields(log.Fields{
			"cert_data": certData,
			"key_data":  keyData,
		}).Error("client-cert or client-key missing")
		return nil, ErrClientCertMissing
	}

	// decode both certificate and key
	cert, err := base64.StdEncoding.DecodeString(*certData)
	if err != nil {
		return nil, util.TraceLogging(err, "unable to decode client-cert")
	}

	key, err := base64.StdEncoding.DecodeString(*keyData)
	if err != nil {
		return nil, util.TraceLogging(err, "unable to decode client-key")
	}

	// parse a public/private key pair from a pair of PEM encoded data
	x509, err := tls.X509KeyPair(cert, key)
	if err != nil {
		return nil, util.TraceLogging(err, "unable to parse a public/private key pair")
	}

	return []tls.Certificate{x509}, nil
}
