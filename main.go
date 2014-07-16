package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/lestrrat/go-server-starter-listener"
	"github.com/mix3/go-rocket-sample-app/webapp"
)

func main() {
	log.Println("Launch succeeded!")

	listener, _ := ss.NewListener()
	if listener == nil {
		// Fallback if not running under Server::Starter
		var err error
		listener, err = net.Listen(
			"tcp",
			fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")),
		)
		if err != nil {
			panic(fmt.Errorf("Failed to listen to port %s", os.Getenv("PORT")))
		}
	}

	webapp.Start(listener)
}
