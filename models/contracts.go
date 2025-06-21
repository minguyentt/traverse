package models

import "time"

// DTO's
type ContractMetaData struct {
	Contract
	ReviewCounts int `json:"review_counts"`
}

type Contract struct {
	ID         int64               `json:"id"`
	UserID     int64               `json:"user_id"`
	JobTitle   string              `json:"job_title"`
	City       string              `json:"city"`
	Agency     string              `json:"agency"`
	JobDetails *ContractJobDetails `json:"job_details"`
	Reviews    []Review            `json:"reviews"`
	User       User                `json:"user"`
	CreatedAt  time.Time           `json:"created_at"`
	UpdatedAt  time.Time           `json:"updated_at"`
	Version    int                 `json:"version"`
}

type ContractJobDetails struct {
	ContractID       int64  `json:"contract_id"`
	Profession       string `json:"profession"`
	AssignmentLength string `json:"assignment_length"`
	Experience       string `json:"experience"`
}

// payloads
type ContractPayload struct {
	JobTitle   string            `json:"job_title"   validate:"required,max=50"`
	City       string            `json:"city"        validate:"required,max=20"`
	Agency     string            `json:"agency"      validate:"required,max=50"`
	JobDetails *JobDetailPayload `json:"job_details" validate:"required"` // dive tag will validate and go into the nested struct
}

type JobDetailPayload struct {
	Profession       string `json:"profession"        validate:"required,max=30"`
	AssignmentLength string `json:"assignment_length" validate:"required,max=20"`
	Experience       string `json:"experience"        validate:"required,max=20"`
}

type UpdateContractPayload struct {
	JobTitle  *string                 `json:"job_title" validate:"omitempty,max=50"`
	City      *string                 `json:"city"        validate:"omitempty,max=20"`
	Agency    *string                 `json:"agency"      validate:"omitempty,max=50"`
	JobDetail *UpdateJobDetailPayload `json:"job_detail" validate:"omitempty,dive"`
}

type UpdateJobDetailPayload struct {
	Profession       *string `json:"profession"        validate:"omitempty,max=30"`
	AssignmentLength *string `json:"assignment_length" validate:"omitempty,max=20"`
	Experience       *string `json:"experience"        validate:"omitempty,max=20"`
}
