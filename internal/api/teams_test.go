package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cubixle/accounts/internal/api"
	"github.com/cubixle/accounts/internal/models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetTeams(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)

	req := httptest.NewRequest(http.MethodGet, "/teams", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetTeam(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)

	teamRec := createTeam(e)
	if !assert.Equal(t, http.StatusCreated, teamRec.Code) {
		t.Fatal(teamRec.Body.String())
	}

	team := new(models.Team)
	err := json.Unmarshal(teamRec.Body.Bytes(), team)
	if err != nil {
		t.Fatal("WTF mate JSON failed")
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/teams/"+team.ID, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Error(rec.Body.String())
		return
	}

	foundTeam := new(models.Team)
	json.Unmarshal(rec.Body.Bytes(), foundTeam)

	assert.NotEqual(t, foundTeam.Name, "")
}

func TestCreateTeam(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)
	rec := createTeam(e)
	if !assert.Equal(t, http.StatusCreated, rec.Code) {
		t.Fatal(rec.Body.String())
	}
}

func TestUpdateTeam(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)
	createRec := createTeam(e)
	if !assert.Equal(t, http.StatusCreated, createRec.Code) {
		t.Fatal(createRec.Body.String())
		return
	}

	team := new(models.Team)
	json.Unmarshal(createRec.Body.Bytes(), team)

	team.Name = "testing2"
	teamBytes, err := json.Marshal(team)
	if err != nil {
		t.Fatal("")
		return
	}

	req := httptest.NewRequest(http.MethodPut, "/teams/"+team.ID, bytes.NewBuffer(teamBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if !assert.Equal(t, http.StatusAccepted, rec.Code) {
		t.Fatal(rec.Body.String())
		return
	}

	team = new(models.Team)
	json.Unmarshal(rec.Body.Bytes(), team)

	assert.Equal(t, team.Name, "testing2")
}

func createTeam(e *echo.Echo) *httptest.ResponseRecorder {
	teamJSON := `{"name": "testing"}`
	req := httptest.NewRequest(http.MethodPost, "/teams", strings.NewReader(teamJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}
