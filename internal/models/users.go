package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User is a representation of a user that is stored in the database and an api response
type User struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	TeamID    string     `json:"team_id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// BeforeCreate is run before an insert query happens. So we can set some defaults here.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.NewV4()

	scope.SetColumn("id", id.String())
	scope.SetColumn("created_at", time.Now())
	scope.SetColumn("updated_at", time.Now())

	return nil
}

func NewUserModel(db *gorm.DB) *UserModel {
	return &UserModel{db}
}

type UserModel struct {
	db *gorm.DB
}

// GetAll will get all teams under the given parent User id.
func (m *UserModel) GetAll(teamID string) ([]User, error) {
	var users []User
	err := m.db.Where("team_id = ?", teamID).Find(&users).Error

	return users, err
}

// GetByID will get a User by the given id.
func (m *UserModel) GetByID(id string) (User, error) {
	var user User
	err := m.db.Where("id = ?", id).First(&user).Error

	return user, err
}

// Create will create the User in the database.
func (m *UserModel) Create(u User) (User, error) {
	err := m.db.Create(&u).Error
	return u, err
}

// Update will update the User in the database.
func (m *UserModel) Update(u User) (User, error) {
	err := m.db.Save(&u).Error

	return u, err
}

// Delete will HARD delete the User from the database.
func (m *UserModel) Delete(id string) error {
	return m.db.Delete(User{ID: id}).Error
}
