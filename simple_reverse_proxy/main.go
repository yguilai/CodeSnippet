package main

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
)

// 反向代理简单实现

var (
	proxy_addr = "http://127.0.0.1:2001"
	port = ":2000"
)

func handler(w http.ResponseWriter, r *http.Request)  {
	proxy, err := url.Parse(proxy_addr)
	if err != nil {
		log.Println(err)
		return
	}
	r.URL.Scheme = proxy.Scheme
	r.URL.Host = proxy.Host

	transport := http.DefaultTransport
	resp, err := transport.RoundTrip(r)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println(err)
		}
	}()

	for k, val := range resp.Header {
		for _, v := range val {
			w.Header().Add(k, v)
		}
	}

	_, err = bufio.NewReader(resp.Body).WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Println("Start Reverse Proxy Serve On Port " + port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln(err)
	}
}