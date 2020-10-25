package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	s1 := &Server{Addr: "127.0.0.1:2001"}
	s1.Run()
	s2 := &Server{Addr: "127.0.0.1:2002"}
	s2.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type Server struct {
	Addr string
}

func (s *Server) Run() {
	log.Println("Start Server On " + s.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HelloHandler)
	mux.HandleFunc("/base/err", s.ErrorHandler)
	mux.HandleFunc("/t1/t2/timeout", s.TimeoutHandler)

	server := &http.Server{
		Addr: s.Addr,
		WriteTimeout: time.Second * 3,
		Handler: mux,
	}

	go func() {
		log.Fatalln(server.ListenAndServe())
	}()
}

func (s *Server) HelloHandler(w http.ResponseWriter, r *http.Request) {
	p := fmt.Sprintf("http://%s%s\n", s.Addr, r.URL.Path)
	ip := fmt.Sprintf("RemoteAddr=%s,X-Forwarded-For=%v,X-Real-Ip=%v\n", r.RemoteAddr, r.Header.Get("X-Forwarded-For"), r.Header.Get("X-Real-Ip"))
	h := fmt.Sprintf("headers =%v\n", r.Header)
	io.WriteString(w, p)
	io.WriteString(w, ip)
	io.WriteString(w, h)
}

func (s *Server) ErrorHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "Error Handler")
}

func (s *Server) TimeoutHandler(w http.ResponseWriter, req *http.Request)  {
	time.Sleep(6 * time.Second)
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Timeout Handler")
}

