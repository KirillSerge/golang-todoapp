package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/KirillSerge/golang-todoapp/internal/core/logger"
	core_http_request "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHendler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorRespons(err, "failed to get taskID path value")

		return
	}

	taskDoamin, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorRespons(err, "failed to get task")

		return
	}

	response := GetTaskResponse(taskDTOFromDomain(taskDoamin))

	responseHandler.JSONResponse(response, http.StatusOK)
}
