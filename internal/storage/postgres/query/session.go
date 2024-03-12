package query

const (
	CreateSession = `
		INSERT INTO sessions (user_id, refresh_token, finger_print, expire_at)
		VALUES($1, md5(random()::text), $2, $3)
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	UpdateSessionByID = `
		UPDATE sessions
		SET user_id = $2, refresh_token = md5(random()::text), finger_print = $3, expire_at = $4, closed_at = $5
		WHERE sessions.id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	CloseSessionByID = `
		UPDATE sessions
		SET closed_at = NOW()
		WHERE sessions.id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	CloseSessionsByUserID = `
		UPDATE sessions
		SET closed_at = NOW()
		WHERE sessions.user_id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	UpdateSessionRefreshTokenByID = `
		UPDATE sessions
		SET refresh_token = md5(random()::text)
		WHERE sessions.id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	GetSessionByID = `
		SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at FROM sessions
		WHERE sessions.id = $1 AND sessions.closed_at is NULL LIMIT 1
	`

	SearchSessionsTemp = `SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at FROM sessions`

	GetSessionByUserIDAndFingerPrint = `
		SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at FROM sessions
		WHERE sessions.user_id = $1 AND sessions.finger_print = $2 AND sessions.closed_at is NULL LIMIT 1
	`
	GetSessionByUserID = `
		SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at FROM sessions
		WHERE sessions.user_id = $1 AND sessions.closed_at is NULL
	`

	GetSessionByRefreshToken = `
		SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at 
		FROM sessions 
		WHERE refresh_token = $1
	`
)
