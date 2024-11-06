package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

var customTransport = &http.Transport{
	TLSClientConfig: &tls.Config{
		InsecureSkipVerify: true, // Ignorar verificação de certificado (não recomendado em produção)
	},
	Proxy: http.ProxyFromEnvironment, // Usa o proxy configurado no sistema

}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Println("Incoming request")

	targetURL := r.URL
	proxyReq, err := http.NewRequest(r.Method, targetURL.String(), r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	for name, values := range r.Header {
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	resp, err := customTransport.RoundTrip(proxyReq)
	if err != nil {
		http.Error(w, "Error sending proxy request", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	// Log de sucesso
	log.Println("Request successfully proxied to", targetURL.String())
}

func main() {
	http.HandleFunc("/", handleRequest)

	// Start HTTP server
	log.Println("Starting HTTP server on :8010")
	if err := http.ListenAndServe("0.0.0.0:8010", nil); err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}
