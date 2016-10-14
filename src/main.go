package main 

import (
	"flag"
	"go/build"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"chat"
	"fmt"
)

var (
	addr      = flag.String("addr", ":8080", "http service address")
	assets    = flag.String("assets", defaultAssetPath(), "path to assets")
	homeTempl *template.Template
)

func defaultAssetPath() string {
	p, err := build.Default.Import("gary.burd.info/go-websocket-chat", "", build.FindOnly)
	if err != nil {
		return "."
	}

	return p.Dir
}

func homeHandler(resp http.ResponseWriter, req *http.Request) {
	homeTempl.Execute(resp, req.Host)
}

func main() {
	fmt.Println(*addr);
	h := chat.H
	flag.Parse()
	homeTempl = template.Must(template.ParseFiles(filepath.Join("/Users/henry-sun/data/www/my-go-chat/src/chat/home.html")))
	go h.Run()
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", chat.WsHandler)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}