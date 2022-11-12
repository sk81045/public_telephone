package global

import (
	"Hwgen/config"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	H_DB     *gorm.DB
	H_DBList map[string]*gorm.DB
	H_CONFIG config.Server
	H_VIPER  *viper.Viper
	H_LOG    *zap.Logger
	err      error
)
