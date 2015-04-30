package server

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"net/http"
	"sync"

	bosherr "github.com/cloudfoundry/bosh-agent/errors"
	boshlog "github.com/cloudfoundry/bosh-agent/logger"
)

const RegistryListenerLogTag = "RegistryListener"

type Listener struct {
	config          ListenerConfig
	instanceHandler *InstanceHandler
	logger          boshlog.Logger
	listener        net.Listener
	waitGroup       sync.WaitGroup
}

func NewListener(
	config ListenerConfig,
	instanceHandler *InstanceHandler,
	logger boshlog.Logger,
) Listener {
	return Listener{
		config:          config,
		instanceHandler: instanceHandler,
		logger:          logger,
	}
}

func (l *Listener) ListenAndServe() (err error) {
	tcpAddr := &net.TCPAddr{
		IP:   net.ParseIP(l.config.Address),
		Port: l.config.Port,
	}
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return bosherr.WrapError(err, "Starting Registry TCP Listener")
	}

	if l.config.Protocol == "https" {
		certificates, err := tls.LoadX509KeyPair(l.config.TLS.CertFile, l.config.TLS.KeyFile)
		if err != nil {
			return bosherr.WrapError(err, "Loading X509 Key Pair")
		}

		certPool := x509.NewCertPool()
		if l.config.TLS.CACertFile != "" {
			caCert, err := ioutil.ReadFile(l.config.TLS.CACertFile)
			if err != nil {
				return bosherr.WrapError(err, "Loading CA certificate")
			}
			if !certPool.AppendCertsFromPEM(caCert) {
				return bosherr.WrapError(err, "Invalid CA Certificate")
			}
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{certificates},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    certPool,
			MinVersion:   tls.VersionTLS12,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
			SessionTicketsDisabled:   true,
		}

		l.listener = tls.NewListener(tcpListener, tlsConfig)
	} else {
		l.listener = tcpListener
	}

	httpServer := http.Server{}
	mux := http.NewServeMux()
	mux.HandleFunc("/instances/", l.instanceHandler.HandleFunc)
	httpServer.Handler = mux

	l.logger.Debug(RegistryListenerLogTag, "Starting Registry Server at %s://%s:%d", l.config.Protocol, l.config.Address, l.config.Port)
	l.waitGroup.Add(1)
	go func() {
		defer l.waitGroup.Done()
		err := httpServer.Serve(l.listener)
		if err != nil {
			l.logger.Debug(RegistryListenerLogTag, "Unexpected server shutdown: %#v", err)
		}
	}()

	return nil
}

func (l *Listener) WaitForServerToExit() {
	l.waitGroup.Wait()
}
