package api

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/shapled/pitaya"
	"github.com/sirupsen/logrus"
	"time"
	"uas/dao"
)

type (
	AppListRequest struct {
		pitaya.BaseRequest
	}

	AppListResponse struct {
		pitaya.BaseResponse
		Data []*TableApp
	}

	TableApp struct {
		ID		int64	`db:"id" json:"id"`
		App		string	`db:"app" json:"app"`
		Desc	string	`db:"desc" json:"desc"`
		CreatedBy	int64	`db:"created_by" json:"created_by"`
		CreatedAt	time.Time	`db:"created_at" json:"created_at"`
		UpdatedAt	time.Time	`db:"updated_at" json:"updated_at"`
		DeletedAt	time.Time	`db:"deleted_at" json:"deleted_at"`
	}
)

func ListApps(req pitaya.Request) (pitaya.Response, error) {
	apps := make([]*TableApp, 0, 0)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		return db.SelectContext(ctx, &apps,
			`select id,app,description,created_by,created_at,updated_at,deleted_at from uas_app limit ?, ?`,
			0, 10)
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &AppListResponse{Data: apps}, nil
}
