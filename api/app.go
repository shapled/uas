package api

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"uas/dao"

	"github.com/jmoiron/sqlx"
	"github.com/shapled/pitaya"
	"github.com/sirupsen/logrus"
)

const (
	AppStatusEnabled  = 0
	AppStatusDisabled = 1
	AppStatusDeleted  = 2
)

type (
	AppListRequest struct {
		pitaya.BaseRequest
		Pagination
	}

	AppListResponse struct {
		pitaya.BaseResponse
		Total int64       `json:"total"`
		Page  int         `json:"page"`
		Size  int         `json:"size"`
		Data  []*TableApp `json:"data"`
	}

	AddAppRequest struct {
		pitaya.BaseRequest
		App         string `json:"app" validate:"required"`
		Description string `json:"description"`
	}

	AddAppResponse struct {
		pitaya.BaseResponse
		ID int64 `json:"id"`
	}

	ModifyAppRequest struct {
		pitaya.BaseRequest
		ID          int64  `json:"id" validate:"gt=0"`
		App         string `json:"app" validate:"required"`
		Description string `json:"description" validate:"required"`
		Status      int    `json:"status" validate:"lt=2,gte=0"`
	}

	ModifyAppResponse struct {
		pitaya.BaseResponse
	}

	DeleteAppRequest struct {
		pitaya.BaseRequest
		ID int64 `param:"id" validate:"gt=0"`
	}

	DeleteAppResponse struct {
		pitaya.BaseResponse
	}

	TableApp struct {
		ID        int64        `db:"id" json:"id"`
		App       string       `db:"app" json:"app"`
		Desc      string       `db:"description" json:"description"`
		Status    int          `db:"status" json:"status"`
		CreatedBy int64        `db:"created_by" json:"created_by"`
		CreatedAt time.Time    `db:"created_at" json:"created_at"`
		UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
		DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
	}
)

func ListApps(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*AppListRequest)
	request.FormatPageAndSize(5, 20, 20)
	offset, limit := request.CalcOffset()
	var count int64
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		return db.GetContext(ctx, &count, "select count(*) from uas_app where status != ?", AppStatusDeleted)
	})
	if err != nil {
		return nil, err
	}
	apps := make([]*TableApp, 0)
	err = dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		return db.SelectContext(ctx, &apps,
			`select id,app,description,status,created_by,created_at,updated_at,deleted_at 
					from uas_app
					where status != ?
					limit ?, ?`,
			AppStatusDeleted, offset, limit)
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &AppListResponse{Total: count, Page: request.Page, Size: request.Size, Data: apps}, nil
}

func AddApp(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*AddAppRequest)
	var id int64
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "insert into uas_app(app,description,created_by) values(?, ?, ?)",
			request.App, request.Description, 0)
		if err != nil {
			return err
		}
		id, err = res.LastInsertId()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &AddAppResponse{ID: id}, nil
}

func ModifyApp(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*ModifyAppRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "update uas_app set app=?, description=?, status=? where id = ?",
			request.App, request.Description, request.Status, request.ID)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("no such id: %d", request.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ModifyAppResponse{}, nil
}

func DeleteApp(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*DeleteAppRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "update uas_app set status = ? where id = ?",
			AppStatusDeleted, request.ID)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("app not found: %d", request.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &ModifyAppResponse{}, nil
}
