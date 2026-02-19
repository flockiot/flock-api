package config

import (
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Server     ServerConfig     `envPrefix:"SERVER_"`
	Postgres   PostgresConfig   `envPrefix:"POSTGRES_"`
	Redis      RedisConfig      `envPrefix:"REDIS_"`
	ClickHouse ClickHouseConfig `envPrefix:"CLICKHOUSE_"`
	Kafka      KafkaConfig      `envPrefix:"KAFKA_"`
	Loki       LokiConfig       `envPrefix:"LOKI_"`
	Harbor     HarborConfig     `envPrefix:"HARBOR_"`
	S3         S3Config         `envPrefix:"S3_"`
	Zitadel    ZitadelConfig    `envPrefix:"ZITADEL_"`
	Log        LogConfig        `envPrefix:"LOG_"`
}

type ServerConfig struct {
	Host     string `env:"HOST"      envDefault:"0.0.0.0"`
	Port     int    `env:"PORT"      envDefault:"8080"`
	GRPCPort int    `env:"GRPC_PORT" envDefault:"9090"`
}

type PostgresConfig struct {
	DSN string `env:"DSN" envDefault:"postgres://flock:flock@localhost:5432/flock?sslmode=disable"`
}

type RedisConfig struct {
	Addr     string `env:"ADDR"     envDefault:"localhost:6379"`
	Password string `env:"PASSWORD" envDefault:""`
	DB       int    `env:"DB"       envDefault:"0"`
}

type ClickHouseConfig struct {
	DSN string `env:"DSN" envDefault:"clickhouse://flock:flock@localhost:9000/flock"`
}

type KafkaConfig struct {
	Brokers []string `env:"BROKERS" envDefault:"localhost:9092"`
}

type LokiConfig struct {
	Addr string `env:"ADDR" envDefault:"http://localhost:3100"`
}

type HarborConfig struct {
	Addr     string `env:"ADDR"     envDefault:"http://localhost:8081"`
	Username string `env:"USERNAME" envDefault:"admin"`
	Password string `env:"PASSWORD" envDefault:"Harbor12345"`
}

type S3Config struct {
	Endpoint  string `env:"ENDPOINT"   envDefault:"http://localhost:9000"`
	Bucket    string `env:"BUCKET"     envDefault:"flock"`
	AccessKey string `env:"ACCESS_KEY" envDefault:"minioadmin"`
	SecretKey string `env:"SECRET_KEY" envDefault:"minioadmin"`
	Region    string `env:"REGION"     envDefault:"us-east-1"`
}

type ZitadelConfig struct {
	Issuer string `env:"ISSUER" envDefault:"http://localhost:8085"`
	Key    string `env:"KEY"    envDefault:""`
}

type LogConfig struct {
	Level  string `env:"LEVEL"  envDefault:"info"`
	Format string `env:"FORMAT" envDefault:"json"`
}

func Load() (*Config, error) {
	cfg, err := env.ParseAsWithOptions[Config](env.Options{
		Prefix: "FLOCK_",
	})
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
