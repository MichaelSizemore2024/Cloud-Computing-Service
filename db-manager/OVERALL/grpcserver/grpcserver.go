package grpcserver

import (
	"context"
	"errors"

	// remove OVERALL when moved to root
	"dbmanager/OVERALL/config"
	"dbmanager/OVERALL/generated"
	"dbmanager/OVERALL/grpc"
)

type MovieGrpcServer struct {
	// 'generated' name is based on name put in options (change name in future too lazy to recompile)
	generated.UnimplementedDBGenericServer
	service.MovieService
}

func NewGrpcServer(server *grpc.Server) *MovieGrpcServer {
	dbConfig := config.Configurations.DbConfigs
	s := &MovieGrpcServer{
		MovieService: service.NewMovieImpl(dbConfig),
	}
	moviepb.RegisterMoviePlatformServer(server.GrpcServer, s)
	return s
}

func (m *MovieGrpcServer) CreateMovie(ctx context.Context, request *moviepb.MovieRequest) (*moviepb.MovieResponse, error) {
	movieResponse, err := m.MovieService.CreateMovie(ctx, request)
	if err != nil {
		return nil, err
	}
	if movieResponse == nil {
		return &moviepb.MovieResponse{
			Status: moviepb.MovieStatus_FAILED,
		}, nil
	}
	if movieResponse.Id == constants.EmptyString {
		return nil, errors.New("err not captured but nil response received")
	}

	return &moviepb.MovieResponse{
		MovieDetails: buildMovie(movieResponse),
		Status:       moviepb.MovieStatus_CREATED,
	}, nil
}

func buildMovie(movie *models.Movie) *moviepb.MovieDetails {
	return &moviepb.MovieDetails{
		Id:          movie.Id,
		MovieName:   movie.Name,
		Genre:       movie.Genre,
		Description: movie.Desc,
		Ratings:     movie.Rating,
	}
}
