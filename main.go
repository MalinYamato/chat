// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"
	//"github.com/kabukky/httpscerts"

)

var addr = flag.String("addr", ":8080", "http service address")
var homeTemplate = template.Must(template.ParseFiles("home.html"))

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	homeTemplate.Execute(w, r.Host)
}

func main() {

	queue := new(QueueStack)

	// Check if the cert files are available.
	//err := httpscerts.Check("cert.pem", "key.pem")
	// If they are not available, generate new ones.
	//if err != nil {
	//	err = httpscerts.Generate("cert.pem", "key.pem", "localhost")
	//	if err != nil {
	//		log.Fatal("Error: Couldn't create https certs.")
	//	}
	//}

	flag.Parse()
	hub := newHub(*queue)
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServeTLS(*addr,"cert.pem", "key.pem", nil )
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
