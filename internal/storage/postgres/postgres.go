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
						SELECT 1 FROM pg_type WHERE typname = 'role_enum'
				) THEN
						CREATE TYPE role_enum AS ENUM ('admin', 'teacher', 'student');
				END IF;

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
		log.Fatal("Failed to create enums:", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.PasswordReset{},
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
