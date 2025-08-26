package query

const (
	CreateSecret = `
		INSERT INTO secrets (user_id, description, data)
		VALUES($1, $2, $3)
		RETURNING id, user_id, description, data, created_at, updated_at
	`
)
