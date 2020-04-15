package parkings

import (
	"github.com/jinzhu/gorm"
)

type DbStateRepository struct {
	DB *gorm.DB
}

func ProvideDbStateRepository(DB *gorm.DB) DbStateRepository {
	return DbStateRepository{DB: DB}
}

func (sr *DbStateRepository) GetLast() (State, error) {
	var state State
	if err := sr.DB.Last(&state).Error; err != nil {
		return state, err
	}

	return state, nil
}

func (sr *DbStateRepository) Save(state State) error {
	return sr.DB.Save(&state).Error
}

func (sr *DbStateRepository) SaveAll(states []State) error {
	return sr.DB.Transaction(func(tx *gorm.DB) error {
		for _, state := range states {
			if err := tx.Create(&state).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
