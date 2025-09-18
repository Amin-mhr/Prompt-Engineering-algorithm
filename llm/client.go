package llm

import (
	"context"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type Usage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// model: مثلا "gpt-4o-mini"  | temperature: 0.0..2.0
func CallLLM(prompt string, model string, temperature float32) (string, Usage, error) {
	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	req := openai.ChatCompletionRequest{
		Model: model,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleUser, Content: prompt},
		},
		Temperature: temperature,
	}

	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return "", Usage{}, err
	}

	usage := Usage{
		PromptTokens:     resp.Usage.PromptTokens,
		CompletionTokens: resp.Usage.CompletionTokens,
		TotalTokens:      resp.Usage.TotalTokens,
	}

	content := ""
	if len(resp.Choices) > 0 {
		content = resp.Choices[0].Message.Content
	}
	return content, usage, nil
}
