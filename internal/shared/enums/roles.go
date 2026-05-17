package enums

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleUser      Role = "user"
	RoleGuest     Role = "guest"
)

func (r Role) String() string { return string(r) }

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleModerator, RoleUser, RoleGuest:
		return true
	}
	return false
}
