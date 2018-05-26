package api

import (
	"github.com/cubixle/accounts/internal/context"
	"github.com/cubixle/accounts/internal/models"
	"github.com/labstack/echo"
)

// API holds all the route handles together and carries the app context for them all to use.
type API struct {
	appContext *context.AppContext
}

// New registers the routes and returns the router to serve.
func New(c *context.AppContext) *echo.Echo {
	e := echo.New()

	// turn off echo's setup banner.
	e.HideBanner = true
	// turn echo's logger off, we'll use our own (logrus).
	// e.Logger.SetOutput(ioutil.Discard)

	a := &API{appContext: c}

	e.GET("/teams", a.ListTeamsHandler)
	e.POST("/teams", a.CreateTeamHandler)
	e.GET("/teams/:id", a.GetTeamHandler)
	e.PUT("/teams/:id", a.UpdateTeamHandler)
	e.DELETE("/teams/:id", a.DeleteTeamHandler)

	e.GET("/users", a.ListUsersHandler)
	e.POST("/users", a.CreateUserHandler)
	e.GET("/users/:id", a.GetUserHandler)
	e.PUT("/users/:id", a.UpdateUserHandler)
	e.DELETE("/users/:id", a.DeleteUserHandler)

	return e
}

func (a *API) getAuthedUser(c echo.Context) models.User {
	return models.User{}
}
