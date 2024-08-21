package lbserver

import (
	"flag"
	"go-load-balancer/pkg/lb"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func StartServer() {
	var serverList string
	var port int

	flag.StringVar(&serverList, "backends", "", "Load balanced backends, use commas to separate")
	flag.IntVar(&port, "port", 3030, "Port to serve")
	flag.Parse()

	if len(serverList) == 0 {
		log.Fatal("Please provide one or more backends to load balance")
	}

	servers := strings.Split(serverList, ",")
	for _, server := range servers {
		serverUrl, error := url.Parse(server)
		if error != nil {
			log.Fatal(error)
		}

		lb.RegisterServerWithReverseProxy(serverUrl)
	}

	server := http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: http.HandlerFunc(lb.BalanceLoad),
	}

	log.Printf("Load Balancer started at :%d\n", port)
	if error := server.ListenAndServe(); error != nil {
		log.Fatal(error)
	}
}
