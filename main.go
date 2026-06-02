package main

import (
	"fmt"
	"log"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	client := anthropic.NewClient()
messages := []Message{
    {
        Content:     "Congratulations! You won a free iPhone. Click here now!!!",
        ContentType: ContentTypeText,
    },
    {
        Content:     "Hey, are we still on for lunch tomorrow at 12?",
        ContentType: ContentTypeText,
    },
    {
        Content:     "URGENT: Your bank account has been compromised. Verify your details immediately at secure-bank-login.ru",
        ContentType: ContentTypeText,
    },
    {
        Content:     "Don't forget to pick up milk on your way home",
        ContentType: ContentTypeText,
    },
    {
        Content:     "You have been selected for a $1000 gift card. Reply YES to claim. Only 3 left!",
        ContentType: ContentTypeText,
    },

    {
        Content:     "Hi Viktor, just following up on our meeting notes from Tuesday. Let me know if you have any questions.",
        ContentType: ContentTypeEmail,
    },
    {
        Content:     "Dear Customer, your invoice #4521 for $240 is due on June 15th. Please log in to your account to pay.",
        ContentType: ContentTypeEmail,
    },
    {
        Content:     "We noticed unusual activity on your Netflix account. Click here to verify your identity or your account will be suspended.",
        ContentType: ContentTypeEmail,
    },
    {
        Content:     "Hey! Loved your talk at the conference last week. Would love to grab a coffee and chat more about the project.",
        ContentType: ContentTypeEmail,
    },
    {
        Content:     "Earn $5000 a week working from home! No experience needed. Limited spots available. Act now at workfromhome-cash.net",
        ContentType: ContentTypeEmail,
    },
}

	for _, msg := range messages {
		result, err := classify(&client, msg)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("--- %s ---\nIs spam: %v\nScore: %d/100\nReason: %s\n\n",
			msg.ContentType, result.IsSpam, result.Score, result.Reason)
	}

	_ = os.Getenv

}