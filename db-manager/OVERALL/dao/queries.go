package dao

import (
	"context"
	"dbmanager/OVERALL/log"

	"go.uber.org/zap"

	"dbmanager/OVERALL/db"

	daomodels "github.com/go-grpc-service/internal/dao/dao_models"
	"github.com/scylladb/gocqlx/v2/table"
)

const (
	id = "id"
)

// No changes
func getTable(tableName string, cols, partitionKeys, sortKeys []string) db.Table {
	tableMeta := table.Metadata{
		Name:    tableName,
		Columns: cols,
		PartKey: partitionKeys,
		SortKey: sortKeys,
	}
	return table.New(tableMeta)
}

// Just doing insert for now
func insertInDB(ctx context.Context, session db.SessionWrapperService, movies *daomodels.Movies, tableName string, cols, partitionKeys, sortKeys []string) error {
	// not sure how we are going to get these lol
	configTable := getTable(tableName, cols, partitionKeys, sortKeys)
	query := session.Query(configTable.Insert()).BindStruct(movies)
	if err := query.Exec(ctx); err != nil {
		log.Logger.Error("error inserting movie name in db ", zap.String("key", movies.Name))
		return err
	}
	log.Logger.Info("Movie insertion successful", zap.String(id, movies.MovieID))
	return nil
}
