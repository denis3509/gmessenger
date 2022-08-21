package main

import (
	"flag"
	"fmt"
	"log"
	"messenger/internal/test"
)

func help() {
	fmt.Println("recreate-test-db - пересоздать тестовую базу данных")
	flag.PrintDefaults()
}

//  (https://github.com/urfave/cli) - todo cli framework

var command string

func main() {
	flag.Parse()
	args := flag.Args()
	// args = append(args, "recreate-test-db")
	fmt.Println(args)
 
	if len(args) == 0 {
		help()
		log.Fatal("no args specified")
	}
    command = args[0]
	 
	switch command {
	case "recreate-test-db":
		test.DropAndCreateTestDB()
	default:
		help()
	}
}
