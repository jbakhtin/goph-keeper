package query

const (
	CreateSecret = `
		INSERT INTO secrets (uuid, user_id, type, data, metadata)
		VALUES($1, $2, $3, $4, $5)
		RETURNING uuid, user_id, type, data, metadata, created_at, updated_at
	`
)
