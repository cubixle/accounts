package api

import (
	"net/http"

	"github.com/cubixle/accounts/internal/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

// ListUsersHandler will list all Users related to the authed user.
func (a *API) ListUsersHandler(c echo.Context) error {
	UserModel := models.NewUserModel(a.appContext.DB)

	authedUser := a.getAuthedUser(c)

	Users, err := UserModel.GetAll(authedUser.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		a.appContext.Logger.WithError(err).Error("unable to select users from the database")
	}

	return c.JSON(http.StatusOK, Users)
}

// GetUserHandler will get a single User for the given id.
func (a *API) GetUserHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid User id")
	}

	UserModel := models.NewUserModel(a.appContext.DB)
	User, err := UserModel.GetByID(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to access resource")
	}

	return c.JSON(http.StatusOK, User)
}

// CreateUserHandler will create a new User within the database.
func (a *API) CreateUserHandler(c echo.Context) error {
	var t models.User
	if err := c.Bind(&t); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m := models.NewUserModel(a.appContext.DB)
	User, err := m.Create(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create User")
	}

	return c.JSON(http.StatusCreated, User)
}

func (a *API) UpdateUserHandler(c echo.Context) error {
	var t models.User
	if err := c.Bind(&t); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	t.ID = c.Param("id")
	m := models.NewUserModel(a.appContext.DB)

	// check entity exists
	if _, err := m.GetByID(t.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "unable to find User")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to update User")
	}

	user, err := m.Update(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to update User")
	}

	return c.JSON(http.StatusAccepted, user)
}

// DeleteUser will delete the User from the database.
func (a *API) DeleteUserHandler(c echo.Context) error {
	id := c.Param("id")
	m := models.NewUserModel(a.appContext.DB)
	if _, err := m.GetByID(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "unable to find User")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to delete User")
	}

	if err := m.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to delete User")
	}

	return c.NoContent(http.StatusAccepted)
}
