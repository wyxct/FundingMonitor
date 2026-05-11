package config

import (
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type AppConfig struct {
	Port string `mapstructure:"port"`
}

type PostgresConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

type ChainConfig struct {
	RpcUrl     string `mapstructure:"rpc"`
	Contract   string `mapstructure:"contract"`
	StartBlock int64  `mapstructure:"startblock"`
}

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Postgres PostgresConfig `mapstructure:"postgres"`
	Chain    ChainConfig    `mapstructure:"chain"`
}

var C Config

// 全局DB对象，整个项目共用
var DB *gorm.DB

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic("读取配置文件失败: " + err.Error())
	}

	err = viper.Unmarshal(&C)
	if err != nil {
		panic("解析配置文件失败: " + err.Error())
	}
}

func getDsn() string {
	return "host=" + C.Postgres.Host + " port=" + C.Postgres.Port + " user=" + C.Postgres.User + " password=" + C.Postgres.Password + " dbname=" + C.Postgres.DBName + " sslmode=" + C.Postgres.SSLMode
}

// InitDB 初始化数据库（主函数只调用我）
func InitDB() {
	dsn := getDsn()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败: " + err.Error())
	}

	DB = db
	println("✅ 数据库初始化成功")
}
