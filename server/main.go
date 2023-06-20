package main

import (
	"flag"

	"evilqr-server/core"
	"evilqr-server/log"
)

var wwwdir = flag.String("d", "", "www content directory path")

func main() {
	flag.Parse()

	if *wwwdir == "" {
		log.Error("you need to set up www content directory path with '-d' parameter")
		return
	}

	http, err := core.NewHttpServer()
	if err != nil {
		log.Error("http: %s", err)
		return
	}
	http.Run(*wwwdir)
}
