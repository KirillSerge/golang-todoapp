package statistics_postgres_repository

import core_postgres_pool "github.com/KirillSerge/golang-todoapp/internal/core/repository/postgres/pool"

type StatisticsRepository struct {
	pool core_postgres_pool.Pool
}

func NewStatistics(pool core_postgres_pool.Pool) *StatisticsRepository {
	return &StatisticsRepository{
		pool: pool,
	}
}
