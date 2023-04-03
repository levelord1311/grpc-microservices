package config

import (
	"github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/database"
	"time"
)

var _ database.Config = &Database{}

type Database struct {
	DSN             string        `yaml:"DSN"`
	maxOpenConns    int           `yaml:"maxOpenConns"`
	maxIdleConns    int           `yaml:"maxIdleConns"`
	connMaxIdleTime time.Duration `yaml:"connMaxIdleTime"`
	connMaxLifeTime time.Duration `yaml:"connMaxLifeTime"`
	Migrations      string        `yaml:"migrations"`
}

func (d Database) GetDSN() string {
	return d.DSN
}

func (d Database) GetMaxOpenConns() int {
	return d.maxOpenConns
}

func (d Database) GetMaxIdleConns() int {
	return d.maxIdleConns
}

func (d Database) GetConnMaxIdleTime() time.Duration {
	return d.connMaxIdleTime
}

func (d Database) GetConnMaxLifetime() time.Duration {
	return d.connMaxLifeTime
}
