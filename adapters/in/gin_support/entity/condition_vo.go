package entity

import "time"

type ConditionRequest struct {
	Name            string    `json:"name" binding:"omitempty,min=1,max=500"`
	InsertTimeStart time.Time `json:"insert_time_start" binding:"omitempty"`
	InsertTimeEnd   time.Time `json:"insert_time_end" binding:"omitempty"`
	Description     string    `json:"description" binding:"omitempty,min=1,max=1024"`
	Series          string    `json:"series" binding:"omitempty,min=1,max=500"`
	IsPlay          *bool     `json:"is_play" binding:"omitempty"`
	CategoryId      string    `json:"category_id" binding:"omitempty,min=1,max=64"`
}
