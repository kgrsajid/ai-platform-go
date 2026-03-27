package models

import "fmt"

// GradeBand represents the three age-appropriate tiers
type GradeBand string

const (
	BandSprouts   GradeBand = "sprouts"   // Grades 0-4 (ages 6-10)
	BandExplorers GradeBand = "explorers" // Grades 5-9 (ages 10-14)
	BandChampions GradeBand = "champions" // Grades 10-11 (ages 15-17)
)

// GetGradeBand returns the grade band for a given grade level
func GetGradeBand(grade int) GradeBand {
	switch {
	case grade >= 0 && grade <= 4:
		return BandSprouts
	case grade >= 5 && grade <= 9:
		return BandExplorers
	case grade >= 10 && grade <= 11:
		return BandChampions
	default:
		return BandExplorers
	}
}

// CurrencyLabel returns the display name for points in this band
func (b GradeBand) CurrencyLabel() string {
	switch b {
	case BandSprouts:
		return "Stars"
	case BandExplorers:
		return "Gems"
	case BandChampions:
		return "Coins"
	default:
		return "Points"
	}
}

// CurrencyIcon returns the emoji for this band's currency
func (b GradeBand) CurrencyIcon() string {
	switch b {
	case BandSprouts:
		return "⭐"
	case BandExplorers:
		return "💎"
	case BandChampions:
		return "🪙"
	default:
		return "💰"
	}
}

// LevelName returns the display name for a level within this band
func (b GradeBand) LevelName(level int) string {
	switch b {
	case BandSprouts:
		names := []string{"", "Seedling", "Sprout", "Sapling", "Little Tree", "Big Tree",
			"Strong Oak", "Ancient Oak", "Great Oak", "Mighty Oak", "Oak Legend"}
		if level >= 1 && level <= len(names)-1 {
			return names[level]
		}
		return fmt.Sprintf("Level %d", level)
	case BandExplorers:
		names := []string{"", "Novice", "Learner", "Scholar", "Explorer", "Pathfinder",
			"Navigator", "Adventurer", "Trailblazer", "Pioneer", "Vanguard",
			"Sentinel", "Warden", "Ranger", "Champion", "Hero",
			"Legend", "Mythic", "Epic", "Mythical", "Transcendent",
			"Ascendant", "Celestial", "Cosmic", "Infinite", "Omega"}
		if level >= 1 && level <= len(names)-1 {
			return names[level]
		}
		return fmt.Sprintf("Level %d", level)
	case BandChampions:
		names := []string{"", "Bronze", "Bronze+", "Silver", "Silver+", "Silver Elite",
			"Gold", "Gold+", "Gold Elite", "Platinum", "Platinum+",
			"Diamond", "Diamond+", "Diamond Elite", "Master", "Master+",
			"Grandmaster", "Grandmaster+", "Legend", "Legend+", "Mythic",
			"Mythic+", "Transcendent", "Ascendant", "Ascendant+", "Celestial",
			"Cosmic", "Cosmic+", "Omega", "Omega+", "Alpha", "Alpha+",
			"Supreme", "Supreme+", "Ultimate", "Ultimate+", "Apex",
			"Apex+", "Godlike", "Godlike+", "Immortal", "Immortal+",
			"Eternal", "Eternal+", "Infinite", "Infinite+", "Absolute"}
		if level >= 1 && level <= len(names)-1 {
			return names[level]
		}
		return fmt.Sprintf("Level %d", level)
	default:
		return fmt.Sprintf("Level %d", level)
	}
}

// MaxLevel returns the maximum level for this band
func (b GradeBand) MaxLevel() int {
	switch b {
	case BandSprouts:
		return 10
	case BandExplorers:
		return 25
	case BandChampions:
		return 50
	default:
		return 25
	}
}

// XPPerLevel returns the XP needed per level in this band
func (b GradeBand) XPPerLevel() int {
	switch b {
	case BandSprouts:
		return 20
	case BandExplorers:
		return 100
	case BandChampions:
		return 500
	default:
		return 100
	}
}

// QuizPoints calculates points earned for a quiz completion
func (b GradeBand) QuizPoints(percentage float64) int {
	switch b {
	case BandSprouts:
		// 5-15 points based on score
		return int(5 + percentage*0.1)
	case BandExplorers:
		// 10-50 points based on score
		return int(10 + percentage*0.4)
	case BandChampions:
		// 20-50 points based on score
		return int(20 + percentage*0.3)
	default:
		return int(10 + percentage*0.4)
	}
}

// FlashcardPoints calculates points earned for a flashcard session
func (b GradeBand) FlashcardPoints(cardsStudied int) int {
	switch b {
	case BandSprouts:
		// 3-5 points per session
		if cardsStudied >= 20 {
			return 5
		}
		return 3
	case BandExplorers:
		// 5-15 points
		return min(5+cardsStudied/2, 15)
	case BandChampions:
		// 10-25 points
		return min(10+cardsStudied, 25)
	default:
		return min(5+cardsStudied/2, 15)
	}
}

// QuizXP calculates XP earned for a quiz completion
func (b GradeBand) QuizXP(percentage float64) int {
	switch b {
	case BandSprouts:
		return int(percentage * 0.5)
	case BandExplorers:
		return int(percentage)
	case BandChampions:
		return int(percentage * 1.5)
	default:
		return int(percentage)
	}
}

// FlashcardXP calculates XP earned for a flashcard session
func (b GradeBand) FlashcardXP(cardsStudied int) int {
	switch b {
	case BandSprouts:
		return cardsStudied
	case BandExplorers:
		return cardsStudied * 2
	case BandChampions:
		return cardsStudied * 3
	default:
		return cardsStudied * 2
	}
}

// DailyLoginXP returns XP for daily login
func (b GradeBand) DailyLoginXP() int {
	switch b {
	case BandSprouts:
		return 5
	case BandExplorers:
		return 10
	case BandChampions:
		return 25
	default:
		return 10
	}
}

// DailyLoginPoints returns points for daily login
func (b GradeBand) DailyLoginPoints() int {
	switch b {
	case BandSprouts:
		return 3
	case BandExplorers:
		return 5
	case BandChampions:
		return 10
	default:
		return 5
	}
}

// StreakBonusPoints returns bonus points at a given streak day
func (b GradeBand) StreakBonusPoints(streakDay int) int {
	if streakDay%7 != 0 {
		return 0
	}
	switch b {
	case BandSprouts:
		return 5
	case BandExplorers:
		return 10
	case BandChampions:
		return 20
	default:
		return 10
	}
}

// ChatXP returns XP for a chat message
func (b GradeBand) ChatXP() int {
	switch b {
	case BandSprouts:
		return 1
	case BandExplorers:
		return 2
	case BandChampions:
		return 3
	default:
		return 2
	}
}

// AssignmentXP calculates XP for an assignment based on score
func (b GradeBand) AssignmentXP(score int) int {
	// Base XP: 1 per 10 points scored
	xp := score / 10
	if xp < 1 && score > 0 {
		xp = 1
	}

	// Bonus for high scores
	switch b {
	case BandSprouts:
		if score >= 90 {
			xp += 3
		} else if score >= 70 {
			xp += 2
		}
	case BandExplorers:
		if score >= 90 {
			xp += 5
		} else if score >= 70 {
			xp += 3
		} else if score >= 50 {
			xp += 1
		}
	case BandChampions:
		if score >= 90 {
			xp += 7
		} else if score >= 70 {
			xp += 4
		} else if score >= 50 {
			xp += 2
		}
	}

	return xp
}

// AssignmentPoints calculates points for an assignment
func (b GradeBand) AssignmentPoints(score int) int {
	// Points scale with score
	points := score / 20
	if points < 1 && score > 0 {
		points = 1
	}

	// Bonus for high performers
	if score >= 90 {
		points += 2
	} else if score >= 70 {
		points += 1
	}

	return points
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
