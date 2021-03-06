package conf

import (
	"flag"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Name  string `split_words:"true" default:"elf"`
	Host  string `split_words:"true" default:"127.0.0.1"`
	Port  int    `split_words:"true" default:"5000"`
	Debug bool   `split_words:"true" default:"true"`
	Log   bool   `split_words:"true" default:"true"`

	DbHost string `split_words:"true" default:"127.0.0.1"`
	DbPort int    `split_words:"true" default:"3306"`
	DbUser string `split_words:"true" default:"root"`
	DbPwd  string `split_words:"true" default:"root"`
	DbName string `split_words:"true" default:"elf-dev"`

	JwtKey        string `split_words:"true" default:"a-secret-key"`
	JwtTimeout    int    `split_words:"true" default:"168"`
	JwtMaxRefresh int    `split_words:"true" default:"24"`

	RedisHost string `split_words:"true" default:"127.0.0.1"`
	RedisPort int    `split_words:"true" default:"6379"`
	RedisPwd  string `split_words:"true" default:"pass"`
	RedisDb   int    `split_words:"true" default:"0"`
}

var cfg *Config
var wg sync.WaitGroup

func ParseConfig() *Config {
	wg.Add(1)
	defer wg.Done()

	cfg = &Config{}

	flag.StringVar(&cfg.Name, "name", "elf", "server name")
	flag.StringVar(&cfg.Host, "host", "127.0.0.1", "server host")
	flag.IntVar(&cfg.Port, "port", 5000, "server port")
	flag.BoolVar(&cfg.Debug, "debug", true, "debug mode")
	flag.BoolVar(&cfg.Log, "log", true, "output access log")

	flag.StringVar(&cfg.DbHost, "db-host", "127.0.0.1", "database host")
	flag.IntVar(&cfg.DbPort, "db-port", 3306, "database port")
	flag.StringVar(&cfg.DbUser, "db-user", "root", "database user")
	flag.StringVar(&cfg.DbPwd, "db-pwd", "root", "database password")
	flag.StringVar(&cfg.DbName, "db-name", "elf-dev", "database name")

	flag.StringVar(&cfg.JwtKey, "jwt-key", "a-secret-key", "jwt key")
	flag.IntVar(&cfg.JwtTimeout, "jwt-timeout", 7*24, "jwt timeout hours")
	flag.IntVar(&cfg.JwtMaxRefresh, "jwt-max-refresh", 24, "jwt max refresh hours")

	flag.StringVar(&cfg.RedisHost, "redis-host", "127.0.0.1", "redis host")
	flag.IntVar(&cfg.RedisPort, "redis-port", 6379, "redis port")
	flag.StringVar(&cfg.RedisPwd, "redis-pwd", "pass", "redis password")
	flag.IntVar(&cfg.RedisDb, "redis-db", 0, "redis db")

	flag.Parse()

	return cfg
}

func ParserConfigFromEnv() *Config {
	wg.Add(1)
	defer wg.Done()

	cfg = &Config{}
	if err := envconfig.Process("elf", cfg); err != nil {
		panic(err)
	}
	return cfg
}

func GetConfig() *Config {
	wg.Wait()
	return cfg
}
