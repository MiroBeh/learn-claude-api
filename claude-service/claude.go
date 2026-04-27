package claudeApi

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/joho/godotenv"
)

const (
	defaultMaxTokens = 1024
	requestTimeout   = 30 * time.Second
)

var (
	clientOnce sync.Once
	clientInst anthropic.Client
	clientErr  error
)

func getAnthropicClient() (anthropic.Client, error) {
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

		clientInst = anthropic.NewClient(
			option.WithAPIKey(apiKey),
		)
	})

	if clientErr != nil {
		return anthropic.Client{}, clientErr
	}

	return clientInst, nil
}

func callClaudeApiInternal(ctx context.Context, messages []anthropic.MessageParam) (string, error) {
	client, err := getAnthropicClient()
	if err != nil {
		return "", err
	}

	message, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		MaxTokens: defaultMaxTokens,
		Messages:  messages,
		Model:     anthropic.ModelClaudeHaiku4_5,
	})
	if err != nil {
		return "", err
	}

	if len(message.Content) == 0 {
		return "", errors.New("empty response content from API")
	}

	return message.Content[0].Text, nil
}

func CallClaudeApi(role anthropic.MessageParamRole, content string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	messages := []anthropic.MessageParam{
		{Role: role, Content: []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(content)}},
	}

	return callClaudeApiInternal(ctx, messages)
}

func CallClaudeApiWithHistory(messages []anthropic.MessageParam) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	return callClaudeApiInternal(ctx, messages)
}
