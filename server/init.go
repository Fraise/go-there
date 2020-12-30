package server

import (
	"crypto/tls"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-there/config"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"strconv"
)

// Start initializes and starts configured http and https servers.
// It returns (server *http.Server, tlsServer *http.Server). If a server is not configured, nil is returned.
func Start(conf *config.Configuration, e *gin.Engine) (*http.Server, *http.Server) {
	var s *http.Server
	var tlsServer *http.Server

	if conf.Server.HttpListenPort > 0 {
		// HTTP server
		s = &http.Server{
			Addr:    conf.Server.ListenAddress + ":" + strconv.Itoa(conf.Server.HttpListenPort),
			Handler: e,
		}

		go func() {
			if err := s.ListenAndServe(); err != http.ErrServerClosed {
				log.Fatal().Err(err).Send()
			}
		}()
	}

	if conf.Server.HttpsListenPort > 0 {

		// HTTPS server
		tlsServer = &http.Server{
			Addr:    conf.Server.ListenAddress + ":" + strconv.Itoa(conf.Server.HttpsListenPort),
			Handler: e,
		}

		if conf.Server.UseAutoCert {
			if len(conf.Server.Domains) <= 0 {
				log.Fatal().Msg("no domain configured for autocert")
			}

			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(conf.Server.Domains...),
				Cache:      autocert.DirCache(conf.Server.CertCache),
			}

			// Use let's encrypt autocert
			tlsServer.TLSConfig = &tls.Config{
				GetCertificate: m.GetCertificate,
			}

			// Remove provided path when using autocert to avoid conflicts
			conf.Server.CertPath = ""
			conf.Server.KeyPath = ""
		}

		go func() {
			if err := tlsServer.ListenAndServeTLS(conf.Server.CertPath, conf.Server.KeyPath); err != http.ErrServerClosed {
				log.Fatal().Err(err).Send()
			}
		}()
	}

	return s, tlsServer
}
