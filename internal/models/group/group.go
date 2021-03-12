package group

type Group struct {
	ID    string
	Name  string
	Read  bool
	Write bool
	Users []string
}
