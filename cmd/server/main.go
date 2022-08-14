package main

import (
	"fmt"
	"log"
	"messenger/internal/config"
	"os"
	"os/signal"
	"syscall"

	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/lib/pq"
)

func main(){
	cfg := config.GetConfig()
	db, err := dbx.MustOpen("postgres", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()  
    sigs := make(chan os.Signal, 1) 
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM) 
    done := make(chan bool, 1) 
    go func() {
        sig := <-sigs 
        fmt.Printf("Received signal: %s \n",sig)
        done <- true
    }()
 
    <-done
    fmt.Println("exiting")
 
}