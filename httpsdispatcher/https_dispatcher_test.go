package httpsdispatcher_test

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	boshdispatcher "github.com/cloudfoundry/bosh-agent/httpsdispatcher"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var _ = Describe("HTTPSDispatcher", func() {
	var (
		dispatcher *boshdispatcher.HTTPSDispatcher
	)

	BeforeEach(func() {
		logger := boshlog.NewLogger(boshlog.LevelNone)
		serverURL, err := url.Parse("https://127.0.0.1:7788")
		Expect(err).ToNot(HaveOccurred())
		dispatcher = boshdispatcher.NewHTTPSDispatcher(serverURL, logger)

		errChan := make(chan error)
		go func() {
			errChan <- dispatcher.Start()
		}()

		select {
		case err := <-errChan:
			Expect(err).ToNot(HaveOccurred())
		case <-time.After(1 * time.Second):
			// server should now be running, continue
		}
	})

	AfterEach(func() {
		dispatcher.Stop()
		time.Sleep(1 * time.Second)
	})

	It("calls the handler function for the route", func() {
		var hasBeenCalled = false
		handler := func(w http.ResponseWriter, r *http.Request) {
			hasBeenCalled = true
			w.WriteHeader(201)
		}

		dispatcher.AddRoute("/example", handler)

		client := getHTTPClient()
		response, err := client.Get("https://127.0.0.1:7788/example")

		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(BeNumerically("==", 201))
		Expect(hasBeenCalled).To(Equal(true))
	})

	It("returns a 404 if the route does not exist", func() {
		client := getHTTPClient()
		response, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).ToNot(HaveOccurred())
		Expect(response.StatusCode).To(BeNumerically("==", 404))
	})

	// Go's TLS client does not support SSLv3 (so we couldn't test it even if it did)
	PIt("does not allow connections using SSLv3", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionSSL30,
			MaxVersion:         tls.VersionSSL30,
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).To(HaveOccurred())
	})

	It("does allow connections using TLSv1", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS10,
			MaxVersion:         tls.VersionTLS10,
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).ToNot(HaveOccurred())
	})

	It("does allow connections using TLSv1.1", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS11,
			MaxVersion:         tls.VersionTLS11,
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).ToNot(HaveOccurred())
	})

	It("does allow connections using TLSv1.2", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS12,
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).ToNot(HaveOccurred())
	})

	It("does not allow connections using 3DES ciphers", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
			},
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).To(HaveOccurred())
	})

	It("does not allow connections using RC4 ciphers", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
			CipherSuites: []uint16{
				tls.TLS_RSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_RC4_128_SHA,
				tls.TLS_ECDHE_RSA_WITH_RC4_128_SHA,
			},
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).To(HaveOccurred())
	})

	It("does allow connections using AES ciphers", func() {
		handler := func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(200) }
		dispatcher.AddRoute("/example", handler)

		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
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
		}
		client := getHTTPClientWithConfig(tlsConfig)
		_, err := client.Get("https://127.0.0.1:7788/example")
		Expect(err).ToNot(HaveOccurred())
	})
})

func getHTTPClient() http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		// Both CBC & RC4 ciphers can be exploited
		// Mozilla's "Modern" recommended settings only overlap with the golang TLS client on these two ciphers
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		},
		// SSLv3 and TLSv1.0 are considered weak
		// TLS1.1 does not support GCM, so it won't actually be used
		MinVersion: tls.VersionTLS11,
		MaxVersion: tls.VersionTLS12,
	}
	return getHTTPClientWithConfig(tlsConfig)
}

func getHTTPClientWithConfig(tlsConfig *tls.Config) http.Client {
	httpTransport := &http.Transport{TLSClientConfig: tlsConfig}
	return http.Client{Transport: httpTransport}
}
