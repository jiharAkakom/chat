package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/onufert/trace"
)

// templ represents a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// ServeHTTP handle the http request
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {

	var addr = flag.String("addr", ":8080", "Address to serve application from")
	var debug = flag.Bool("debug", false, "Turn on debugging")
	flag.Parse()

	r := newRoom()
	if *debug {
		r.tracer = trace.New(os.Stdout)
	}
	http.Handle("/", &templateHandler{filename: "chat.html"})
	http.Handle("/room", r)
	// start the room for clients to connect to
	go r.run()
	//start the web server
	log.Println("Starting web server on port ", *addr)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
