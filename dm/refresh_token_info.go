package dm

type RefreshTokenInfo struct {
	Token       string `json:"-" db:"token"`
	UserID      int    `json:"-" db:"user_id"`
	IP          string `json:"ip" db:"ip_address"`
	IssuedAt    string `json:"issued-at" db:"issued_at"`
	ExpiresAt   string `json:"expires-at" db:"expires_at"`
	Invalidated bool   `json:"invalidated" db:"invalidated"`
}
