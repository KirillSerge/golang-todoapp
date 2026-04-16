package users_transport_http

import (
	"net/http"

	core_logger "github.com/KirillSerge/golang-todoapp/internal/core/logger"
	core_http_response "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/response"
	core_http_utils "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/utils"
)

type GetUserResponse UserDTOResponse

func (h *UsersHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHendler := core_http_response.NewHTTPResponseHendler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHendler.ErrorRespons(err, "failed to get userID path value")
		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHendler.ErrorRespons(err, "failed to get user")
		return
	}

	response := GetUserResponse(userDTOFromDomain(user))

	responseHendler.JSONResponse(response, http.StatusOK)
}
