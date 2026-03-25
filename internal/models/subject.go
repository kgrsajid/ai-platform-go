package models

// Subject represents a school subject with trilingual names
type Subject struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"size:100;not null" json:"name"`
	NameKZ    string `gorm:"size:100" json:"name_kz"`
	NameRU    string `gorm:"size:100" json:"name_ru"`
	Icon      string `gorm:"size:100" json:"icon"`
	Color     string `gorm:"size:7" json:"color"`
	MinGrade  int    `gorm:"default:0" json:"min_grade"`
	MaxGrade  int    `gorm:"default:11" json:"max_grade"`
	IsCore    bool   `gorm:"default:false" json:"is_core"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

// TestSubject is a junction table linking tests to subjects
type TestSubject struct {
	TestID    uint    `gorm:"primaryKey" json:"test_id"`
	Test      Test    `gorm:"foreignKey:TestID" json:"-"`
	SubjectID uint    `gorm:"primaryKey" json:"subject_id"`
	Subject   Subject `gorm:"foreignKey:SubjectID" json:"subject"`
}
