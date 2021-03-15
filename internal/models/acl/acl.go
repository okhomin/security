package acl

type Acl struct {
	ID     string
	UserID string
	Read   bool
	Write  bool
}
