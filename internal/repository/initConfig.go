package repository

import "github.com/spf13/viper"


func initConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.AddConfigPath("C:/Users/User/GolandProjects/CatsCrud/internal/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

