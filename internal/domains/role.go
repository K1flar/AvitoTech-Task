package domains

type Role int

type RoleKey string

const (
	UserRole  Role = iota
	AdminRole Role = iota
)

func (r Role) String() string {
	return [...]string{"user", "admin"}[r]
}
