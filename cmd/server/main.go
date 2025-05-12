package main

import (
	"fmt"
	"log"
	"time"
	"ewa/internal/handler"
	"ewa/internal/middleware"
	"ewa/internal/model"
	"ewa/internal/repository"
	"ewa/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// 连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(&model.User{}, &model.Model{}, &model.UserModelPermission{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 初始化Redis服务
	redisConfig := &service.RedisConfig{
		Host:            viper.GetString("redis.host"),
		Port:            viper.GetInt("redis.port"),
		Password:        viper.GetString("redis.password"),
		DB:              viper.GetInt("redis.db"),
		PoolSize:        viper.GetInt("redis.pool_size"),
		MinIdleConns:    viper.GetInt("redis.min_idle_conns"),
		PublicModelTTL:  time.Duration(viper.GetInt("redis.cache.public_model_ttl")) * time.Second,
		UserInfoTTL:     time.Duration(viper.GetInt("redis.cache.user_info_ttl")) * time.Second,
		PermissionTTL:   time.Duration(viper.GetInt("redis.cache.permission_ttl")) * time.Second,
		RateLimitTTL:    time.Duration(viper.GetInt("redis.cache.rate_limit_ttl")) * time.Second,
	}
	redisService, err := service.NewRedisService(redisConfig)
	if err != nil {
		log.Fatalf("Failed to initialize Redis service: %v", err)
	}

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	modelRepo := repository.NewModelRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)

	// 初始化OSS服务
	ossService := service.NewAliyunOSSService(
		viper.GetString("oss.endpoint"),
		viper.GetString("oss.access_key_id"),
		viper.GetString("oss.access_key_secret"),
		viper.GetString("oss.bucket"),
		viper.GetString("oss.region"),
		viper.GetInt64("oss.url_expire"),
		viper.GetInt64("oss.upload_expire"),
	)

	// 初始化服务
	userService := service.NewUserService(userRepo, redisService)
	modelService := service.NewModelService(modelRepo, permissionRepo, ossService, redisService)
	permissionService := service.NewPermissionService(permissionRepo, redisService)

	// 初始化处理器
	userHandler := handler.NewUserHandler(userService, viper.GetString("jwt.secret"), viper.GetDuration("jwt.expire"))
	modelHandler := handler.NewModelHandler(modelService)
	permissionHandler := handler.NewPermissionHandler(permissionService)

	// 初始化路由
	r := gin.Default()

	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/register", userHandler.Register)
		public.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(viper.GetString("jwt.secret")))
	{
		// 用户相关
		auth.PUT("/profile", userHandler.UpdateProfile)

		// 模型相关
		auth.POST("/models", modelHandler.Create)
		auth.PUT("/models/:id", modelHandler.Update)
		auth.DELETE("/models/:id", modelHandler.Delete)
		auth.GET("/models/:id", modelHandler.Get)
		auth.GET("/models", modelHandler.List)
		auth.GET("/models/:id/file", modelHandler.GetModelFile)
		auth.POST("/models/upload-url", modelHandler.GetUploadURL)

		// 权限相关
		auth.POST("/permissions", permissionHandler.Grant)
		auth.PUT("/permissions/:id", permissionHandler.Update)
		auth.DELETE("/permissions/:id", permissionHandler.Revoke)
		auth.GET("/models/:model_id/permissions", permissionHandler.List)
	}

	// 启动服务器
	port := viper.GetString("server.port")
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 