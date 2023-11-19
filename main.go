package main

import (
	"context"
	"log"
	"os"

	"goast/core"

	"github.com/sashabaranov/go-openai"
)

// func makeMessages(apiFile string, targetFile string, referenceFiles []string, moderation string) ([]openai.ChatCompletionMessage, error) {
// 	var messages []openai.ChatCompletionMessage
// 	if len(referenceFiles) != 0 {
// 		messages = append(messages, openai.ChatCompletionMessage{
// 			Role:    "system",
// 			Content: "These are the code which already exist in the directory and have relation with the file you are going to generate. ",
// 		})
// 	}
// 	for _, referenceFile := range referenceFiles {
// 		err := addFileContent(&messages, referenceFile)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	messages = append(messages, openai.ChatCompletionMessage{
// 		Role:    "user",
// 		Content: "This is the API file you are going to generate with. Generate the code witch implements the each functions by detail for all APIs defined in the API file. To be consistent with existing files and with the number and type of function arguments and return values.\n",
// 	})

// 	if moderation != "" {
// 		messages = append(messages, openai.ChatCompletionMessage{
// 			Role:    "user",
// 			Content: moderation,
// 		})
// 	}

// 	err := addFileContent(&messages, apiFile)
// 	if err != nil {
// 		return nil, err
// 	}

// 	messages = append(messages, openai.ChatCompletionMessage{
// 		Role:    "assistant",
// 		Content: targetFile,
// 	})

// 	return messages, nil
// }

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	client := openai.NewClient(apiKey)
	ctx := context.Background()

	params := core.InitializeParameters{
		Model:  "gpt-4-turbo",
		Client: client,
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

	generated, err := codeGenerator.Predict(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(generated)
}
