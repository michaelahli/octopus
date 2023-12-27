package mainpkg

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

func ServiceConfig(filename string) (*viper.Viper, error) {
	if filename == "" {
		return nil, fmt.Errorf("filename required")
	}
	cfg := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	cfg.SetConfigFile(filename)
	cfg.SetConfigType("ini")
	cfg.AutomaticEnv()
	return cfg, cfg.ReadInConfig()
}
