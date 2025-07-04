package main

import (
	"ai_code_gen/internals/agents"
	"context"
	"flag"
	"fmt"
	"os"
)

func main() {
	openAPIKey := flag.String("openai-key", "", "OpenAI API Key")
	outputDir := flag.String("output-dir", "./output", "Output directory for generated files")
	basePackage := flag.String("base-package", "github.com/user/app", "Base package for generated files")
	workerCount := flag.Int("worker-count", 4, "Number of workers to use for file generation")

	flag.Parse()
	if *openAPIKey == "" {
		*openAPIKey = "{OPENAI_API_SECRET_KEY}"
		if *openAPIKey == "" {
			fmt.Println("Please provide with OpenAI API Key")
			os.Exit(1)
		}
	}

	cntxt := context.Background()

	openAIClient := agents.NewOpenAPI(cntxt, *openAPIKey, nil)
	agents.NewAgent(cntxt, openAIClient, *outputDir, *basePackage, *workerCount)
	fmt.Println("READY")
}
