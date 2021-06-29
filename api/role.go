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
	RoleStatusEnabled  = 0
	RoleStatusDisabled = 1
	RoleStatusDeleted  = 2
)

type (
	ListAppRolesRequest struct {
		pitaya.BaseRequest
		AppID	int64	`json:"app_id"`
	}

	ListAppRolesResponse struct {
		pitaya.BaseResponse
		Roles	[]*RoleWithPermissions	`json:"roles"`
	}

	AddRoleRequest struct {
		pitaya.BaseRequest
		Role         string `json:"role" validate:"required"`
		AppID		int64	`json:"app_id"`
		Description string `json:"description"`
	}

	AddRoleResponse struct {
		pitaya.BaseResponse
		ID int64 `json:"id"`
	}

	AddRolePermissionRequest struct {
		pitaya.BaseRequest
		ID int64 `json:"id"`
		PermissionID int64	`json:"permission_id"`
	}

	AddRolePermissionResponse struct {
		pitaya.BaseResponse
	}

	ModifyRoleRequest struct {
		pitaya.BaseRequest
		ID          int64  `json:"id" validate:"gt=0"`
		Role         string `json:"role" validate:"required"`
		Description string `json:"description" validate:"required"`
		Status      int    `json:"status" validate:"lt=2,gte=0"`
	}

	ModifyRoleResponse struct {
		pitaya.BaseResponse
	}

	DeleteRoleRequest struct {
		pitaya.BaseRequest
		ID int64 `param:"id" validate:"gt=0"`
	}

	DeleteRoleResponse struct {
		pitaya.BaseResponse
	}

	DeleteRolePermissionRequest struct {
		pitaya.BaseRequest
		ID int64 `param:"id" validate:"gt=0"`
		PermissionID int64 `param:"permission_id" validate:"gt=0"`
	}

	DeleteRolePermissionResponse struct {
		pitaya.BaseResponse
	}

	TableRole struct {
		ID        int64        `db:"id" json:"id"`
		AppID	  int64		   `db:"app_id" json:"app_id"`
		Role      string       `db:"role" json:"role"`
		Desc      string       `db:"description" json:"description"`
		Status    int          `db:"status" json:"status"`
		CreatedBy int64        `db:"created_by" json:"created_by"`
		CreatedAt time.Time    `db:"created_at" json:"created_at"`
		UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
		DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
	}

	RoleWithPermissions struct {
		TableRole
		Permissions []*TablePermission	`db:"permissions" json:"permissions"`
	}
)

func ListAppRoles(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*ListAppRolesRequest)
	roles := make([]*RoleWithPermissions, 0)
	permissions := make([]*PermissionWithRoleID, 0)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		err := db.SelectContext(ctx, &roles,
			`select id,` + "`role`" + `,description,status,created_by,created_at,updated_at,deleted_at 
					from uas_role
					where status != ? and app_id = ?
					order by id`,
			RoleStatusDeleted, request.AppID)
		if err != nil {
			return err
		}
		return db.SelectContext(ctx, &permissions,
			`select p.id as id, 
       				rp.role_id as role_id, 
       				p.permission as permission,
       				p.description as description, 
       				p.status as status, 
       				p.created_by as created_by,
       				p.created_at as created_at, 
       				p.updated_at as updated_at, 
       				p.deleted_at as p.deleted_at 
				from uas_role r 
				    left join uas_role_permission rp on r.app_id = ? and r.status != ? and r.id = rp.role_id
				    left join uas_permission p on p.status != ? and rp.permissions_id = p.id
				order by role_id, id`,
			request.AppID, RoleStatusDeleted, PermissionStatusDeleted)
	})
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	// i 是 roles 的下标， j 是 permissions 的下标，因为两个数组 id 都有序，可以直接匹配
	for i, j := 0, 0; i<len(roles); i++ {
		roles[i].Permissions = make([]*TablePermission, 0)
		for ; j < len(permissions) && permissions[j].RoleID == roles[i].ID; j++ {
			roles[i].Permissions = append(roles[i].Permissions, &permissions[j].TablePermission)
		}
	}
	return &ListAppRolesResponse{Roles: roles}, nil
}

func AddRole(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*AddRoleRequest)
	var id int64
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "insert into uas_role(app_id,`role`,description,created_by) values(?, ?, ?, ?)",
			request.AppID, request.Role, request.Description, 0)
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
	return &AddRoleResponse{ID: id}, nil
}

func AddRolePermission(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*AddRolePermissionRequest)
	var id int64
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx,
			"insert into uas_role_permission(role_id, permission_id) values(?, ?)",
			request.ID, request.PermissionID)
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
	return &AddRolePermissionResponse{}, nil
}

func ModifyRole(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*ModifyRoleRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "update uas_role set `role`=?, description=?, status=? where id = ?",
			request.Role, request.Description, request.Status, request.ID)
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
	return &ModifyRoleResponse{}, nil
}

func DeleteRole(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*DeleteRoleRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "update uas_role set status = ? where id = ?",
			RoleStatusDeleted, request.ID)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("Role not found: %d", request.ID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &DeleteRoleResponse{}, nil
}

func DeleteRolePermission(req pitaya.Request) (pitaya.Response, error) {
	request := req.(*DeleteRolePermissionRequest)
	err := dao.Dao(func(ctx context.Context, db *sqlx.DB) error {
		res, err := db.ExecContext(ctx, "delete from uas_role_permission where role_id = ? and permission_id = ?",
			request.ID, request.PermissionID)
		if err != nil {
			return err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if n != 1 {
			return fmt.Errorf("role %d's permission %d not found", request.ID, request.PermissionID)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &DeleteRolePermissionResponse{}, nil
}
