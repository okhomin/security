package group

type Group struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Read  bool     `json:"read"`
	Write bool     `json:"write"`
	Users []string `json:"users"`
}
