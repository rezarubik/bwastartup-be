package transaction

import "time"

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string // can be "pending", "success", "failed"
	Code       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
