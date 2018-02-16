package main

import (
	"crypto/tls"
	"io"
	"log"
	"net"
	"net/http"
	"time"
	"fmt"
)

const (
	HTTPProtocol  = "http"
	HTTPSProtocol = "https"
)

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	dest_conn, err := net.DialTimeout("tcp", r.Host, 10*time.Second)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}
	client_conn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
	}
	go transfer(dest_conn, client_conn)
	go transfer(client_conn, dest_conn)

	log.Printf("handling HTTPS request: %9.2f\n", time.Now().Sub(start).Seconds())
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)

	fmt.Printf("Proxy [%s] | handled in  %9.2f sec\n", req.RequestURI, time.Now().Sub(start).Seconds())
}

func transfer(destination io.WriteCloser, source io.ReadCloser) {
	defer destination.Close()
	defer source.Close()
	io.Copy(destination, source)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}


func showBootInfo() {
	var logo string = `
,------.,------.
|  .---'|  .--. ',--.--. ,---.,--.  ,--.,--. ,--.
|  "--, |  '--' ||  .--'| .-. |\  "'  /  \  '  /
|  "---.|  | --' |  |   ' '-' '/  /.  \   \   '
"------""--'     "--'    "---''--'  '--'.-'  /
                                        "---'     `
	fmt.Println(logo)
	fmt.Println("---------------------------------------------------")
	fmt.Printf("PEM path:  %s\n", conf.PemPath)
	fmt.Printf("KEY path:  %s\n", conf.KeyPath)
	fmt.Printf("PROTO mod: %s\n", conf.Proto)
	fmt.Printf("PORT:      %s\n", conf.Port)
	fmt.Println("---------------------------------------------------")
	fmt.Println("Awaiting for connections...")

}

func main() {
	showBootInfo()

	if conf.Proto != HTTPProtocol && conf.Proto != HTTPSProtocol{
		log.Fatalf("Protocol must be either % or %s, got %s", HTTPProtocol, HTTPSProtocol, conf.Proto)
	}

	server := &http.Server{
		Addr: ":" + conf.Port,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodConnect {
				handleTunneling(w, r)
			} else {
				handleHTTP(w, r)
			}
		}),
		// Disable HTTP/2.
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler)),
	}

	if conf.Proto == HTTPProtocol {
		log.Fatal(server.ListenAndServe())
	} else {
		log.Fatal(server.ListenAndServeTLS(conf.PemPath, conf.KeyPath))
	}
}
