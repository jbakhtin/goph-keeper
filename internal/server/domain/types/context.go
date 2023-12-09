package types

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

var (
	ContextKeyUserID    = contextKey("user_id")
	ContextKeySessionID = contextKey("session_id")
)
