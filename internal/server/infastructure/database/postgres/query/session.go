package query

const (
	CreateSession = `
		INSERT INTO sessions (user_id, refresh_token, finger_print, expire_at, created_at)
		VALUES($1, md5(random()::text), $2, $3, NOW())
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	UpdateSessionByID = `
		UPDATE sessions
		SET refresh_token = md5(random()::text), expire_at = $2, updated_at = NOW()
		WHERE sessions.id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	CloseSessionByID = `
		UPDATE sessions
		SET closed_at = NOW(), updated_at = NOW()
		WHERE sessions.id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	CloseSessionsByUserID = `
		UPDATE sessions
		SET closed_at = NOW(), updated_at = NOW()
		WHERE sessions.user_id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	UpdateSessionRefreshTokenById = `
		UPDATE sessions
		SET updated_at = NOW()
		WHERE sessions.id = $1
		RETURNING id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at
	`

	GetSessionByUserIDAndFingerPrint = `
		SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at FROM sessions
		WHERE sessions.user_id = $1 AND sessions.finger_print = $2 AND sessions.closed_at is NULL LIMIT 1
	`

	GetSessionByRefreshToken = `
		SELECT id, user_id, refresh_token, finger_print, expire_at, created_at, closed_at, updated_at 
		FROM sessions 
		WHERE refresh_token = $1
	`
)
