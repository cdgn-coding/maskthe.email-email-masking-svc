package services

import (
	"email-masking-svc/src/application/events"
	"email-masking-svc/src/business/repositories"
	"email-masking-svc/src/infrastructure/configuration"
	"encoding/json"
)

type RedirectEmail struct {
	maskRepository     repositories.MaskRepository
	sendEmailPublisher events.Publisher
	logger             configuration.Logger
}

func NewRedirectEmail(maskRepository repositories.MaskRepository, sendEmailPublisher events.Publisher, logger configuration.Logger) *RedirectEmail {
	return &RedirectEmail{maskRepository: maskRepository, sendEmailPublisher: sendEmailPublisher, logger: logger}
}

func (r RedirectEmail) Execute(command *events.ReceivedEmail) error {
	r.logger.Info("Redirecting email.")
	mask, err := r.maskRepository.GetByAlias(command.To)
	if err != nil {
		r.logger.Info("Error fetching mask. ", err)
		return err
	}

	message := events.SendEmail{
		Subject:     command.Subject,
		From:        "noreply@maskthe.email",
		To:          mask.Target,
		PlainText:   command.PlainText,
		HTML:        command.HTML,
		Attachments: command.Attachments,
	}

	messageString, err := json.Marshal(message)
	if err != nil {
		return err
	}

	r.logger.Info("Publishing email to send.", string(messageString))
	err = r.sendEmailPublisher.Dispatch(string(messageString))
	if err != nil {
		return err
	}

	return nil
}
