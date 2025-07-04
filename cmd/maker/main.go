package main

import (
	"ai_code_gen/internals/agents"
	"context"
	"fmt"
)

func main() {
	key := "{OPENAI_API_SECRET_KEY}"

	cntxt := context.Background()

	prompt := "Write a simple todo program"

	openAIClient := agents.NewOpenAPI(cntxt, key, nil)

	res, err := openAIClient.Query("", prompt)

	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v", res)
}
