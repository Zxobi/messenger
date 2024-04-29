package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Clients  `yaml:"clients"`
	Services `yaml:"services" env-required:"true"`
	Env      string `yaml:"env" env-default:"local"`
	LogLevel string `yaml:"log_level" env-default:"info"`
}

type GrpcConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

type Services struct {
	Auth     AuthConfig     `yaml:"auth"`
	User     UserConfig     `yaml:"user"`
	Chat     ChatConfig     `yaml:"chat"`
	Frontend FrontendConfig `yaml:"frontend"`
}

type AuthConfig struct {
	GrpcConfig `yaml:"grpc"`
	TokenTTL   time.Duration `yaml:"token_ttl" env-default:"1h"`
	Secret     string        `yaml:"secret" env-default:"supa-dupa-secret"`
}

type UserConfig struct {
	GrpcConfig `yaml:"grpc"`
	Storage    UsersMongo `yaml:"storage"`
}

type ChatConfig struct {
	GrpcConfig `yaml:"grpc"`
}

type UsersMongo struct {
	Timeout         time.Duration `yaml:"timeout"`
	ConnectUri      string        `yaml:"connect_uri"`
	DbName          string        `yaml:"db_name"`
	UsersCollection string        `yaml:"users_collection"`
}

type FrontendConfig struct {
	GrpcConfig   `yaml:"grpc"`
	WsPort       int           `yaml:"ws_port"`
	WsBasePath   string        `yaml:"ws_base_path" env-default:"ws"`
	SendBuffSize int           `yaml:"send_buff_size" env-default:"128"`
	RBuffSize    int           `yaml:"read_buff_size" env-default:"4096"`
	WBuffSize    int           `yaml:"write_buff_size" env-default:"4096"`
	HsTimeout    time.Duration `yaml:"hs_timeout" env-default:"30s"`
	MsgLimit     int64         `yaml:"msg_limit" env-default:"4096"`
	WriteWait    time.Duration `yaml:"write_wait" env-default:"5s"`
	PongWait     time.Duration `yaml:"pong_wait" env-default:"5s"`
}

type Clients struct {
	Auth     ClientConfig `yaml:"auth"`
	User     ClientConfig `yaml:"user"`
	Chat     ClientConfig `yaml:"chat"`
	Frontend ClientConfig `yaml:"frontend"`
}

type ClientConfig struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retries_count"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exists: " + path)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
