package main

import (
	"context"
	"log"
	"os"

	"goast/core"

	"github.com/sashabaranov/go-openai"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	params := core.InitializeParameters{
		Model:  "gpt-4-1106-preview",
		Client: client,
		// Temperature: 0.7,
	}
	codeGenerator, err := core.NewCodeGenerator(params)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = codeGenerator.LoadContext(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	generated, err := codeGenerator.Predict(ctx, "Please Generate the controller layer, usecase layer, model layer, dao layer codes for the API file.")
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(generated)
}
