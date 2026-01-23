package config

import "fmt"

type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
	SSLMode  string `env:"DB_SSLMODE" envDefault:"disable"`
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		d.User,
		d.Password,
		d.Host,
		d.Name,
		d.SSLMode,
	)
}
