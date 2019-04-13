package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	hub "github.com/DATA-DOG/golang-websocket-hub"
	"github.com/phogolabs/log"
	"gopkg.in/fsnotify.v1"
)

// Reloader reloads the page
type Reloader struct {
	watcher *fsnotify.Watcher
	server  *hub.Hub
}

// LiveReloader reloads a webpage
func LiveReloader(next http.Handler) http.Handler {
	reloader := NewReloader()
	return reloader.ServeHTTP(next)
}

// NewReloader creates a new reloader
func NewReloader() *Reloader {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}

	watcher.Add(".")

	reloader := &Reloader{
		server:  hub.New(os.Stdout, "*"),
		watcher: watcher,
	}

	go reloader.notify()
	go reloader.server.Run()

	return reloader
}

// ServeHTTP serves the reloader
func (l *Reloader) ServeHTTP(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/livereload.js" {
			l.script(w, r)
			return
		}

		if path == "/livereload" {
			l.server.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (l *Reloader) script(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/javascript")

	buffer := &bytes.Buffer{}

	fmt.Fprintln(buffer, "const liveUrl = 'ws://' + window.location.hostname + ':' + window.location.port + '/livereload';")
	fmt.Fprintln(buffer, "const socket = new WebSocket(liveUrl);")
	fmt.Fprintln(buffer, "socket.addEventListener('message', function (event) {")
	fmt.Fprintln(buffer, "  const message = JSON.parse(event.data);")
	fmt.Fprintln(buffer, "  if (message.topic === 'notify' && message.data === 'reload') {")
	fmt.Fprintln(buffer, "    window.document.location.reload(true);")
	fmt.Fprintln(buffer, "  }")
	fmt.Fprintln(buffer, "});")

	io.Copy(w, buffer)
}

func (l *Reloader) notify() {
	for {
		select {
		case event := <-l.watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				l.reload(event)
			}
		}
	}
}

func (l *Reloader) reload(event fsnotify.Event) {
	msg := &hub.Message{
		Topic: "notify",
		Data:  "reload",
	}

	log.WithField("event", event.Name).Info("reloading")
	l.server.Broadcast <- msg
}
