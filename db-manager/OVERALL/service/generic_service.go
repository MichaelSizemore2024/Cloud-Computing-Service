package service

import (
	"context"

	"dbmanager/OVERALL/config"
	"dbmanager/OVERALL/dao"

	"github.com/go-grpc-service/internal/models"
	"github.com/go-grpc-service/resources/moviepb"
)

type MovieService interface {
	CreateMovie(ctx context.Context, request *moviepb.MovieRequest) (*models.Movie, error)
}

type MovieServiceImpl struct {
	movieDao dao.MovieDao
}

func NewMovieImpl(dbConfigs config.DbConfigs) MovieService {
	return &MovieServiceImpl{
		movieDao: dao.NewMovieDaoImpl(dbConfigs),
	}
}

func (m *MovieServiceImpl) CreateMovie(ctx context.Context, request *moviepb.MovieRequest) (*models.Movie, error) {
}
