package models

type CreateSchedule struct {
	Description string `json:"description" binding:"required"`
	Daymonth    string `json:"day_month" binding:"required"`
	Priority    string `json:"prioity" binding:"required"`
}
