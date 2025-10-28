package models

// UserPreferences defines the structure for user preferences stored in JSON
type UserPreferences struct {
	// Email preferences
	MarketingEmails    bool `json:"marketing_emails"`
	OrderNotifications bool `json:"order_notifications"`
	BlogNotifications  bool `json:"blog_notifications"`

	// Display preferences
	Language string `json:"language"` // en, es, fr, etc.
	Timezone string `json:"timezone"` // UTC, America/New_York, etc.
	Theme    string `json:"theme"`    // light, dark, auto

	// Privacy settings
	ProfileVisibility   string `json:"profile_visibility"` // public, private
	ShowEmail           bool   `json:"show_email"`
	ShowPurchaseHistory bool   `json:"show_purchase_history"`
}

// DefaultPreferences returns default user preferences
func DefaultPreferences() UserPreferences {
	return UserPreferences{
		MarketingEmails:     true,
		OrderNotifications:  true,
		BlogNotifications:   true,
		Language:            "en",
		Timezone:            "UTC",
		Theme:               "light",
		ProfileVisibility:   "public",
		ShowEmail:           false,
		ShowPurchaseHistory: false,
	}
}
