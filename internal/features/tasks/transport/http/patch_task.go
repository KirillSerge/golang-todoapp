package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/KirillSerge/golang-todoapp/internal/core/domain"
	core_logger "github.com/KirillSerge/golang-todoapp/internal/core/logger"
	core_http_request "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/types"
)

type PathTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

// PatchTasks      godoc
// @Summary      Обновить задачу
// @Description  Обновляет информацию об уже существующей задаче
// @Description  Изменение информации об уже существующем в системе пользователи
// @Description  ### Логика обновления полей (Three-state logic):
// @Description  1. **Поле не переданно**:`description` игнорируется, значение в БД не меняется
// @Description  2. **Явно передано значение**: `"description":"выйти на прогулку с бобиком"` - устанавливает новый номер телефона в БД
// @Description  3. **Передан null**: `"description":null` - очищает поле в БД (set to NULL)
// @Description  Ограничение: `title`и `completed` не может быть выставлен как null
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id path int true "ID изменяемой задачи"
// @Param        request body PathTaskRequest true "PatchTask тело запроса"
// @Success      200 {object} GetTasksResponse "Усрешно измененная задача"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      404 {object} core_http_response.ErrorResponse "Task not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Conflict"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /tasks/{id} [patch]
func (r *PathTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}

		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("'Title' must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 1000 {
				return fmt.Errorf("'Description' must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}

	return nil
}

type PatchUserResponse TaskDTOResponse

func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHendler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorRespons(err, "failed to get taskID path value")

		return
	}

	var request PathTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorRespons(err, "failed to decode and validate HTTP request")

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorRespons(err, "failed to patch task")

		return
	}

	response := PatchUserResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PathTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
