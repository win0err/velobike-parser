package parkings

import "github.com/jinzhu/gorm"

type StateRepository struct {
	DB *gorm.DB
}

func ProvideStateRepostiory(DB *gorm.DB) StateRepository {
	return StateRepository{DB: DB}
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
