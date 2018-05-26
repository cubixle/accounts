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
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

func TestGetUsers(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestGetUser(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)

	userRec := createUser(e)
	if !assert.Equal(t, http.StatusCreated, userRec.Code) {
		t.Fatal(userRec.Body.String())
	}

	User := new(models.User)
	err := json.Unmarshal(userRec.Body.Bytes(), User)
	if err != nil {
		t.Fatal("WTF mate JSON failed")
		return
	}

	req := httptest.NewRequest(http.MethodGet, "/users/"+User.ID, nil)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	if !assert.Equal(t, http.StatusOK, rec.Code) {
		t.Error(rec.Body.String())
		return
	}

	foundUser := new(models.User)
	json.Unmarshal(rec.Body.Bytes(), foundUser)

	assert.NotEqual(t, foundUser.Name, "")
}

func TestCreateUser(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)
	rec := createUser(e)
	if !assert.Equal(t, http.StatusCreated, rec.Code) {
		t.Fatal(rec.Body.String())
	}
}

func TestUpdateUser(t *testing.T) {
	ctx, tearDownDB := appContext()
	defer tearDownDB()

	e := api.New(ctx)
	createRec := createUser(e)
	if !assert.Equal(t, http.StatusCreated, createRec.Code) {
		t.Fatal(createRec.Body.String())
		return
	}

	User := new(models.User)
	json.Unmarshal(createRec.Body.Bytes(), User)

	User.Name = "testing2"
	userBytes, err := json.Marshal(User)
	if err != nil {
		t.Fatal("")
		return
	}

	req := httptest.NewRequest(http.MethodPut, "/users/"+User.ID, bytes.NewBuffer(userBytes))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if !assert.Equal(t, http.StatusAccepted, rec.Code) {
		t.Fatal(rec.Body.String())
		return
	}

	User = new(models.User)
	json.Unmarshal(rec.Body.Bytes(), User)

	assert.Equal(t, User.Name, "testing2")
}

func createUser(e *echo.Echo) *httptest.ResponseRecorder {
	userJSON := `{"name": "testing"}`
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	return rec
}
