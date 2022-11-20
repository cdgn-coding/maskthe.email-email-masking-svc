package consumers

import (
	"email-masking-svc/src/application/events"
	"email-masking-svc/src/application/services"
	"email-masking-svc/src/infrastructure/configuration"
	"encoding/json"
	"errors"
	"fmt"
)

type EmailConsumer struct {
	logger        configuration.Logger
	redirectEmail *services.RedirectEmail
}

func NewEmailConsumer(logger configuration.Logger, redirectEmail *services.RedirectEmail) *EmailConsumer {
	return &EmailConsumer{logger: logger, redirectEmail: redirectEmail}
}

var PayloadNotValid = errors.New("payload not valid")

func (e *EmailConsumer) Invoke(payload string) error {
	command := &events.ReceivedEmail{}
	err := json.Unmarshal([]byte(payload), command)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Unable to parse payload. %s", payload))
		return PayloadNotValid
	}

	err = e.redirectEmail.Execute(command)

	if err != nil {
		e.logger.Error(fmt.Sprintf("Unable to redirect email. cause: %v", err))
		return err
	}

	return nil
}
