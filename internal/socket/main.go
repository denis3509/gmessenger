// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func ServeHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

// func main() {
// 	flag.Parse()
// 	hub := NewHub(log)
// 	go hub.Run()
// 	http.HandleFunc("/", ServeHome)
// 	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
// 		ServeWs(hub, w, r)
// 	})
// 	err := http.ListenAndServe(*addr, nil)
// 	if err != nil {
// 		log.Fatal("ListenAndServe: ", err)
// 	}
// }
