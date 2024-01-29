package types

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// ToDo: need move to infrastructure
var (
	ContextKeyUserID    = contextKey("user_id")
	ContextKeySessionID = contextKey("session_id")
)
