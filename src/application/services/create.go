package services

import (
	"email-masking-svc/src/business/entities"
	"email-masking-svc/src/business/repositories"
)

type CreateMaskCommand struct {
	Alias  string
	Target string
}

type CreateMask struct {
	maskRepository repositories.MaskRepository
}

func NewCreateMask(maskRepository repositories.MaskRepository) *CreateMask {
	return &CreateMask{maskRepository: maskRepository}
}

func (c CreateMask) Execute(command CreateMaskCommand) error {
	_, err := c.maskRepository.CreateMask(&entities.EmailMask{
		Alias:  command.Alias,
		Target: command.Target,
	})

	if err != nil {
		return err
	}

	return nil
}
