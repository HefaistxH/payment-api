package config

const (
	CheckUserByEmailQuery      = "EXISTS(SELECT 1 FROM customers WHERE email = ?)"
	GetCredByEmailQuery        = "SELECT * FROM credentials WHERE email = ?"
	CheckCustomerBalanceQuery  = "SELECT balance FROM customers WHERE id = ?"
	CheckMerchantBalanceQuery  = "SELECT balance FROM merchants WHERE id = ?"
	UpdateCustomerBalanceQuery = "UPDATE customers SET balance = ? WHERE id = ?"
	UpdateMerchantBalanceQuery = "UPDATE merchants SET balance = ? WHERE id = ?"
	AddHistoryQuery            = "INSERT INTO history (customer_id, merchant_id, activity, amount, messages,time_stamp) VALUES (?, ?, ?, ?, ?, ?)"
)
