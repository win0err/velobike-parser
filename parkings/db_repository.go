package parkings

import (
	"github.com/jinzhu/gorm"
	"time"
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

func (sr *DbStateRepository) FindByTimeRange(from, to time.Time) ([]State, error) {
	var states []State
	if err := sr.DB.
		Set("gorm:auto_preload", true).
		Find(&states, "time >= ? AND time <= ?", from, to).Error;
		err != nil {
		return states, err
	}

	return states, nil
}

func (sr *DbStateRepository) FindAll() ([]State, error) {
	var states []State
	if err := sr.DB.
		Set("gorm:auto_preload", true).
		Find(&states).Error;
		err != nil {
		return states, err
	}

	return states, nil
}
