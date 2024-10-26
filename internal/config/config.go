package config

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	EnvLocal = "local"
	EnvProd  = "prod"
)

type Config struct {
	Env           string        `yaml:"env" env:"ENV" env-default:"local"`
	SqlConnection SqlConnection `yaml:"database"`
	AuthConfig    AuthConfig    `yaml:"auth"`
	GRPC GRPCConfig `yaml:"grpc"`
}

type SqlConnection struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port"`
	Name     string `yaml:"name"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

type AuthConfig struct {
	JWT          JWTConfig `yaml:"jwt"`
	PasswordSalt string    `yaml:"password_salt"`
}
type JWTConfig struct {
	AccessTokenTTL  time.Duration `yaml:"access_token_ttl"`
	RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl"`
	Secret          string        `yaml:"secret"`
}

type GRPCConfig struct{
	Port int `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}
func MustLoad() *Config{
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}
	return MustLoadPath(configPath)
}
func MustLoadPath(filePath string) *Config {
	if filePath == "" {
		log.Fatal("CONFIG_PATH not set")
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist\n", filePath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		log.Fatalf("cannot read config %s", filePath)
	}
	return &cfg
}

func fetchConfigPath() string{
	var res string
	flag.StringVar(&res, "config","","path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}
	return res
}