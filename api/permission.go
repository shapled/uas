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
	PermissionStatusEnabled  = 0
	PermissionStatusDisabled = 1
	PermissionStatusDeleted  = 2
)

type (
	ListAppPermissionsRequest struct {
		pitaya.BaseRequest
		AppID	int64		`json:"app_id"`
	}

	ListAppPermissionsResponse struct {
		pitaya.BaseResponse
		Permissions  []*TablePermission `json:"permissions"`
	}

	AddPermissionRequest struct {
		pitaya.BaseRequest
		Permission         string `json:"permission" validate:"required"`
		Description string `json:"description"`
	}

	AddPermissionResponse struct {
		pitaya.BaseResponse
		ID int64 `json:"id"`
	}

	ModifyPermissionRequest struct {
		pitaya.BaseRequest
		ID          int64  `json:"id" validate:"gt=0"`
		Permission         string `json:"permission" validate:"required"`
		Description string `json:"description" validate:"required"`
		Status      int    `json:"status" validate:"lt=2,gte=0"`
	}

	ModifyPermissionResponse struct {
		pitaya.BaseResponse
	}

	DeletePermissionRequest struct {
		pitaya.BaseRequest
		ID int64 `param:"id" validate:"gt=0"`
	}

	DeletePermissionResponse struct {
		pitaya.BaseResponse
	}

	TablePermission struct {
		ID        int64        `db:"id" json:"id"`
		Permission       string       `db:"permission" json:"permission"`
		Desc      string       `db:"description" json:"description"`
		Status    int          `db:"status" json:"status"`
		CreatedBy int64        `db:"created_by" json:"created_by"`
		CreatedAt time.Time    `db:"created_at" json:"created_at"`
		UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
		DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
	}

	PermissionWithRoleID struct {
		TablePermission
		RoleID		int64		`db:"role_id" json:"role_id"`
	}
)

func ListAppPermissions(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*ListAppPermissionsRequest)
	permissions := make([]*TablePermission, 0)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		return db.SelectContext(ctx, &permissions,
			`select id,permission,description,status,created_by,created_at,updated_at,deleted_at 
					from uas_permission
					where status != ? and id in (
						select permission_id
						from uas_role_permission
						where role_id in (
							select id
							from role
							where status != ? and app_id = ?
						)
					)
					limit ?, ?`,
			PermissionStatusDeleted, RoleStatusDeleted, request.AppID)
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return &ListAppPermissionsResponse{Permissions: permissions}, nil
}

func AddPermission(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*AddPermissionRequest)
	var id int64
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "insert into uas_permission(permission,description,created_by) values(?, ?, ?)",
			request.Permission, request.Description, 0)
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
	return &AddPermissionResponse{ID: id}, nil
}

func ModifyPermission(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*ModifyPermissionRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "update uas_permission set permission=?, description=?, status=? where id = ?",
			request.Permission, request.Description, request.Status, request.ID)
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
	return &ModifyPermissionResponse{}, nil
}

func DeletePermission(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*DeletePermissionRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "update uas_permission set status = ?, deleted_at = ? where id = ?",
			PermissionStatusDeleted, time.Now(), request.ID)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("permission not found: %d", request.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &DeletePermissionResponse{}, nil
}
