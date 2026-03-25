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
		&models.TestOption{},
		&models.CardHolder{},
		&models.Card{},
		&models.Game{},
		&models.Category{},
		&models.GameResult{},
		&models.Subject{},
		&models.TestSubject{},
		&models.UserProgress{},
		&models.PointTransaction{},
		&models.DailyActivity{},
		&models.Reward{},
		&models.Redemption{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedCategories(db)
	seedSubjects(db)
	seedRewards(db)

	return db, nil
}

func seedSubjects(db *gorm.DB) {
	var count int64
	db.Model(&models.Subject{}).Count(&count)
	if count > 0 {
		return
	}

	subjects := []models.Subject{
		// Primary School (0-4)
		{Name: "Mathematics", NameKZ: "Математика", NameRU: "Математика", Icon: "🔢", Color: "#FF6B6B", MinGrade: 0, MaxGrade: 11, IsCore: true, SortOrder: 1},
		{Name: "Kazakh Language", NameKZ: "Қазақ тілі", NameRU: "Казахский язык", Icon: "🇰🇿", Color: "#4ECDC4", MinGrade: 0, MaxGrade: 11, IsCore: true, SortOrder: 2},
		{Name: "Russian Language", NameKZ: "Орыс тілі", NameRU: "Русский язык", Icon: "🇷🇺", Color: "#95E1D3", MinGrade: 0, MaxGrade: 11, IsCore: true, SortOrder: 3},
		{Name: "English Language", NameKZ: "Ағылшын тілі", NameRU: "Английский язык", Icon: "🇬🇧", Color: "#FFE66D", MinGrade: 1, MaxGrade: 11, IsCore: false, SortOrder: 4},
		{Name: "Natural Science", NameKZ: "Табиғаттану", NameRU: "Природоведение", Icon: "🌿", Color: "#A8E6CF", MinGrade: 1, MaxGrade: 4, IsCore: false, SortOrder: 5},
		{Name: "Art", NameKZ: "Бейнелеу өнері", NameRU: "Изобразительное искусство", Icon: "🎨", Color: "#FF8B94", MinGrade: 0, MaxGrade: 4, IsCore: false, SortOrder: 6},
		{Name: "Music", NameKZ: "Музыка", NameRU: "Музыка", Icon: "🎵", Color: "#DDA0DD", MinGrade: 0, MaxGrade: 4, IsCore: false, SortOrder: 7},
		// Basic Secondary (5-9)
		{Name: "Algebra", NameKZ: "Алгебра", NameRU: "Алгебра", Icon: "📐", Color: "#6C5CE7", MinGrade: 5, MaxGrade: 9, IsCore: true, SortOrder: 8},
		{Name: "Geometry", NameKZ: "Геометрия", NameRU: "Геометрия", Icon: "📏", Color: "#A29BFE", MinGrade: 7, MaxGrade: 9, IsCore: true, SortOrder: 9},
		{Name: "Kazakh Literature", NameKZ: "Қазақ әдебиеті", NameRU: "Казахская литература", Icon: "📖", Color: "#74B9FF", MinGrade: 5, MaxGrade: 11, IsCore: true, SortOrder: 10},
		{Name: "Russian Literature", NameKZ: "Орыс әдебиеті", NameRU: "Русская литература", Icon: "📚", Color: "#81ECEC", MinGrade: 5, MaxGrade: 11, IsCore: true, SortOrder: 11},
		{Name: "Biology", NameKZ: "Биология", NameRU: "Биология", Icon: "🧬", Color: "#55E6C1", MinGrade: 5, MaxGrade: 11, IsCore: false, SortOrder: 12},
		{Name: "Chemistry", NameKZ: "Химия", NameRU: "Химия", Icon: "⚗️", Color: "#FECA57", MinGrade: 7, MaxGrade: 11, IsCore: false, SortOrder: 13},
		{Name: "Physics", NameKZ: "Физика", NameRU: "Физика", Icon: "⚡", Color: "#48DBFB", MinGrade: 7, MaxGrade: 11, IsCore: false, SortOrder: 14},
		{Name: "History of Kazakhstan", NameKZ: "Қазақстан тарихы", NameRU: "История Казахстана", Icon: "📜", Color: "#FF9FF3", MinGrade: 5, MaxGrade: 11, IsCore: true, SortOrder: 15},
		{Name: "World History", NameKZ: "Дүниежүзі тарихы", NameRU: "Всемирная история", Icon: "🌍", Color: "#F368E0", MinGrade: 5, MaxGrade: 9, IsCore: false, SortOrder: 16},
		{Name: "Geography", NameKZ: "География", NameRU: "География", Icon: "🗺️", Color: "#1DD1A1", MinGrade: 6, MaxGrade: 11, IsCore: false, SortOrder: 17},
		{Name: "Informatics", NameKZ: "Информатика", NameRU: "Информатика", Icon: "💻", Color: "#54A0FF", MinGrade: 5, MaxGrade: 11, IsCore: false, SortOrder: 18},
		// Senior (10-11)
		{Name: "Advanced Mathematics", NameKZ: "Математика (жетілдірілген)", NameRU: "Математика (углубленная)", Icon: "📊", Color: "#C44569", MinGrade: 10, MaxGrade: 11, IsCore: true, SortOrder: 19},
		{Name: "UNT Prep", NameKZ: "ҰБТ дайындық", NameRU: "Подготовка к ЕНТ", Icon: "🎯", Color: "#F97F51", MinGrade: 11, MaxGrade: 11, IsCore: false, SortOrder: 20},
	}

	if err := db.Create(&subjects).Error; err != nil {
		log.Println("Failed to seed subjects:", err)
		return
	}
	log.Println("Subjects seeded successfully:", count, "existing,", len(subjects), "new")
}

func seedCategories(db *gorm.DB) {
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count > 0 {
		return
	}

	categories := []models.Category{
		{Name: "Математика"},
		{Name: "Информатика"},
		{Name: "Физика"},
		{Name: "Программирование"},
		{Name: "Frontend"},
		{Name: "Backend"},
		{Name: "Электроника"},
		{Name: "Компьютеры"},
		{Name: "Алгоритмы"},
		{Name: "Базы данных"},
		{Name: "Сети"},
	}

	if err := db.Create(&categories).Error; err != nil {
		log.Println("Failed to seed categories:", err)
		return
	}
	log.Println("Categories seeded successfully")
}

func seedRewards(db *gorm.DB) {
	var count int64
	db.Model(&models.Reward{}).Count(&count)
	if count > 0 {
		return
	}

	rewards := []models.Reward{
		// Food (0-4 Sprouts)
		{Title: "Dodo Pizza Slice", Description: "Free pizza slice at Dodo Pizza", Category: "food", PartnerName: "Dodo Pizza", PointCost: 50, MinGrade: 0, MaxGrade: 11, TotalStock: 100, IsActive: true},
		{Title: "Popeyes Chicken Meal", Description: "Chicken combo at Popeyes", Category: "food", PartnerName: "Popeyes", PointCost: 200, MinGrade: 0, MaxGrade: 11, TotalStock: 50, IsActive: true},
		{Title: "McDonald's Happy Meal", Description: "Happy Meal at McDonald's", Category: "food", PartnerName: "McDonald's", PointCost: 100, MinGrade: 0, MaxGrade: 9, TotalStock: 75, IsActive: true},
		{Title: "Starbucks Drink", Description: "Any drink at Starbucks", Category: "food", PartnerName: "Starbucks", PointCost: 150, MinGrade: 5, MaxGrade: 11, TotalStock: 60, IsActive: true},

		// Retail (0-4)
		{Title: "School Supplies Pack", Description: "Notebooks, pens, and pencils", Category: "retail", PartnerName: "Barnes & Noble", PointCost: 80, MinGrade: 0, MaxGrade: 4, TotalStock: 40, IsActive: true},
		{Title: "Book Voucher", Description: "Any book up to $15", Category: "retail", PartnerName: "Amazon", PointCost: 120, MinGrade: 0, MaxGrade: 11, TotalStock: 30, IsActive: true},

		// Entertainment (5-9)
		{Title: "Movie Ticket", Description: "Free movie ticket at Cinema City", Category: "entertainment", PartnerName: "Cinema City", PointCost: 250, MinGrade: 5, MaxGrade: 11, TotalStock: 25, IsActive: true},
		{Title: "Karaoke Voucher", Description: "2 hours of karaoke", Category: "entertainment", PartnerName: "Karaoke Center", PointCost: 300, MinGrade: 5, MaxGrade: 11, TotalStock: 15, IsActive: true},

		// Education (10-11)
		{Title: "STEM Workshop", Description: "One STEM workshop session", Category: "education", PartnerName: "STEM Lab", PointCost: 400, MinGrade: 5, MaxGrade: 11, TotalStock: 10, IsActive: true},
		{Title: "Tutoring Session", Description: "1-on-1 tutoring session", Category: "education", PartnerName: "SDU AI School", PointCost: 500, MinGrade: 5, MaxGrade: 11, TotalStock: 20, IsActive: true},

		// Virtual (all grades)
		{Title: "Custom Avatar", Description: "Unlock a unique avatar skin", Category: "virtual", PartnerName: "SDU AI School", PointCost: 75, MinGrade: 0, MaxGrade: 11, TotalStock: -1, IsActive: true},
		{Title: "Robot Upgrade", Description: "Upgrade your robot to next level", Category: "virtual", PartnerName: "SDU AI School", PointCost: 200, MinGrade: 0, MaxGrade: 11, TotalStock: -1, IsActive: true},
		{Title: "Gems × 50", Description: "50 virtual gems for the shop", Category: "virtual", PartnerName: "SDU AI School", PointCost: 100, MinGrade: 0, MaxGrade: 11, TotalStock: -1, IsActive: true},
	}

	if err := db.Create(&rewards).Error; err != nil {
		log.Println("Failed to seed rewards:", err)
		return
	}
	log.Println("Rewards seeded successfully:", len(rewards), "rewards added")
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
