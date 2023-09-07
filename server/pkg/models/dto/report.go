package dto

import "time"

type Report struct {
	ID             int64
	DateDiagnostic time.Time
}

type GrpcReport struct {
	UserId            string `json:"userId,omitempty"`
	DateStartSymptoms int64  `json:"dateStartSymptoms,omitempty"`
	DateDiagnostic    int64  `json:"dateDiagnostic,omitempty"`
	DateReport        int64  `json:"dateReport"`
}
