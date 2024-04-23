package storage

import (
	"database/sql"
	"log"
	_ "modernc.org/sqlite"
	"os"
)

type Storage struct {
	config         *Config
	db             *sql.DB
	taskRepository *TaskRepository
}

func New(config *Config) *Storage {
	return &Storage{
		config: config,
	}
}

func (s *Storage) Open() error {
	os.Remove(s.config.Dsn)

	db, err := sql.Open("sqlite", s.config.Dsn)
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS task (id INTEGER NOT NULL PRIMARY KEY autoincrement, task TEXT, tags TEXT, due datetime);
		DELETE FROM task;
		`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	s.db = db

	log.Println("DB connection created successfully")
	log.Println("DB file is", s.config.Dsn)

	return nil
}

func (s *Storage) Task() *TaskRepository {
	if s.taskRepository == nil {
		s.taskRepository = &TaskRepository{
			storage: s,
		}
	}

	return s.taskRepository
}
