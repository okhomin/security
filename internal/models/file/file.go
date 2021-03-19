package file

type File struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Content string   `json:"content"`
	Groups  []string `json:"groups"`
	Acls    []string `json:"acls"`
}
