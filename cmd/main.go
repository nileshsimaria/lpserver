package main

import (
	"flag"

	lp "github.com/nileshsimaria/lpserver"
)

var (
	host = flag.String("host", "127.0.0.1", "host name or ip")
	port = flag.Int("port", 50052, "grpc server port")
)

func main() {
	flag.Parse()
	s := lp.NewLPServer(*host, *port)
	s.StartServer()
}
