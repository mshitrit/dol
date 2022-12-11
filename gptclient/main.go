package main

import (
	"fmt"

	"gptchatapi/gpt"
)

func main() {
	//create API_KEY at: https://beta.openai.com/account/api-keys
	apiKey := "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"

	gptClient := gpt.NewClient(apiKey)

	text, _ := gptClient.SendRequest("Who is the strongest avenger ?")
	fmt.Println(text)

}
