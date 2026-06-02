package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
)

type aiResponse struct {
	IsSpam bool 	`json:"is_spam"`
	Score  int 		`json:"score"`
	Reason string 	`json:"reason"`
}

func classify(client *anthropic.Client, msg Message) (SpamResult, error) {
	if msg.Content == "" {
		return SpamResult{}, fmt.Errorf("content cannot be empty")
	}

	if msg.ContentType != ContentTypeText && msg.ContentType != ContentTypeEmail {
		return SpamResult{}, fmt.Errorf("Unknown content type: %s", msg.ContentType)
	}

	prompt := buildPrompt(msg)

	response, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 256,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return SpamResult{}, fmt.Errorf("anthropic API error: %w", err)
	}

	rawText := cleanJSON(response.Content[0].Text)
	
	var aiResp aiResponse
	if err := json.Unmarshal([]byte(rawText), &aiResp); err != nil {
		return SpamResult{}, fmt.Errorf("failed to parse AI response %w", err)
	}

	return SpamResult{
		IsSpam: aiResp.IsSpam,
		Score: aiResp.Score,
		Reason: aiResp.Reason,
		ContentType: msg.ContentType,
	}, nil
}

func buildPrompt(msg Message) string {
	contentLabel := "text message"
	if msg.ContentType == ContentTypeEmail {
		contentLabel = "email"
	}

	return fmt.Sprintf(`You are a spam classifier. Analyze the following %s and respond ONLY with a JSON object, no markdown, no explanation outside the JSON.
		
	JSON format:
		{
		"is_spam": true or false,
		"score": a number from 0 to 100 where 0 is definitely not spam and 100 is definitely spam,
		"reason": "one sentence explanation"
		}
	%s to classify: %s`, 
		contentLabel, contentLabel, msg.Content)
}

func cleanJSON(s string) string {
    s = strings.TrimSpace(s)
    s = strings.TrimPrefix(s, "```json")
    s = strings.TrimPrefix(s, "```")
    s = strings.TrimSuffix(s, "```")
    return strings.TrimSpace(s)
}