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
	fmt.Print("Hello I am your bot, ask me anything: ")

	var messages []anthropic.MessageParam

	for {
		userInput := getTextInput()
		messages = append(messages, anthropic.MessageParam{Role: anthropic.MessageParamRoleUser, Content: []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(userInput)}})

		response, err := claudeApi.CallClaudeApiWithHistory(messages)

		messages = append(messages, anthropic.MessageParam{Role: anthropic.MessageParamRoleAssistant, Content: []anthropic.ContentBlockParamUnion{anthropic.NewTextBlock(response)}})

		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", response)
	}

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
