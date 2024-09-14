package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/will-lol/http-to-aws-proxy/handler"
)

func main() {
	argsWithoutProg := os.Args[1:]
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage of %s: Pass the lambda endpoint as the first argument\n", os.Args[0])
		flag.PrintDefaults()
	}
	if len(argsWithoutProg) == 0 {
		flag.Usage()
		panic("No lambda endpoint provided")
	}

	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	var eventJson events.LambdaFunctionURLRequest
	err = json.Unmarshal(stdin, &eventJson)
	if err != nil {
		panic(err)
	}

	handler := handler.NewLambdaHandler(eventJson, argsWithoutProg[0])

	http.HandleFunc("/", handler.HandlerFunc)
	log.Println("Listening on http://localhost:5544")
	http.ListenAndServe(":5544", nil)
}
