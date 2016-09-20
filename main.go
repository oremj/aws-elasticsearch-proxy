package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/oremj/aws-elasticsearch-proxy/proxy"
)

var addr = flag.String("addr", ":8000", "Listen address")
var esEndpoint = flag.String("es-endpoint", "", "ElasticSearch endpoint e.g., my-es.us-east-1.es.amazonaws.com")
var region = flag.String("region", "us-east-1", "AWS region")
var debug = flag.Bool("debug", false, "Enables additional logging")

func main() {
	flag.Parse()
	if *esEndpoint == "" {
		fmt.Println("-es-endpoint must be defined\n")
		flag.Usage()
		os.Exit(2)
	}
	signer := v4.NewSigner(defaults.Get().Config.Credentials)

	if *debug {
		signer.Debug = aws.LogDebugWithSigning
		signer.Logger = aws.NewDefaultLogger()
	}

	es := proxy.NewElasticSearch(*esEndpoint, *region, signer)
	log.Fatal(http.ListenAndServe(*addr, es))
}
