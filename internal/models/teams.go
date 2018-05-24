package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// Team is a representation of a user that is stored in the database and an api response
type Team struct {
	ID   string `json:"id" gorm:"column(id);index"`
	Name string `json:"name" gorm:"column(name)" validate:"required"`
	// if parent team id is empty then we must assume that it is a parent.
	ParentTeamID string     `json:"-" gorm:"column(parent_team_id);index"`
	CreatedAt    time.Time  `json:"created_at" gorm:"column(created_at)"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"column(updated_at)"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"column(deleted_at);index"`
}

// BeforeCreate is run before an insert query happens. So we can set some defaults here.
func (t *Team) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4()

	scope.SetColumn("id", id.String())
	scope.SetColumn("created_at", time.Now())
	scope.SetColumn("updated_at", time.Now())

	return nil
}

func NewTeamModel(db *gorm.DB) *TeamModel {
	return &TeamModel{db}
}

type TeamModel struct {
	db *gorm.DB
}

// GetAll will get all teams under the given parent team id.
func (m *TeamModel) GetAll(parentTeamID string) ([]Team, error) {
	var teams []Team
	err := m.db.Where("parent_team_id = ?", parentTeamID).Find(&teams).Error

	return teams, err
}

// GetByID will get a team by the given id.
func (m *TeamModel) GetByID(id string) (Team, error) {
	var team Team
	err := m.db.Where("id = ?", id).First(&team).Error

	return team, err
}

// Create will create the team in the database.
func (m *TeamModel) Create(t Team) (Team, error) {
	err := m.db.Create(&t).Error
	return t, err
}

// Update will update the team in the database.
func (m *TeamModel) Update(t Team) (Team, error) {
	err := m.db.Save(&t).Error

	return t, err
}

// Delete will HARD delete the team from the database.
func (m *TeamModel) Delete(id string) error {
	return m.db.Delete(Team{ID: id}).Error
}
