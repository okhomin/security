package acl

type Acl struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
	Read   bool   `json:"read"`
	Write  bool   `json:"write"`
}
