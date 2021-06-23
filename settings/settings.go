package settings

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type dbSettings struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type uasSettings struct {
	DB *dbSettings `yaml:"db"`
}

var UASSettings = &uasSettings{
	DB: &dbSettings{
		Host:     "localhost",
		Port:     3306,
		Name:     "uas",
		User:     "uas",
		Password: "123456",
	},
}

func SyncFromConfigFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		logrus.Error(err)
		return fmt.Errorf("config not exists: %s", filepath)
	}
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		logrus.Error(err)
		return fmt.Errorf("read file failed: %s", filepath)
	}
	if err := yaml.Unmarshal(content, UASSettings); err != nil {
		logrus.Error(err)
		return fmt.Errorf("read file failed: %s", filepath)
	}
	return nil
}
