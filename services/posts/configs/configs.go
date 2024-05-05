package configs

import (
	"context"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Configs struct {
	Port  string
	Host  string
	DbUri string
}

func InitConfigs() *Configs {
	port, host, dbUri := readEnv()

	configs := &Configs{Port: port, Host: host, DbUri: dbUri}
	return configs
}

func (configs *Configs) SetupDB() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(configs.DbUri))

	if err != nil {
		panic(err)
	}

	return client
}

func readEnv() (string, string, string) {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("configs")
	viper.SetConfigType("json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	port := viper.GetString("dev.port")
	host := viper.GetString("dev.host")
	dbUri := viper.GetString("dev.dbUri")

	return port, host, dbUri
}
