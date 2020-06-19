package config

import (
	"github.com/spf13/viper"
)

const (
	//ExplorerEndpointKey defines the explorer to fetch data
	ExplorerEndpointKey = "EXPLORER_ENDPOINT"
	//TaxiEndpointKey defines the Taxi backend
	TaxiEndpointKey = "TAXI_ENDPOINT"
	// ConfigPathKey defines the config path
	ConfigPathKey = "CONFIG_PATH"
)

var vip *viper.Viper

func init() {
	vip = viper.New()
	vip.SetEnvPrefix("TAXI")
	vip.AutomaticEnv()

	vip.SetDefault(ExplorerEndpointKey, "https://blockstream.info/liquid/api")
	vip.SetDefault(TaxiEndpointKey, "https://3moyhezvi3.execute-api.eu-west-1.amazonaws.com/production")
	vip.SetDefault(ConfigPathKey, "taxi-cli.config")

	vip.SetConfigName(vip.GetString(ConfigPathKey))
	vip.AddConfigPath(".")
	vip.ReadInConfig()
}

// GetString ...
func GetString(key string) string {
	return vip.GetString(key)
}

// GetInt ...
func GetInt(key string) int {
	return vip.GetInt(key)
}

// GetBool ...
func GetBool(key string) bool {
	return vip.GetBool(key)
}
