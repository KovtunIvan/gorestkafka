package main

import (
	"flag"
	"gorestkafka/pkg/api/handlers"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

var (
	addr     = flag.String("addr", ":8081", "TCP address to listen to")
	compress = flag.Bool("compress", false, "Response compression")
)

func main() {
	flag.Parse()

	pH := handlers.ProduceHandler
	cH := handlers.ConsumeHandler

	if *compress {
		pH = fasthttp.CompressHandler(pH)
	}

	r := router.New()
	r.POST("/produce/{topic}", pH)
	r.GET("/consume/{topic}", cH)

	if err := fasthttp.ListenAndServe(*addr, r.Handler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}

}
