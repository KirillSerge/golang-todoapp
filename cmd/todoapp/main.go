package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/KirillSerge/golang-todoapp/internal/core/config"
	core_logger "github.com/KirillSerge/golang-todoapp/internal/core/logger"
	core_pgx_pool "github.com/KirillSerge/golang-todoapp/internal/core/repository/postgres/pool/pgx"
	core_http_middleware "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/middlrware"
	core_http_server "github.com/KirillSerge/golang-todoapp/internal/core/transport/http/server"
	statistics_postgres_repository "github.com/KirillSerge/golang-todoapp/internal/features/staristics/repository/postgres"
	statistics_service "github.com/KirillSerge/golang-todoapp/internal/features/staristics/service"
	statistics_transport_http "github.com/KirillSerge/golang-todoapp/internal/features/staristics/transport/http"
	tasks_postgres_repository "github.com/KirillSerge/golang-todoapp/internal/features/tasks/repository/postgres"
	tasks_service "github.com/KirillSerge/golang-todoapp/internal/features/tasks/service"
	tasks_transport_http "github.com/KirillSerge/golang-todoapp/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/KirillSerge/golang-todoapp/internal/features/users/repository/postgres"
	user_service "github.com/KirillSerge/golang-todoapp/internal/features/users/service"
	users_transport_http "github.com/KirillSerge/golang-todoapp/internal/features/users/transport/http"
	web_fs_repository "github.com/KirillSerge/golang-todoapp/internal/features/web/repository/file_system"
	web_service "github.com/KirillSerge/golang-todoapp/internal/features/web/service"
	web_transport_http "github.com/KirillSerge/golang-todoapp/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/KirillSerge/golang-todoapp/docs"
)

// @title       Golang Todo API
// @version     1.0
// @description Todo Application REST-API scheme
// @host        127.0.0.1:5050
// @BasePath    /api/v1
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

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

	logger.Debug("application time zone", zap.Any("zone", time.Local))

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

	logger.Debug("initialzing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initialzing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatistics(pool)
	statisticsService := statistics_service.NewStatisticsService(statisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initialzing feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initiazling HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)
	apiVersionRouterV1 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouterV1.RegisterRouter(usersTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRouter(tasksTransportHTTP.Routes()...)
	apiVersionRouterV1.RegisterRouter(statisticsTransportHTTP.Routes()...)

	/*apiVersionRouterV2 := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion2, core_http_middleware.Dummy("api v2 middleware"))
	apiVersionRouterV2.RegisterRouter(usersTransportHTTP.Routes()...)*/

	httpServer.RegisterAPIRouters(apiVersionRouterV1)
	//apiVersionRouterV2)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)
	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
