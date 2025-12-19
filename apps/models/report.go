package models

type StatisticsResponse struct {
	TotalStudents     int                   `json:"total_students"`
	TotalLectures     int                   `json:"total_lectures"`
	TotalAchievements int                   `json:"total_achievements"`
	AchievementStats  AchievementStatistics `json:"achievement_stats"`
}

type AchievementStatistics struct {
	Draft     int `json:"draft"`
	Submitted int `json:"submitted"`
	Verified  int `json:"verified"`
	Rejected  int `json:"rejected"`
}

type StudentReportResponse struct {
	Student           Students               `json:"student"`
	TotalAchievements int                    `json:"total_achievements"`
	AchievementStats  AchievementStatistics  `json:"achievement_stats"`
	Achievements      []AchievementReference `json:"achievements"`
}
