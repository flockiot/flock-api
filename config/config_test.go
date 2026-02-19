package config

import (
	"testing"
)

func TestLoadDefaults(t *testing.T) {
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 8080)
	}
	if cfg.Server.GRPCPort != 9090 {
		t.Errorf("Server.GRPCPort = %d, want %d", cfg.Server.GRPCPort, 9090)
	}
	if cfg.Postgres.DSN != "postgres://flock:flock@localhost:5432/flock?sslmode=disable" {
		t.Errorf("Postgres.DSN = %q, want default", cfg.Postgres.DSN)
	}
	if cfg.Redis.Addr != "localhost:6379" {
		t.Errorf("Redis.Addr = %q, want %q", cfg.Redis.Addr, "localhost:6379")
	}
	if cfg.Redis.Password != "" {
		t.Errorf("Redis.Password = %q, want empty", cfg.Redis.Password)
	}
	if cfg.Redis.DB != 0 {
		t.Errorf("Redis.DB = %d, want %d", cfg.Redis.DB, 0)
	}
	if cfg.ClickHouse.DSN != "clickhouse://flock:flock@localhost:9000/flock" {
		t.Errorf("ClickHouse.DSN = %q, want default", cfg.ClickHouse.DSN)
	}
	if len(cfg.Kafka.Brokers) != 1 || cfg.Kafka.Brokers[0] != "localhost:9092" {
		t.Errorf("Kafka.Brokers = %v, want [localhost:9092]", cfg.Kafka.Brokers)
	}
	if cfg.Loki.Addr != "http://localhost:3100" {
		t.Errorf("Loki.Addr = %q, want %q", cfg.Loki.Addr, "http://localhost:3100")
	}
	if cfg.Harbor.Addr != "http://localhost:8081" {
		t.Errorf("Harbor.Addr = %q, want %q", cfg.Harbor.Addr, "http://localhost:8081")
	}
	if cfg.Harbor.Username != "admin" {
		t.Errorf("Harbor.Username = %q, want %q", cfg.Harbor.Username, "admin")
	}
	if cfg.S3.Endpoint != "http://localhost:9000" {
		t.Errorf("S3.Endpoint = %q, want %q", cfg.S3.Endpoint, "http://localhost:9000")
	}
	if cfg.S3.Bucket != "flock" {
		t.Errorf("S3.Bucket = %q, want %q", cfg.S3.Bucket, "flock")
	}
	if cfg.S3.AccessKey != "minioadmin" {
		t.Errorf("S3.AccessKey = %q, want %q", cfg.S3.AccessKey, "minioadmin")
	}
	if cfg.S3.Region != "us-east-1" {
		t.Errorf("S3.Region = %q, want %q", cfg.S3.Region, "us-east-1")
	}
	if cfg.Zitadel.Issuer != "http://localhost:8085" {
		t.Errorf("Zitadel.Issuer = %q, want %q", cfg.Zitadel.Issuer, "http://localhost:8085")
	}
	if cfg.Zitadel.Key != "" {
		t.Errorf("Zitadel.Key = %q, want empty", cfg.Zitadel.Key)
	}
	if cfg.Log.Level != "info" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "info")
	}
	if cfg.Log.Format != "text" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "text")
	}
}

func TestLoadEnvOverrides(t *testing.T) {
	t.Setenv("FLOCK_SERVER_HOST", "127.0.0.1")
	t.Setenv("FLOCK_SERVER_PORT", "3000")
	t.Setenv("FLOCK_SERVER_GRPC_PORT", "3001")
	t.Setenv("FLOCK_POSTGRES_DSN", "postgres://prod:secret@db:5432/flock")
	t.Setenv("FLOCK_REDIS_ADDR", "redis:6379")
	t.Setenv("FLOCK_REDIS_PASSWORD", "s3cret")
	t.Setenv("FLOCK_REDIS_DB", "2")
	t.Setenv("FLOCK_CLICKHOUSE_DSN", "clickhouse://prod:secret@ch:9000/flock")
	t.Setenv("FLOCK_LOKI_ADDR", "http://loki:3100")
	t.Setenv("FLOCK_HARBOR_ADDR", "https://harbor.example.com")
	t.Setenv("FLOCK_HARBOR_USERNAME", "robot")
	t.Setenv("FLOCK_HARBOR_PASSWORD", "token123")
	t.Setenv("FLOCK_S3_ENDPOINT", "https://s3.example.com")
	t.Setenv("FLOCK_S3_BUCKET", "prod-flock")
	t.Setenv("FLOCK_S3_ACCESS_KEY", "AKIA123")
	t.Setenv("FLOCK_S3_SECRET_KEY", "secret456")
	t.Setenv("FLOCK_S3_REGION", "eu-west-1")
	t.Setenv("FLOCK_ZITADEL_ISSUER", "https://auth.example.com")
	t.Setenv("FLOCK_ZITADEL_KEY", "zkey")
	t.Setenv("FLOCK_LOG_LEVEL", "debug")
	t.Setenv("FLOCK_LOG_FORMAT", "json")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Errorf("Server.Host = %q, want %q", cfg.Server.Host, "127.0.0.1")
	}
	if cfg.Server.Port != 3000 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 3000)
	}
	if cfg.Server.GRPCPort != 3001 {
		t.Errorf("Server.GRPCPort = %d, want %d", cfg.Server.GRPCPort, 3001)
	}
	if cfg.Postgres.DSN != "postgres://prod:secret@db:5432/flock" {
		t.Errorf("Postgres.DSN = %q, want override", cfg.Postgres.DSN)
	}
	if cfg.Redis.Addr != "redis:6379" {
		t.Errorf("Redis.Addr = %q, want %q", cfg.Redis.Addr, "redis:6379")
	}
	if cfg.Redis.Password != "s3cret" {
		t.Errorf("Redis.Password = %q, want %q", cfg.Redis.Password, "s3cret")
	}
	if cfg.Redis.DB != 2 {
		t.Errorf("Redis.DB = %d, want %d", cfg.Redis.DB, 2)
	}
	if cfg.ClickHouse.DSN != "clickhouse://prod:secret@ch:9000/flock" {
		t.Errorf("ClickHouse.DSN = %q, want override", cfg.ClickHouse.DSN)
	}
	if cfg.Loki.Addr != "http://loki:3100" {
		t.Errorf("Loki.Addr = %q, want %q", cfg.Loki.Addr, "http://loki:3100")
	}
	if cfg.Harbor.Addr != "https://harbor.example.com" {
		t.Errorf("Harbor.Addr = %q, want %q", cfg.Harbor.Addr, "https://harbor.example.com")
	}
	if cfg.Harbor.Username != "robot" {
		t.Errorf("Harbor.Username = %q, want %q", cfg.Harbor.Username, "robot")
	}
	if cfg.Harbor.Password != "token123" {
		t.Errorf("Harbor.Password = %q, want %q", cfg.Harbor.Password, "token123")
	}
	if cfg.S3.Endpoint != "https://s3.example.com" {
		t.Errorf("S3.Endpoint = %q, want %q", cfg.S3.Endpoint, "https://s3.example.com")
	}
	if cfg.S3.Bucket != "prod-flock" {
		t.Errorf("S3.Bucket = %q, want %q", cfg.S3.Bucket, "prod-flock")
	}
	if cfg.S3.AccessKey != "AKIA123" {
		t.Errorf("S3.AccessKey = %q, want %q", cfg.S3.AccessKey, "AKIA123")
	}
	if cfg.S3.SecretKey != "secret456" {
		t.Errorf("S3.SecretKey = %q, want %q", cfg.S3.SecretKey, "secret456")
	}
	if cfg.S3.Region != "eu-west-1" {
		t.Errorf("S3.Region = %q, want %q", cfg.S3.Region, "eu-west-1")
	}
	if cfg.Zitadel.Issuer != "https://auth.example.com" {
		t.Errorf("Zitadel.Issuer = %q, want %q", cfg.Zitadel.Issuer, "https://auth.example.com")
	}
	if cfg.Zitadel.Key != "zkey" {
		t.Errorf("Zitadel.Key = %q, want %q", cfg.Zitadel.Key, "zkey")
	}
	if cfg.Log.Level != "debug" {
		t.Errorf("Log.Level = %q, want %q", cfg.Log.Level, "debug")
	}
	if cfg.Log.Format != "json" {
		t.Errorf("Log.Format = %q, want %q", cfg.Log.Format, "json")
	}
}

func TestLoadKafkaBrokersCommaSeparated(t *testing.T) {
	t.Setenv("FLOCK_KAFKA_BROKERS", "broker1:9092,broker2:9092,broker3:9092")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	expected := []string{"broker1:9092", "broker2:9092", "broker3:9092"}
	if len(cfg.Kafka.Brokers) != len(expected) {
		t.Fatalf("Kafka.Brokers length = %d, want %d", len(cfg.Kafka.Brokers), len(expected))
	}
	for i, got := range cfg.Kafka.Brokers {
		if got != expected[i] {
			t.Errorf("Kafka.Brokers[%d] = %q, want %q", i, got, expected[i])
		}
	}
}

func TestLoadPartialOverride(t *testing.T) {
	t.Setenv("FLOCK_SERVER_PORT", "4000")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.Server.Port != 4000 {
		t.Errorf("Server.Port = %d, want %d", cfg.Server.Port, 4000)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %q, want default %q", cfg.Server.Host, "0.0.0.0")
	}
	if cfg.Server.GRPCPort != 9090 {
		t.Errorf("Server.GRPCPort = %d, want default %d", cfg.Server.GRPCPort, 9090)
	}
}
