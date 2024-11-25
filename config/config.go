package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	logging "github.com/haitao-sun03/logging/config"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

var config Config
var DB *gorm.DB
var RedisClient *redis.Client

func Init() {
	// 设置viper读取配置文件
	viper.SetConfigName("config")  // 配置文件的名称（不需要后缀）
	viper.SetConfigType("yaml")    // 配置文件的类型
	viper.AddConfigPath("config/") // 配置文件所在的路径

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	// 解析配置到结构体
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	// 解析出日志配置
	var loggingConfig logging.LoggingConfig
	if err := viper.UnmarshalKey("logging", &loggingConfig); err != nil {
		panic(fmt.Errorf("unable to decode into struct, %v", err))
	}

	InitDatabase()
	InitRedis()
	//将日志配置传递给日志模块并初始化
	logging.InitLogging(loggingConfig)
}

func InitDatabase() {

	// 使用配置初始化GORM数据库连接
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB.User,
		config.DB.Password,
		config.DB.Host,
		config.DB.Port,
		config.DB.DBName,
		// config.DB.SSLMode,
	)
	// 配置 GORM Logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 使用标准输出作为日志写入器
		logger.Config{
			LogLevel:      logger.Info,             // 设置日志级别为 Info
			SlowThreshold: 1000 * time.Millisecond, // 慢 SQL 阈值
			Colorful:      true,                    // 启用彩色输出
		},
	)

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Printf("failed to connect database: %v", err)
	}
	if DB.Error != nil {
		fmt.Print("DB.Error")
	}

	fmt.Println("init DB success : ", *DB)

	// 自动迁移模式
	// DB.AutoMigrate(&User{})
}

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.Redis.Address,
		Password: config.Redis.Password, // 没有密码可以为空字符串
		DB:       config.Redis.DB,       // 使用默认DB 0
	})

	// 测试连接
	ctx := context.Background()
	pong, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis connection successful:", pong)

}
