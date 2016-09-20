package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/oremj/aws-elasticsearch-proxy/proxy"
)

var addr = flag.String("addr", "127.0.0.1:8080", "Listen address")
var esEndpoint = flag.String("es-endpoint", "", "ElasticSearch endpoint e.g., my-es.us-east-1.es.amazonaws.com")
var region = flag.String("region", "us-east-1", "AWS region")

func main() {
	flag.Parse()
	if *esEndpoint == "" {
		fmt.Println("-es-endpoint must be defined\n")
		flag.Usage()
		os.Exit(2)
	}
	signer := v4.NewSigner(defaults.Get().Config.Credentials)
	es := proxy.NewElasticSearch(*esEndpoint, *region, signer)
	log.Fatal(http.ListenAndServe(*addr, es))
}
