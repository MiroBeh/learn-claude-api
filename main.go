package main

import (
	"fmt"

	claudeApi "github.com/mirobeh/learn-claude-api/claude-service"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
	translation, err := translate("Brücke", "englisch")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Translation: %s\n", translation)
}

func translate(word string, language string) (string, error) {
	response, err := claudeApi.CallClaudeApi(anthropic.MessageParamRoleUser, fmt.Sprintf("Translate the word '%s' to %s.", word, language))
	if err != nil {
		return "", err
	}
	return response, nil
}
