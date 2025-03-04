package connector

import (
	"errors"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Connector struct {
	db *gorm.DB
}

func NewDBConnector(dbtype, dsn string) (*Connector, error) {
	db, err := initDB(dbtype, dsn)
	if err != nil {
		return nil, err
	}

	return &Connector{
		db: db,
	}, nil
}

func initDB(dbtype, dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	switch dbtype {
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgres":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, errors.New("unknown database type")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %w", dbtype, err)
	}

	if db == nil {
		return nil, fmt.Errorf("failed to initialize database connection for %s", dbtype)
	}

	return db, nil
}

func (c *Connector) AutoMigrate(models ...interface{}) error {
	return c.db.AutoMigrate(models...)
}

func (c *Connector) DB() *gorm.DB {
	return c.db
}
