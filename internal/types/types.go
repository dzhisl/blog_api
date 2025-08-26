package types

var RequestIDKey = "request_id"

type Role string
type Status string
type TokenType string

var (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"

	StatusOk Status = "ok"

	TokenAccess  TokenType = "access"
	TokenRefresh TokenType = "refresh"
)
