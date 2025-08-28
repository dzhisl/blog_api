package types

var RequestIDKey = "request_id"

type Role string
type Status string
type TokenType string

var (
	RoleUser       Role = "user"
	RoleModerator  Role = "moderator"
	RoleAdmin      Role = "admin"
	RoleSuperAdmin Role = "super_admin"

	StatusOk     Status = "ok"
	StatusBanned Status = "banned"

	TokenAccess  TokenType = "access"
	TokenRefresh TokenType = "refresh"
)

var RolesMap = map[int]Role{
	0: RoleUser,
	1: RoleModerator,
	2: RoleAdmin,
	3: RoleSuperAdmin,
}
