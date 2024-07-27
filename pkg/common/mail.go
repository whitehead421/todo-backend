package common

import (
	"context"
	"fmt"

	"github.com/mailjet/mailjet-apiv3-go"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func SendActivationMail(reader *kafka.Reader, ctx context.Context) {
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			zap.L().Error("Failed to read message from Kafka", zap.Error(err))
			continue
		}

		err = sendEmail(string(msg.Key), string(msg.Value))
		if err != nil {
			zap.L().Error("Failed to send activation email", zap.Error(err))
			continue
		}

		zap.L().Info(
			"Sent activation email",
			zap.String("email", string(msg.Key)),
		)
	}
}

func sendEmail(toEmail, token string) error {
	env := GetEnvironmentVariables()
	mailjetClient := mailjet.NewMailjetClient(env.MailjetAPIKey, env.MailjetSecretKey)

	email := &mailjet.InfoSendMail{
		FromEmail: env.SenderEmail,
		FromName:  "Todo App",
		Subject:   "Verify Your Account",
		TextPart:  fmt.Sprintf("Please use the following token to activate your account: %s", token),
		HTMLPart:  fmt.Sprintf("<h3>Please use the following link to activate your account:</h3><a target='_blank' href='http://%s:%s/verify?token=%s'>Activate</a>", env.ApplicationHost, env.AuthPort, token),
		Recipients: []mailjet.Recipient{
			{Email: toEmail},
		},
	}

	_, err := mailjetClient.SendMail(email)
	if err != nil {
		return err
	}

	return nil
}
