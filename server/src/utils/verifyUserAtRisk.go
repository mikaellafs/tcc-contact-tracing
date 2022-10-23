package utils

import "time"

func VerifyUserAtRisk(dateContact, dateDiagnostic time.Time, distance float32, maxTimeDiffToConsiderAtRisk time.Duration) bool {
	diff := time.Time.Sub(dateContact, dateDiagnostic)
	if diff < 0 {
		diff = -diff
	}

	return diff < maxTimeDiffToConsiderAtRisk && distance <= 200
}
