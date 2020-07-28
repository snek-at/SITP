package client

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/machinebox/graphql"
)

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Fatal(errors.New(key + " not set"))
	return ""
}

// SendToOPS send a datastring to the OPS SITP endpoint
func SendToOPS(data string) {
	fmt.Println("GraphQL Client Received Data ")
	ioutil.WriteFile("tmp.json", []byte(data), 0644)
	fmt.Print(data)

	// assign env vars
	OPSURL := getEnv("OPS_URL")
	OPSTOKEN := getEnv("OPS_TOKEN")

	// create a client (safe to share across requests)
	client := graphql.NewClient(OPSURL)

	// make a request
	req := graphql.NewRequest(`
    query ($token: String!, $key: String!) {
        items (id:$key) {
            field1
            field2
            field3
        }
    }
	`)

	// set any variables
	req.Var("token", OPSTOKEN)

	// set header fields
	req.Header.Set("Authorization", OPSTOKEN)

	// define a Context for the request
	// Ref: https://golang.org/pkg/context/
	ctx := context.Background()

	// run it and capture the response
	var respData ResponseStruct
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
}

// ResponseStruct not defined yet
type ResponseStruct struct{}
