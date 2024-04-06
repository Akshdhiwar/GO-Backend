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
	SelectProductIdFromId         = "SELECT id FROM products WHERE id = $1"
	SelectAllFromID               = "SELECT id , created_at , updated_at , title , price , category , description , image , rating , count FROM products WHERE id=$1"
	SelectAllProductsLimit        = "SELECT * FROM products LIMIT $1 OFFSET $2"
	SelectAllProducts             = "SELECT * FROM products "
	SelectProductDetailsFromTitle = ` SELECT id , created_at , updated_at , title , price , category , description , image , rating , count  FROM products WHERE title = $1 LIMIT 1`
	SaveNewProduct                = ` INSERT INTO products ( title, price, category, image, description, rating , count) VALUES ($1, $2, $3, $4, $5 , $6 , $7)`
	DeleteProduct                 = ` DELETE FROM products WHERE id = $1`
	SelectIdFromProductsMismatch  = `SELECT id FROM products WHERE title = $1 AND id != $2`
	UpdateProduct                 = ` UPDATE products SET title = $1, price = $2, description = $3, category = $4, image = $5 , rating =$6, count =$7 WHERE id = $8`
)
