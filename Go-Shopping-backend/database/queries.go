package database

const (
	SelectUserIdFromEmail      = "SELECT id FROM users WHERE email=$1"
	SaveUserPassword           = `INSERT INTO users (email, password) VALUES ($1, $2)`
	SelectUserDetailsFromEmail = "SELECT id , password FROM users WHERE email=$1"
	SelectUserIdFromID         = "SELECT id FROM users WHERE id=$1"
	UpdateCartId               = "UPDATE users SET cart_id=$1 WHERE id=$2"
	SelectCartIdFromId         = "SELECT cart_id FROM users WHERE id = $1"
)

// cart related queries
const (
	SelectCartIdFromUserId       = "SELECT id FROM carts WHERE user_id=$1"
	SaveCart                     = "INSERT INTO carts (user_id, products, created_at, updated_at) VALUES ($1, $2 , $3 , $4)"
	SelectCartDetailsFromUserId  = "SELECT products , user_id FROM carts WHERE user_id = $1"
	UpdateCart                   = `UPDATE carts SET products = $1, updated_at = $2 WHERE user_id = $3`
	SelectProductsFromUserID     = "SELECT products FROM carts WHERE user_id = $1"
	UpdateCartProductWhereUserId = "UPDATE carts SET products = $1 WHERE user_id = $2"
)

// product related queries
const (
	SelectProductIdFromId = "SELECT id FROM products WHERE id = $1"
	SelectAllFromID       = "SELECT * FROM products WHERE id=$1"
)
