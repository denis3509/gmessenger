package main

import (
	"fmt"
	"messenger/pkg/log"
	"messenger/internal/config"
	"messenger/internal/socket"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

func main(){
	cfg := config.GetConfig()

	log := log.NewLogger(os.Stdout)

	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// http 
	http.HandleFunc("/", socket.ServeHome)  
	// socket  
	hub := socket.NewHub(log)
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Info("serve ws to ",r.RemoteAddr)
		socket.ServeWs(hub, w, r)
	})
	go func () {
		err  = http.ListenAndServe(":" + strconv.Itoa(cfg.Port) , nil)
		if err != nil {
			log.Fatal("Can't start server: ", err)
		} 
	}()
	log.Infof("server listening on localhost:%d",cfg.Port)
	// signal 
    sigs := make(chan os.Signal, 1) 
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) 
    done := make(chan bool, 1) 
    go func() {
        sig := <-sigs 
        fmt.Printf("Received signal: %s \n",sig)
		hub.CloseAll()
        done <- true
    }() 
    <-done
    fmt.Println("exiting")
 
}