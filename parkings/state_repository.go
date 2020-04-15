package parkings

import "github.com/jinzhu/gorm"

type StateRepository struct {
	DB *gorm.DB
}

func ProvideStateRepository(DB *gorm.DB) StateRepository {
	return StateRepository{DB: DB}
}

func (sr *StateRepository) GetLast() State {
	var state State
	sr.DB.Last(&state)

	return state
}

func (sr *StateRepository) Save(state State) State {
	sr.DB.Save(&state)

	return state
}

func (sr *StateRepository) SaveAll(states []State) []State {
	for _, state := range states {
		sr.Save(state)
	}

	return states
}
