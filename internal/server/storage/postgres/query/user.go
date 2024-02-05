package query

const (
	CreateUser = `
		INSERT INTO users (email, password, created_at)
		VALUES($1, $2, NOW())
		RETURNING id, email, password, created_at, updated_at
	`

	UpdateUser = `
		INSERT INTO users (email, password, updated_at)
		VALUES($1, $2, NOW())
		RETURNING id, email, password, created_at, updated_at
	`

	GetUserByID = `
		SELECT id, email, password, created_at, updated_at FROM users
		WHERE users.id = $1
	`

	GetUserByEmail = `
		SELECT id, email, password, created_at, updated_at FROM users
		WHERE users.email = $1
	`

	GetUsers = `
		SELECT id, email, created_at, updated_at FROM users
	`

	SearchUserTemp = `SELECT id, email, password, created_at, updated_at FROM users`
)
