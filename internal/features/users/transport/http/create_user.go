package users_transport_http

import (
	"net/http"

	"github.com/KirillSerge/golang-todoapp/internal/core/domain"
	core_logger "github.com/KirillSerge/golang-todoapp/internal/core/logger"
	core_http_request "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"                     example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"    example:"+79998887766"`
}

type CreateUserRespons UserDTOResponse

// CreateUser    godoc
// @Summary      Создать пользователя
// @Description  Создать нового пользователя в систему
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body CreateUserRequest true "CreateUser тело запросов"
// @Success      201 {object} CreateUserRespons "Успешно созданный пользователь"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHendler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorRespons(err, "failed to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorRespons(err, "failed to create user")

		return
	}

	response := CreateUserRespons(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUnitialized(dto.FullName, dto.PhoneNumber)
}
