package config

type Config struct {
	Kafka
	Postgres
	Redis
}

func Get() (*Config, error) {
	kafka, err := getKafkaConfig()
	if err != nil {
		return nil, err
	}
	postgres, err := getPostgresConfig()
	if err != nil {
		return nil, err
	}
	redis, err := getRedisConfig()
	if err != nil {
		return nil, err
	}
	return &Config{
		Kafka:    kafka,
		Postgres: postgres,
		Redis:    redis,
	}, nil
}
