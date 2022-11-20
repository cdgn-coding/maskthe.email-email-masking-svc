package repositories

import (
	"email-masking-svc/src/business/entities"
	"gorm.io/gorm"
)

type PostgresMaskRepository struct {
	db *gorm.DB
}

func (p PostgresMaskRepository) GetByAlias(alias string) (*entities.EmailMask, error) {
	mask := &entities.EmailMask{}
	err := p.db.First(mask, "alias = ?", alias).Error
	if err != nil {
		return nil, err
	}

	return mask, nil
}

func NewPostgresMaskRepository(db *gorm.DB) *PostgresMaskRepository {
	return &PostgresMaskRepository{db: db}
}

func (p PostgresMaskRepository) CreateMask(mask *entities.EmailMask) (*entities.EmailMask, error) {
	result := p.db.Create(mask)
	return mask, result.Error
}
