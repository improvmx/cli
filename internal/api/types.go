package api

type Domain struct {
	Display       string   `json:"display"`
	Name          string   `json:"domain"`
	Active        bool     `json:"active"`
	NotificationEmail string `json:"notification_email"`
	Whitelabel    string   `json:"whitelabel"`
	Added         int64    `json:"added"`
	Aliases       []Alias  `json:"aliases,omitempty"`
}

type Alias struct {
	Alias   string `json:"alias"`
	Forward string `json:"forward"`
	ID      int    `json:"id"`
}

type LogEntry struct {
	ID        string            `json:"id"`
	Created   int64             `json:"created"`
	CreatedAt string            `json:"created_at"`
	Events    []LogEvent        `json:"events"`
	Forward   LogForward        `json:"forward"`
	Hostname  string            `json:"hostname"`
	MessageID string            `json:"messageId"`
	Recipient LogRecipient      `json:"recipient"`
	Sender    LogSender         `json:"sender"`
	Subject   string            `json:"subject"`
	Transport string            `json:"transport"`
}

type LogEvent struct {
	Code    int    `json:"code"`
	Created int64  `json:"created"`
	Local   string `json:"local"`
	Message string `json:"message"`
	Server  string `json:"server"`
	Status  string `json:"status"`
}

type LogForward struct {
	Email string `json:"email"`
}

type LogRecipient struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LogSender struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type SMTPCredential struct {
	Created  int64  `json:"created"`
	Usage    int    `json:"usage"`
	Username string `json:"username"`
}

type Account struct {
	BillingEmail string `json:"billing_email"`
	CancelsOn    int64  `json:"cancels_on"`
	CompanyName  string `json:"company_name"`
	CompanyVAT   string `json:"company_vat"`
	Country      string `json:"country"`
	Email        string `json:"email"`
	Last4        string `json:"last4"`
	Limits       AccountLimits `json:"limits"`
	Lock         bool   `json:"lock"`
	Password     bool   `json:"password"`
	Plan         AccountPlan `json:"plan"`
	Premium      bool   `json:"premium"`
	Privacy      bool   `json:"privacy"`
	RenewDate    int64  `json:"renew_date"`
}

type AccountLimits struct {
	Aliases    int `json:"aliases"`
	DailyQuota int `json:"daily_quota"`
	Domains    int `json:"domains"`
	Ratelimit  int `json:"ratelimit"`
}

type AccountPlan struct {
	Display string `json:"display"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
}

type Rule struct {
	ID      string                 `json:"id"`
	Type    string                 `json:"type"`
	Config  map[string]interface{} `json:"config"`
	Rank    float64                `json:"rank"`
	Active  bool                   `json:"active"`
	Created int64                  `json:"created"`
}
