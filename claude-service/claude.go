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

type RequestOptions struct {
	MaxTokens    int64
	StopSequence string
	SystemPrompt string
}

func resolveRequestOptions(options []RequestOptions) (int64, []string, string) {
	maxTokens := int64(defaultMaxTokens)
	var stopSequences []string
	var systemPrompt string

	if len(options) > 0 {
		if options[0].MaxTokens > 0 {
			maxTokens = options[0].MaxTokens
		}

		if options[0].StopSequence != "" {
			stopSequences = []string{options[0].StopSequence}
		}

		if options[0].SystemPrompt != "" {
			systemPrompt = options[0].SystemPrompt
		}
	}

	return maxTokens, stopSequences, systemPrompt
}

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

func callClaudeApiInternal(ctx context.Context, messages []anthropic.MessageParam, maxTokens int64, stopSequences []string, systemPrompt string) (string, error) {
	client, err := getAnthropicClient()
	if err != nil {
		return "", err
	}

	params := anthropic.MessageNewParams{
		MaxTokens: maxTokens,
		Messages:  messages,
		Model:     anthropic.ModelClaudeHaiku4_5,
	}

	if len(stopSequences) > 0 {
		params.StopSequences = stopSequences
	}

	if systemPrompt != "" {
		params.System = []anthropic.TextBlockParam{{Text: systemPrompt}}
	}

	message, err := client.Messages.New(ctx, params)
	if err != nil {
		return "", err
	}

	if len(message.Content) == 0 {
		return "", errors.New("empty response content from API")
	}

	return message.Content[0].Text, nil
}

func CallClaudeApi(role anthropic.MessageParamRole, content string, options ...RequestOptions) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	maxTokens, stopSequences, systemPrompt := resolveRequestOptions(options)

	messages := []anthropic.MessageParam{
		{Role: role, Content: []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(content)}},
	}

	return callClaudeApiInternal(ctx, messages, maxTokens, stopSequences, systemPrompt)
}

func CallClaudeApiWithHistory(messages []anthropic.MessageParam, options ...RequestOptions) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	maxTokens, stopSequences, systemPrompt := resolveRequestOptions(options)

	return callClaudeApiInternal(ctx, messages, maxTokens, stopSequences, systemPrompt)
}
