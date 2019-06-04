package profiles

// ChurnProfile represents a single churn model record for a customer
type ChurnProfile struct {
	CustomerID       string
	Partner          string
	Dependents       string
	Tenure           int
	PhoneService     string
	MultipleLines    string
	InternetService  string
	OnlineSecurity   string
	OnlineBackup     string
	DeviceProtection string
	TechSupport      string
	StreamingTV      string
	StreamingMovies  string
	Contract         string
	PaperlessBilling string
	PaymentMethod    string
	MonthlyCharges   float64
	TotalCharges     float64
}