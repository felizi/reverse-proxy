/*
Generate private key (.key)
# Key considerations for algorithm "RSA" ≥ 2048-bit
    openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
    openssl ecparam -genkey -name secp384r1 -out server.key
Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
    openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
*/
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	log.Println(":: Reverse Proxy ::")
	rawURL := flag.String("u", "", "URL to proxy")
	port := flag.Int("p", 443, "Port to serve. Default: 443")
	certFile := flag.String("c", "server.crt", "Cert file. Default: server.crt")
	keyFile := flag.String("k", "server.key", "Key file. Default: server.key")

	flag.Parse()

	target, err := url.Parse(*rawURL)
	if err != nil || !isURL(*rawURL) {
		log.Fatal(fmt.Sprintf("Invalid URL: '%v' ", *rawURL), err)
	}

	log.Printf("Configuring proxy: %v", target)
	proxy := httputil.NewSingleHostReverseProxy(target)

	proxy.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	log.Printf("Proxing 'https://localhost:%v' to '%v' with cert '%v' and key '%v'", *port, *rawURL, *certFile, *keyFile)
	err = http.ListenAndServeTLS(fmt.Sprintf(":%v", *port), *certFile, *keyFile, proxy)
	if err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
