package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/croisade/chimichanga/pkg/models"
	"github.com/croisade/chimichanga/pkg/services"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCreateRun(t *testing.T) {
	run := &models.Run{Speed: 6.0, Lap: 0, Distance: 3.0, Time: "30:00", Incline: 0.0, AccountId: "123"}
	response := &models.Run{}

	jsonValue, _ := json.Marshal(run)
	req, _ := http.NewRequest("POST", "/run/create", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	fmt.Println(w)
	json.Unmarshal(w.Body.Bytes(), response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, run.Speed, response.Speed)
}

func TestGetRun(t *testing.T) {
	var response *models.Run
	run := &models.Run{Speed: 6.0, Lap: 0, Distance: 3.0, Time: "30:00", Incline: 0.0, AccountId: "234"}
	createdRun, _ := runService.CreateRun(run)
	request := &services.RunRequest{AccountId: createdRun.AccountId, RunId: createdRun.RunId}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("GET", "/run", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, run.AccountId, response.AccountId)

}

func TestUpdateRun(t *testing.T) {
	run := &models.Run{Speed: 6.0, Lap: 0, Distance: 3.0, Time: "30:00", Incline: 0.0, AccountId: "567"}
	createdRun, _ := runService.CreateRun(run)
	request := &services.RunUpdateRequest{AccountId: createdRun.AccountId, RunId: createdRun.RunId, Lap: 1}
	var response *models.Run

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("PUT", "/run/update", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, 1, response.Lap)
}
func TestDeleteRun(t *testing.T) {
	run := &models.Run{Speed: 6.0, Lap: 0, Distance: 3.0, Time: "30:00", Incline: 0.0, AccountId: "890"}
	createdRun, _ := runService.CreateRun(run)
	request := &services.RunRequest{AccountId: createdRun.AccountId, RunId: createdRun.RunId}

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("DELETE", "/run/delete", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	_, err := runService.GetRun(request)
	assert.NotNil(t, err)
}

func TestFetchRun(t *testing.T) {
	runsCollection.DeleteMany(ctx, bson.D{{}})
	run := &models.Run{Speed: 6.0, Lap: 0, Distance: 3.0, Time: "30:00", Incline: 0.0, AccountId: "123"}
	createdRun, _ := runService.CreateRun(run)

	request := &services.RunFetchRequest{AccountId: createdRun.AccountId}
	var response []*models.Run

	jsonValue, _ := json.Marshal(request)
	req, _ := http.NewRequest("GET", "/run/fetch", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Println(response)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, response[0].AccountId, run.AccountId)
}
