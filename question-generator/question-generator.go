package main

import (
	"errors"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	claudeApi "github.com/mirobeh/learn-claude-api/claude-service"
)

func main() {
	topic := "Werder Bremen"
	numQuestions := 3

	genQuestions, err := generateQuestions(topic, numQuestions)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", genQuestions)
}

func generateQuestions(topic string, numQuestions int) (string, error) {
	systemPrompt := fmt.Sprintf("You are an expert on %s. Generate thought-provoking questions about this topic.", topic)
	promt := fmt.Sprintf("Generate %d questions about %s as a numbered list.", numQuestions, topic)

	options := claudeApi.RequestOptions{
		MaxTokens:    500,
		StopSequence: fmt.Sprintf("%d", numQuestions+1),
		SystemPrompt: systemPrompt,
	}

	output, err := claudeApi.CallClaudeApi(anthropic.MessageParamRoleUser, promt, options)

	if err != nil {
		return "", errors.New("Problem calling the claude api")
	}

	return output, nil
}
