package database

import (
	"fmt"
	"log"

	"github.com/clarity/backend/config"
	"github.com/clarity/backend/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database interface {
	GetConnection() *gorm.DB
	Migrate() error
	Close() error
}

type SQLiteDB struct {
	conn *gorm.DB
}

func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
	switch cfg.Type {
	case "sqlite":
		return newSQLiteDB(cfg)
	default:
		return newSQLiteDB(cfg)
	}
}

func newSQLiteDB(cfg *config.DatabaseConfig) (Database, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	log.Printf("Connected to SQLite database at %s", cfg.Path)

	return &SQLiteDB{conn: db}, nil
}

func (s *SQLiteDB) GetConnection() *gorm.DB {
	return s.conn
}

func (s *SQLiteDB) Migrate() error {
	return s.conn.AutoMigrate(
		&models.User{},
		&models.OTPStore{},
		&models.HealthRecord{},
		&models.DoctorConversation{},
	)
}

func (s *SQLiteDB) Close() error {
	sqlDB, err := s.conn.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CloudBackendFactory - placeholder for future cloud support (AWS, GCP, Azure)
type CloudBackendFactory struct {
	Provider string
}

func (cbf *CloudBackendFactory) GetBackend() string {
	switch cbf.Provider {
	case "aws":
		return "AWS RDS"
	case "gcp":
		return "Google Cloud SQL"
	case "azure":
		return "Azure Database"
	default:
		return "Local SQLite"
	}
}
