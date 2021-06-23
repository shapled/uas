package main

import (
	"fmt"
	"github.com/shapled/pitaya"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"uas/api"
	"uas/dao"
	"uas/settings"
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
	Use: "server",
	Short: "Start server",
	Run: func(cmd *cobra.Command, args []string) {
		if yamlConfig != "" {
			if err := settings.SyncFromConfigFile(yamlConfig); err != nil {
				logrus.Fatal(err)
			}
		}

		server := pitaya.NewServer()
		server.GET("/app/", api.ListApps, &api.AppListRequest{})
		if err := server.Start(":10086"); err != nil {
			logrus.Fatal(err)
		}
	},
}

var migrateCmd = &cobra.Command{
	Use: "migrate",
	Short: "Migrate from sql file",
	Run: func(cmd *cobra.Command, args []string) {
		if yamlConfig != "" {
			if err := settings.SyncFromConfigFile(yamlConfig); err != nil {
				logrus.Fatal(err)
			}
		}

		if listSQLFiles {
			names, err := dao.ListSQLFiles();
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
