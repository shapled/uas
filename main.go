package main

import (
	"fmt"
	"net/http"
	"strings"
	"uas/api"
	"uas/dao"
	"uas/settings"

	"github.com/labstack/echo/v4/middleware"
	"github.com/shapled/pitaya"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var yamlConfig string
var listSQLFiles bool
var catSQLFile string
var execSQLFile string

var rootCmd = &cobra.Command{
	Use:   "uas",
	Short: "Uas is a sso server",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		if yamlConfig != "" {
			if err := settings.SyncFromConfigFile(yamlConfig); err != nil {
				logrus.Fatal(err)
			}
		}

		respWrapper := func(resp pitaya.Response) interface{} {
			return map[string]interface{} {
				"error": 0,
				"message": "",
				"data": resp,
			}
		}
		errWrapper := func(err error) interface {} {
			return err
		}
		server := pitaya.NewServerWithArgs(http.StatusBadRequest, respWrapper, errWrapper)
		// app
		server.GET("/app/", api.ListApps, &api.AppListRequest{})
		server.POST("/app/", api.AddApp, &api.AddAppRequest{})
		server.PUT("/app/", api.ModifyApp, &api.ModifyAppRequest{})
		server.DELETE("/app/", api.DeleteApp, &api.DeleteAppRequest{})
		// role
		server.GET("/role/", api.ListAppRoles, &api.ListAppRolesRequest{})
		server.POST("/role/", api.AddRole, &api.AddRoleRequest{})
		server.POST("/role/permission/", api.AddRolePermission, &api.AddRolePermissionRequest{})
		server.PUT("/role/", api.ModifyRole, &api.ModifyRoleRequest{})
		server.DELETE("/role/", api.DeleteRole, &api.DeleteRoleRequest{})
		server.DELETE("/role/permission/", api.DeleteRolePermission, &api.DeleteRolePermissionRequest{})
		// permission
		server.GET("/permission/", api.ListAppPermissions, &api.ListAppPermissionsRequest{})
		server.POST("/permission/", api.AddPermission, &api.AddPermissionRequest{})
		server.PUT("/permission/", api.ModifyPermission, &api.ModifyPermissionRequest{})
		server.DELETE("/permission/", api.DeletePermission, &api.DeletePermissionRequest{})
		// allow cors
		server.Echo.Use(middleware.CORS())
		if err := server.Start(":10086"); err != nil {
			logrus.Fatal(err)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate from sql file",
	Run: func(cmd *cobra.Command, args []string) {
		if yamlConfig != "" {
			if err := settings.SyncFromConfigFile(yamlConfig); err != nil {
				logrus.Fatal(err)
			}
		}

		if listSQLFiles {
			names, err := dao.ListSQLFiles()
			if err != nil {
				logrus.Fatal(err)
			}
			fmt.Println(strings.Join(names, "\n"))
		}

		if catSQLFile != "" {
			content, err := dao.CatSQLFile(catSQLFile)
			if err != nil {
				logrus.Fatal(err)
			}
			fmt.Println(content)
		}

		if execSQLFile != "" {
			if err := dao.ExecSQLFile(execSQLFile); err != nil {
				logrus.Fatal(err)
			}
		}
	},
}

func main() {
	rootCmd.PersistentFlags().StringVarP(&yamlConfig, "config", "c", "", "config yaml file")
	migrateCmd.Flags().BoolVarP(&listSQLFiles, "list", "", false, "list available migration sql files")
	migrateCmd.Flags().StringVarP(&catSQLFile, "cat", "", "", "cat migration sql file")
	migrateCmd.Flags().StringVarP(&execSQLFile, "exec", "", "", "exec migration sql file")

	rootCmd.AddCommand(serverCmd, migrateCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
