package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	claudeApi "github.com/mirobeh/learn-claude-api/claude-service"
)

func main() {
	fmt.Print("Enter your question: ")
	question := getTextInput()

	response, err := claudeApi.CallClaudeApi(anthropic.MessageParamRoleUser, question)

	if err != nil {
		panic(err)
	}

	fmt.Print(response)
}

func getTextInput() string {
	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	text = strings.TrimSpace(text)

	return text
}
