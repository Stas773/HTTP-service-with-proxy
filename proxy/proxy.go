package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"os/signal"
	"time"
)

const proxyAddr string = ":9000"

var counter int = 0

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	port1 := flag.String("port1", "8081", "port number")
	port2 := flag.String("port2", "8082", "port number")
	flag.Parse()
	proxy := httputil.ReverseProxy{
		Director: func(req *http.Request) {
			if counter == 0 {
				req.URL.Scheme = "http"
				req.URL.Host = "localhost:" + *port1
				counter++
			} else {
				req.URL.Scheme = "http"
				req.URL.Host = "localhost:" + *port2
				counter--
			}
		},
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})

	ProxyServer := &http.Server{
		Addr:    proxyAddr,
		Handler: &proxy,
	}

	go func() {
		fmt.Printf("Proxy started on port: localhost%s and listening ports: %s, %s\n", proxyAddr, *port1, *port2)
		log.Fatal(http.ListenAndServe(proxyAddr, nil))
	}()

	<-done
	fmt.Println("Stop signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := ProxyServer.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown error %s\n", err)
	}
	fmt.Println("Proxy server stoped")
}
