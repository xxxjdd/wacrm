package config

import (
	"context"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	// Load .env file if exists
	godotenv.Load()
	
	// Debug: print what DSN is being used
	dsn := os.Getenv("DB_DSN")
	log.Printf("DB_DSN from env: %s", dsn)
	if dsn == "" {
		dsn = "root:password@tcp(localhost:3306)/wacrm?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	log.Println("Database connected successfully")
	return db, nil
}

var RedisClient *Redis

type Redis struct {
	client interface{}
}

func InitRedis() {
	// For now, we'll skip Redis and use in-memory store
	// In production, connect to Redis for session management
	log.Println("Redis initialization skipped (using in-memory)")
}

func GetRedis() *Redis {
	return RedisClient
}

// SessionStore - in-memory session storage (for development)
type SessionStore struct {
	sessions map[string]*UserSession
}

type UserSession struct {
	UserID    uint
	Username  string
	Email     string
	Token     string
	ExpiresAt time.Time
}

var Sessions = &SessionStore{
	sessions: make(map[string]*UserSession),
}

func (s *SessionStore) Set(token string, session *UserSession) {
	s.sessions[token] = session
}

func (s *SessionStore) Get(token string) (*UserSession, bool) {
	session, ok := s.sessions[token]
	if !ok {
		return nil, false
	}
	if time.Now().After(session.ExpiresAt) {
		delete(s.sessions, token)
		return nil, false
	}
	return session, true
}

func (s *SessionStore) Delete(token string) {
	delete(s.sessions, token)
}

func (s *SessionStore) Cleanup() {
	for token, session := range s.sessions {
		if time.Now().After(session.ExpiresAt) {
			delete(s.sessions, token)
		}
	}
}

// Start session cleanup
func StartSessionCleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			Sessions.Cleanup()
		}
	}()
}

type Database struct {
	*gorm.DB
}

func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

var DBContext = context.Background()

func GetDB() *gorm.DB {
	return DB
}

func MustGetDB() *gorm.DB {
	if DB == nil {
		panic("database not initialized")
	}
	return DB
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func init() {
	// Auto cleanup sessions on startup
	StartSessionCleanup()
}
