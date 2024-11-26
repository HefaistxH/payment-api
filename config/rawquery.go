package config

const (
	CheckUserByEmailQuery      = "SELECT EXISTS (SELECT 1 FROM credentials WHERE email = $1)"
	GetCredByEmailQuery        = "SELECT id, user_id, email, password, role FROM credentials WHERE email = $1"
	CheckCustomerBalanceQuery  = "SELECT balance FROM customers WHERE id = $1"
	CheckMerchantBalanceQuery  = "SELECT balance FROM merchants WHERE id = $1"
	UpdateCustomerBalanceQuery = "UPDATE customers SET balance = $1 WHERE id = $2"
	UpdateMerchantBalanceQuery = "UPDATE merchants SET balance = $1 WHERE id = $2"
	AddHistoryQuery            = "INSERT INTO history (customer_id, merchant_id, activity, amount, message, timestamp) VALUES ($1, $2, $3, $4, $5, $6)"
)
