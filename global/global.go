package global

import (
	"Hwgen/config"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	H_DB     *gorm.DB
	H_DBList map[string]*gorm.DB
	H_CONFIG config.Server
	H_VIPER  *viper.Viper
	err      error
)
