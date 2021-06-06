package main

import (
	"github.com/hellodoge/delivery-manager/dm"
	"github.com/hellodoge/delivery-manager/internal/cache"
	"github.com/hellodoge/delivery-manager/internal/handler"
	"github.com/hellodoge/delivery-manager/internal/repository"
	"github.com/hellodoge/delivery-manager/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // postgres support
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading .env: %s", err)
	}

	db, err := repository.NewPostgresDB(repository.DatabaseConfig{
		Host:     viper.GetString("db.host"),
		Port:     uint16(viper.GetInt("db.port")),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to init db: %s", err)
	}

	redisConnOptions := cache.RedisConnectionOptions{
		Addr:     viper.GetString("redis.addr"),
		Port:     uint16(viper.GetUint("redis.port")),
		Password: os.Getenv("REDIS_PASSWORD"),
	}

	cacheStorage := cache.NewStorage(redisConnOptions, cache.StorageConfig{
		RTConfig: cache.RefreshTokensConfig{
			Expiration: viper.GetDuration("refresh-tokens.expires"),
			DB:         viper.GetInt("redis.db.refresh-tokens"),
		},
	})

	repo := repository.NewRepository(db)
	services := service.NewService(repo, cacheStorage, service.Config{
		AuthConfig: service.AuthServiceConfig{
			TokenLifetime: viper.GetDuration("jwt.expires"),
			CheckHash:     viper.GetBool("jwt.check-hash"),
		},
	})
	handlers := handler.NewHandler(services)

	server := dm.InitServer(dm.ServerConfig{
		Port:    uint16(viper.GetInt("port")),
		Timeout: viper.GetDuration("timeout"),
	}, handlers.InitRoutes())

	if err := server.Run(); err != nil {
		logrus.Fatalf("error occurred while running http server: %s", err)
	}
}

func initConfig() error {
	viper.SetDefault("port", dm.DefaultPort)
	viper.SetDefault("timeout", dm.DefaultTimeout)
	viper.SetDefault("jwt.expires", service.DefaultTokenLifetime)

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
