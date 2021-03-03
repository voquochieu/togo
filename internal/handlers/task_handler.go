package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/manabie-com/togo/internal/models"
	"github.com/manabie-com/togo/internal/services"
	"github.com/manabie-com/togo/internal/utils"
)

type TaskHandler struct {
	TaskService services.TaskService
}

func (h *TaskHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	log.Println(req.Method, req.URL.Path)
	switch req.Method {
	case http.MethodGet:
		userID, ok := utils.UserIDFromCtx(req.Context())
		if !ok {
			utils.RespondWithError(resp, http.StatusUnauthorized, "Invalid user ID")
			return
		}
		from := req.FormValue("from")
		start, err := time.Parse(time.RFC3339, from)
		if err != nil {
			log.Printf("User %s get tasks with invalid start time %v", userID, from)
			utils.RespondWithError(resp, http.StatusBadRequest, "Invalid start time")
			return
		}

		to := req.FormValue("to")
		end, err := time.Parse(time.RFC3339, to)
		if err != nil {
			log.Printf("User %s get tasks with invalid end time %v", userID, from)
			utils.RespondWithError(resp, http.StatusBadRequest, "Invalid end time")
			return
		}
		tasks, err := h.TaskService.ListTasks(req.Context(), userID, start, end)
		if err != nil {
			utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
			return
		}
		data := struct {
			Data *[]models.Task `json:"data"`
		}{
			tasks,
		}
		utils.RespondWithJSON(resp, http.StatusOK, data)
	case http.MethodPost:
		task := &models.Task{}
		err := json.NewDecoder(req.Body).Decode(task)
		defer req.Body.Close()
		if err != nil {
			utils.RespondWithError(resp, http.StatusBadRequest, "Invalid request")
			return
		}
		userID, ok := utils.UserIDFromCtx(req.Context())
		if !ok {
			utils.RespondWithError(resp, http.StatusUnauthorized, "Invalid user ID")
			return
		}

		task, err = h.TaskService.AddTask(req.Context(), task, userID)
		if err != nil {
			if strings.Contains(err.Error(), "You are limited") {
				utils.RespondWithError(resp, http.StatusTooManyRequests, err.Error())
				return
			}
			utils.RespondWithError(resp, http.StatusInternalServerError, err.Error())
			return
		}
		data := struct {
			Data *models.Task `json:"data"`
		}{
			task,
		}
		utils.RespondWithJSON(resp, http.StatusOK, data)
	}
	return
}
