package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	certFile = flag.String("cert", "path to pem file", "A PEM eoncoded certificate file.")
	keyFile  = flag.String("key", "path to pem file", "A PEM encoded private key file.")
	caFile   = flag.String("CA", "path to pem file", "A PEM eoncoded CA's certificate file.")
)

func main() {
	flag.Parse()

	// Load client cert
	cert, err := tls.LoadX509KeyPair(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}

	// Load CA cert
	caCert, err := ioutil.ReadFile(*caFile)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		//InsecureSkipVerify: true,
	}
	tlsConfig.BuildNameToCertificate()
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	// Do GET something
	resp, err := client.Get("https://api/url")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Write response to a slice of PatchStatus Structs
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var output []PatchStatus
	err = json.Unmarshal(data, &output)
	if err != nil {
		panic(err)
	}
}
