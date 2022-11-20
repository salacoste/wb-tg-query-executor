package config

import "github.com/spf13/viper"

type Config struct {
	Database struct {
		Host   	  string
		Name	  string
		User	  string
		Password  string
		Port	  int
	}
}

const (
	kDefaultPostgresPort = 5432
)

func LoadConfig(path string) (config Config, err error) {
    viper.AddConfigPath(path)
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    viper.AutomaticEnv()
	viper.BindEnv("Database.Host", "DATABASE_HOST")
	viper.BindEnv("Database.Name", "DATABASE_NAME")
	viper.BindEnv("Database.User", "DATABASE_USER")
	viper.BindEnv("Database.Password", "DATABASE_PASSWORD")
	viper.BindEnv("Database.Port", "DATABASE_PORT")
	
	viper.SetDefault("Database.Port", kDefaultPostgresPort)

    viper.ReadInConfig()
    err = viper.Unmarshal(&config)
    return
}