package postgres

import (
	"fmt"
	"log"
	"project-go/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	const op = "storage.postgres.New"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	err = db.Exec(`
		DO $$
		BEGIN
				IF NOT EXISTS (
						SELECT 1 FROM pg_type WHERE typname = 'difficulty_enum'
				) THEN
						CREATE TYPE difficulty_enum AS ENUM ('easy', 'medium', 'hard');
				END IF;

				IF NOT EXISTS (
						SELECT 1 FROM pg_type WHERE typname = 'status_enum'
				) THEN
						CREATE TYPE status_enum AS ENUM ('pending', 'success', 'error');
				END IF;
		END
		$$;

	`).Error

	if err != nil {
		log.Fatal("Failed to create role_enum:", err)
	}
	err = db.AutoMigrate(
		&models.User{},
		&models.TeacherStudent{},
		&models.SessionHistory{},
		&models.ChatMessage{},
		&models.Test{},
		&models.TestQuestion{},
		&models.TestResult{},
		&models.TestView{},
		&models.CardHolder{},
		&models.Card{},
		&models.Game{},
		&models.Category{},
		&models.GameResult{},
		&models.TestOption{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	return db, nil
}

// func (s *Storage) SaveURL(alias, urlToSave string) (int64, error) {
// 	const op = "storage.postgres.saveURL"

// 	var id int64
// 	// err := s.db.QueryRow(
// 	// 	"INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id",
// 	// 	urlToSave, alias,
// 	// ).Scan(&id)
// 	// if err != nil {
// 	// 	if pqErr, ok := err.(*pq.Error); ok {
// 	// 		switch pqErr.Code {
// 	// 		case "23505": // unique_violation
// 	// 			return 0, fmt.Errorf("duplicate key error on constraint %s: %w", pqErr.Constraint, err)
// 	// 		case "23503": // foreign_key_violation
// 	// 			return 0, fmt.Errorf("foreign key violation on constraint %s: %w", pqErr.Constraint, err)
// 	// 		default:
// 	// 			return 0, fmt.Errorf("pq error (%s): %s", pqErr.Code, pqErr.Message)
// 	// 		}
// 	// 	}
// 	// 	return 0, fmt.Errorf("%s: %w", op, err)
// 	// }

// 	return id, nil
// }

// func (s *Storage) AliasExists(alias string) (bool, error) {
// 	const op = "storage.postgres.aliasExists"
// 	stmt, err := s.db.Prepare("SELECT * from url WHERE alias = $1")
// 	if err != nil {
// 		return false, fmt.Errorf("%s: %w", op, err)
// 	}
// 	var temp string
// 	err = stmt.QueryRow(alias).Scan(&temp)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return false, nil
// 		}
// 		return false, fmt.Errorf("%s: %w", op, err)
// 	}

// 	return true, nil
// }

// func (s *Storage) GetUrl(alias string) (string, error) {
// 	const op = "storage.postgres.getUrl"
// 	stmt, err := s.db.Prepare("SELECT url from url WHERE alias = $1")
// 	if err != nil {
// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}
// 	defer stmt.Close()

// 	var resUrl string
// 	err = stmt.QueryRow(alias).Scan(&resUrl)
// 	if err != nil {
// 		return "", fmt.Errorf("%s: %w", op, err)
// 	}

// 	return resUrl, nil
// }

// func (s *Storage) DeleteUrl(alias string) error {
// 	const op = "storage.postgres.deleteUrl"
// 	stmt, err := s.db.Prepare("DELETE from url WHERE alias = $1")
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}
// 	defer stmt.Close()

// 	_, err = stmt.Exec(alias)
// 	if err != nil {
// 		return fmt.Errorf("%s: %w", op, err)
// 	}

// 	return nil
// }
