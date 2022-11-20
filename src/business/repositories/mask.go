package repositories

import "email-masking-svc/src/business/entities"

type MaskRepository interface {
	CreateMask(mask *entities.EmailMask) (*entities.EmailMask, error)
	GetByAlias(alias string) (*entities.EmailMask, error)
}
