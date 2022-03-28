package config

import (
	"github.com/spf13/viper"
)

// These constants represent the keys used to setup the configuration.
// They are used to access its values using viper.Get.
const (
	DebugKey    = "debug"
	UsernameKey = "username"
	PasswordKey = "password"
	ImagesKey   = "images"
)

// Config has the configurations for the application.
// It's initialized with ReadConfig().
type Config struct {
	Debug    bool
	Username string
	Password string
	Images   []map[string]interface{}
}

func init() {
	viper.SetConfigFile("bob.yaml")
	viper.AddConfigPath(".")

	viper.SetDefault(DebugKey, false)
	viper.SetDefault(UsernameKey, "")
	viper.SetDefault(PasswordKey, "")
	viper.SetDefault(ImagesKey, []interface{}{})
}

// ReadConfig wraps viper.ReadInConfig, but guarantees to setup Viper before reading the config file.
func ReadConfig() (*Config, error) {
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
