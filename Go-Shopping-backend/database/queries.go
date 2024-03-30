package database

// users related queries for user apis
const (
	SelectIdFromEmail          = "SELECT id FROM users WHERE email=$1"
	SaveUserPassword           = `INSERT INTO users (email, password) VALUES ($1, $2)`
	SelectUserDetailsFromEmail = "SELECT id , password FROM users WHERE email=$1"
)
