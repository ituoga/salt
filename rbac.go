package salt

type RBACContext interface {
	WithToken(string) RBACContext
	WithResource(string) RBACContext
	Can(action string) bool
}
