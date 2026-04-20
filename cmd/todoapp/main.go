package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	core_logger "github.com/KirillSerge/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/KirillSerge/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/middlrware"
	core_http_server "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/server"
	users_postgres_repository "github.com/KirillSerge/golang-todoapp/internal/features/users/repository/postgres"
	user_service "github.com/KirillSerge/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/KirillSerge/golang-todoapp/internal/features/users/transport/http"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initiazling postgres connection pool")

	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initialzing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := user_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initiazling HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRouter(usersTransportHTTP.Routes()...)

	/*apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2, core_http_middleware.Dummy("api v2 middleware"))
	apiVersionRouterV2.RegisterRouter(usersTransportHTTP.Routes()...)*/

	httpServer.RegisterAPIRouters(apiVersionRouterV1)
	//apiVersionRouterV2)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
