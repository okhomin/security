package config

import "github.com/spf13/viper"

type Config struct {
	Port      string
	JWTKey    string
	RootLogin string
	DBURL     string
	Cost      int
	Pepper    string
}

func New() Config {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8888")
	viper.SetDefault("ROOT_LOGIN", "root")
	viper.SetDefault("JWT_KEY", "jwt_key")
	viper.SetDefault("COST", 13)
	viper.Set("PEPPER", "pepper")
	viper.SetDefault("DB_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")

	return Config{
		Port:      viper.GetString("PORT"),
		JWTKey:    viper.GetString("JWT_KEY"),
		RootLogin: viper.GetString("ROOT_LOGIN"),
		DBURL:     viper.GetString("DB_URL"),
		Cost:      viper.GetInt("COST"),
		Pepper:    viper.GetString("PEPPER"),
	}
}
