package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/frostyslav/lseg-demo/app/webserver"
	"github.com/sirupsen/logrus"
)

var (
	port *int
)

func init() {
	l := logrus.New()
	l.Formatter = &logrus.TextFormatter{FullTimestamp: true, ForceColors: true}
	l.Level = logrus.DebugLevel
	l.Out = os.Stdout

	log.SetOutput(l.Writer())

	port = flag.Int("port", 8080, "port number")
}

func main() {
	flag.Parse()

	log.Printf("Listening on localhost:%d\n", *port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), webserver.Serve))
}
