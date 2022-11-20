package services

import (
	"email-masking-svc/src/application/events"
	"email-masking-svc/src/business/repositories"
	"encoding/json"
	"fmt"
)

type RedirectEmail struct {
	maskRepository     repositories.MaskRepository
	sendEmailPublisher events.Publisher
}

func NewRedirectEmail(maskRepository repositories.MaskRepository, sendEmailPublisher events.Publisher) *RedirectEmail {
	return &RedirectEmail{maskRepository: maskRepository, sendEmailPublisher: sendEmailPublisher}
}

func (r RedirectEmail) Execute(command *events.ReceivedEmail) error {
	mask, err := r.maskRepository.GetByAlias(command.To)
	if err != nil {
		return err
	}

	subject := fmt.Sprintf("Message From: %s. %s", command.From, command.Subject)
	message := events.SendEmail{
		Subject:     subject,
		From:        mask.Alias,
		To:          mask.Target,
		PlainText:   command.PlainText,
		HTML:        command.HTML,
		Attachments: command.Attachments,
	}

	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = r.sendEmailPublisher.Dispatch(string(messageString))
	if err != nil {
		return err
	}

	return nil
}
