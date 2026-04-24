package types

type Template struct {
	Filename  string   `json:"filename"`
	Name      string   `yaml:"name" json:"name"`
	About     string   `yaml:"about" json:"about"`
	Title     string   `yaml:"title" json:"title"`
	Labels    []string `yaml:"labels" json:"labels"`
	Assignees []string `yaml:"assignees" json:"assignees"`
	Body      string   `json:"body"`
}

func (t Template) DisplayName() string {
	if t.Name != "" {
		return t.Name
	}
	return t.Filename
}
