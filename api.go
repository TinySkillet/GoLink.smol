package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type GoLinkServer struct {
	listenAddr string
	store      *RedisStore
}

func NewGoLinkServer(addr string) *GoLinkServer {
	store := NewRedisStore()
	return &GoLinkServer{addr, store}
}

func (g *GoLinkServer) Run() {

	router := http.NewServeMux()

	// server static files
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// handle redirect
	router.HandleFunc("/{id}", g.HandleRedirect)

	// templates
	g.loadTemplates(router)

	// receives htmx request
	router.HandleFunc("POST /make-it-smol", g.shortenURL)

	server := &http.Server{
		Addr:    g.listenAddr,
		Handler: router,
	}

	go func() {
		log.Println("Server started on PORT", server.Addr)
		log.Fatal(server.ListenAndServe())
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	<-sigChan
	fmt.Println("Received terminate, graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	server.Shutdown(ctx)
}

func (g *GoLinkServer) loadTemplates(router *http.ServeMux) {

	tmpl := template.Must(template.ParseGlob("./templates/*"))

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	})

	router.HandleFunc("/not-found", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "notfound.html", nil)
		if err != nil {
			log.Fatal(err)
		}
	})
}
