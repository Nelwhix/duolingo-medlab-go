package mailer

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SESMailer struct {
	client *sesv2.Client
	from   string
}

func NewSESMailer(client *sesv2.Client, from string) *SESMailer {
	return &SESMailer{client: client, from: from}
}

func (m *SESMailer) SendPasswordResetEmail(ctx context.Context, to, token string) error {
	resetUrl := fmt.Sprintf("%s?token=%s", os.Getenv("APP_URL"), token)
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(m.from),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String("Reset your password"),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(fmt.Sprintf("<p>Click <a href=\"%s\">here</a> to reset your password.</p>", resetUrl)),
					},
					Text: &types.Content{
						Data: aws.String(fmt.Sprintf("Reset your password by visiting: %s", resetUrl)),
					},
				},
			},
		},
	}

	_, err := m.client.SendEmail(ctx, input)
	return err
}
