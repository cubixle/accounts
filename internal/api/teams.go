package api

import (
	"net/http"

	"github.com/cubixle/accounts/internal/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"gopkg.in/go-playground/validator.v9"
)

// ListTeamsHandler will list all teams related to the authed user.
func (a *API) ListTeamsHandler(c echo.Context) error {
	teamModel := models.NewTeamModel(a.appContext.DB)

	authedUser := a.getAuthedUser(c)

	teams, err := teamModel.GetAll(authedUser.TeamID)
	if err != nil {
		// log
	}

	return c.JSON(http.StatusOK, teams)
}

// GetTeamHandler will get a single team for the given id.
func (a *API) GetTeamHandler(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid team id")
	}

	teamModel := models.NewTeamModel(a.appContext.DB)
	team, err := teamModel.GetByID(id)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.ErrNotFound
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to access resource")
	}

	return c.JSON(http.StatusOK, team)
}

// CreateTeamHandler will create a new team within the database.
func (a *API) CreateTeamHandler(c echo.Context) error {
	var t models.Team
	if err := c.Bind(&t); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	m := models.NewTeamModel(a.appContext.DB)
	team, err := m.Create(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to create team")
	}

	return c.JSON(http.StatusCreated, team)
}

func (a *API) UpdateTeamHandler(c echo.Context) error {
	var t models.Team
	if err := c.Bind(&t); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(&t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	t.ID = c.Param("id")
	m := models.NewTeamModel(a.appContext.DB)

	// check entity exists
	if _, err := m.GetByID(t.ID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "unable to find team")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to update team")
	}

	team, err := m.Update(t)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to update team")
	}

	return c.JSON(http.StatusAccepted, team)
}

// DeleteTeam will delete the team from the database.
func (a *API) DeleteTeamHandler(c echo.Context) error {
	id := c.Param("id")
	m := models.NewTeamModel(a.appContext.DB)
	if _, err := m.GetByID(id); err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "unable to find team")
		}

		return echo.NewHTTPError(http.StatusInternalServerError, "unable to delete team")
	}

	if err := m.Delete(id); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "unable to delete team")
	}

	return c.NoContent(http.StatusAccepted)
}
