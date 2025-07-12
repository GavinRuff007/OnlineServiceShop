package repository

import (
	"RestGoTest/src/config"
	"RestGoTest/src/constant"
	"RestGoTest/src/database"
	"RestGoTest/src/filter"
	"RestGoTest/src/helper/service_errors"
	"RestGoTest/src/pkg/logging"
	"RestGoTest/src/util"
	"context"
	"database/sql"
	"time"

	"gorm.io/gorm"
)

const softDeleteExp string = "id = ? and deleted_by is null"

type BaseRepository[TEntity any] struct {
	database *gorm.DB
	logger   logging.Logger
	preloads []database.PreloadEntity
}

func NewBaseRepository[TEntity any](cfg *config.Config, preloads []database.PreloadEntity) *BaseRepository[TEntity] {
	return &BaseRepository[TEntity]{
		database: database.GetDb(),
		logger:   logging.NewLogger(cfg),
		preloads: preloads,
	}
}

func (r BaseRepository[TEntity]) Create(ctx context.Context, entity TEntity) (TEntity, error) {
	tx := r.database.WithContext(ctx).Begin()
	err := tx.
		Create(&entity).
		Error
	if err != nil {
		tx.Rollback()
		r.logger.Error(logging.Mysql, logging.Insert, err.Error(), nil)
		return entity, err
	}
	tx.Commit()

	return entity, nil
}

func (r BaseRepository[TEntity]) Update(ctx context.Context, id int, entity map[string]interface{}) (TEntity, error) {
	snakeMap := map[string]interface{}{}
	for k, v := range entity {
		snakeMap[util.ToSnakeCase(k)] = v
	}
	snakeMap["modified_by"] = &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true}
	snakeMap["modified_at"] = sql.NullTime{Valid: true, Time: time.Now().UTC()}
	model := new(TEntity)
	tx := r.database.WithContext(ctx).Begin()
	if err := tx.Model(model).
		Where(softDeleteExp, id).
		Updates(snakeMap).
		Error; err != nil {
		tx.Rollback()
		r.logger.Error(logging.Mysql, logging.Update, err.Error(), nil)
		return *model, err
	}
	tx.Commit()
	return *model, nil
}

func (r BaseRepository[TEntity]) Delete(ctx context.Context, id int) error {
	tx := r.database.WithContext(ctx).Begin()

	model := new(TEntity)

	deleteMap := map[string]interface{}{
		"deleted_by": &sql.NullInt64{Int64: int64(ctx.Value(constant.UserIdKey).(float64)), Valid: true},
		"deleted_at": sql.NullTime{Valid: true, Time: time.Now().UTC()},
	}

	if ctx.Value(constant.UserIdKey) == nil {
		return &service_errors.ServiceError{EndUserMessage: constant.PermissionDenied}
	}
	if cnt := tx.
		Model(model).
		Where(softDeleteExp, id).
		Updates(deleteMap).
		RowsAffected; cnt == 0 {
		tx.Rollback()
		r.logger.Error(logging.Mysql, logging.Update, constant.RecordNotFound, nil)
		return &service_errors.ServiceError{EndUserMessage: constant.RecordNotFound}
	}
	tx.Commit()
	return nil
}

func (r BaseRepository[TEntity]) GetById(ctx context.Context, id int) (TEntity, error) {
	model := new(TEntity)
	db := database.Preload(r.database, r.preloads)
	err := db.
		Where(softDeleteExp, id).
		First(model).
		Error
	if err != nil {
		return *model, err
	}
	return *model, nil
}

func (r BaseRepository[TEntity]) GetByFilter(ctx context.Context, req filter.PaginationInputWithFilter) (int64, *[]TEntity, error) {
	model := new(TEntity)
	var items *[]TEntity

	db := database.Preload(r.database, r.preloads)
	query := database.GenerateDynamicQuery[TEntity](&req.DynamicFilter)
	sort := database.GenerateDynamicSort[TEntity](&req.DynamicFilter)
	var totalRows int64 = 0

	db.
		Model(model).
		Where(query).
		Count(&totalRows)

	err := db.
		Where(query).
		Offset(req.GetOffset()).
		Limit(req.GetPageSize()).
		Order(sort).
		Find(&items).
		Error

	if err != nil {
		return 0, &[]TEntity{}, err
	}
	return totalRows, items, err

}
