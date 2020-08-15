package client

import (
	"context"

	"github.com/machinebox/graphql"
)

// SendToOPS send a datastring to the OPS SITP endpoint
func SendToOPS(url string, token string, data string) {

	// create a client (safe to share across requests)
	client := graphql.NewClient(url)

	// define a Context for the request
	// Ref: https://golang.org/pkg/context/
	ctx := context.Background()

	req := graphql.NewRequest(`
		mutation tokenAuth($username: String!, $password: String!) {
				tokenAuth(username: $username, password: $password) {
					token
				}
			}
	`)

	req.Var("username", "cisco")
	req.Var("password", "ciscocisco")

	// run it and capture the response
	var respData map[string]map[string]interface{}
	client.Run(ctx, req, &respData)
	if err := client.Run(ctx, req, &respData); err != nil {
		// log.Fatal(err)
	}

	JWT := respData["tokenAuth"]["token"]

	// make a request
	req = graphql.NewRequest(`
		mutation publishPipelineData(
			$token: String!
			$url: String!
			$values: GenericScalar!
		) {
			opsPipelineFormPage(token: $token, url: $url, values: $values) {
				result
			}
		}
	`)

	type valuesStruct struct {
		PipelineToken string `json:"pipeline_token"`
		RawData       string `json:"raw_data"`
	}

	values := valuesStruct{token, data}

	// set any variables
	req.Var("token", JWT)
	req.Var("url", "/pipeline/")
	req.Var("values", values)

	// run it and capture the response
	client.Run(ctx, req, &respData)
	if err := client.Run(ctx, req, &respData); err != nil {
		// log.Fatal(err)
	}
}

// ResponseStruct not defined yet
// map[opsPipelineFormPage:map[result:FAIL]]
type ResponseStruct struct{}
