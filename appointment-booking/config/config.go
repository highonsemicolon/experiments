package config

type Config struct {
	DBHost     string `koanf:"db.host" validate:"required"`
	DBPort     string `koanf:"db.port" validate:"required"`
	DBUser     string `koanf:"db.user" validate:"required"`
	DBPassword string `koanf:"db.password" validate:"required"`
	DBName     string `koanf:"db.name" validate:"required"`
	ServerPort string `koanf:"server.port" validate:"required"`
}
