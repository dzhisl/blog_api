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
	RoleOwner      Role = "owner"

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
	4: RoleOwner,
}

func IsValidRole(role Role) bool {
	for _, r := range RolesMap {
		if role == r {
			return true
		}
	}
	return false
}

// CompareRoles compares two roles based on their hierarchy defined in rolesMap.
// It returns:
//
//	-1 if role1 is lower than role2,
//	 0 if both roles are equal,
//	 1 if role1 is higher than role2.
//
// If either role is not found in rolesMap, it returns 0.
func CompareRoles(role1, role2 Role) int {
	var idx1, idx2 int = -1, -1

	for k, v := range RolesMap {
		if v == role1 {
			idx1 = k
		}
		if v == role2 {
			idx2 = k
		}
	}

	if idx1 == -1 || idx2 == -1 {
		// unknown role → consider them equal
		return 0
	}

	switch {
	case idx1 < idx2:
		return -1
	case idx1 > idx2:
		return 1
	default:
		return 0
	}
}
