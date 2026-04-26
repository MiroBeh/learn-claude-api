package claudeApi

import (
	"context"
	"os"
	"sync"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/joho/godotenv"
)

var (
	clientOnce sync.Once
	clientInst *anthropic.Client
	clientErr  error
)

func getAnthropicClient() (*anthropic.Client, error) {
	clientOnce.Do(func() {
		if err := godotenv.Load(); err != nil {
			clientErr = err
			return
		}

		apiKey := os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			clientErr = os.ErrInvalid
			return
		}

		client := anthropic.NewClient(
			option.WithAPIKey(apiKey),
		)

		clientInst = &client
	})

	if clientErr != nil {
		return nil, clientErr
	}

	return clientInst, nil
}

func CallClaudeApi(role anthropic.MessageParamRole, content string) (string, error) {
	client, err := getAnthropicClient()
	if err != nil {
		return "", err
	}

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages: []anthropic.MessageParam{
			{Role: role, Content: []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(content)}},
		},
		Model: anthropic.ModelClaudeHaiku4_5,
	})
	if err != nil {
		return "", err
	}

	return message.Content[0].Text, nil
}

func CallClaudeApiWithHistory(messages []anthropic.MessageParam) (string, error) {
	client, err := getAnthropicClient()
	if err != nil {
		return "", err
	}

	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		MaxTokens: 1024,
		Messages:  messages,
		Model:     anthropic.ModelClaudeHaiku4_5,
	})
	if err != nil {
		return "", err
	}

	return message.Content[0].Text, nil
}
