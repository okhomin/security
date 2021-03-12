package file

type File struct {
	ID      string
	Name    string
	Mode    uint16
	Content string
	Groups  []string
}
