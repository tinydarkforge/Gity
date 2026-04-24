package types

type AgentStatus string

const (
	StatusNeedsInfo AgentStatus = "needs_info"
	StatusReady     AgentStatus = "ready"
)

type AgentResponse struct {
	Status    AgentStatus `json:"status"`
	Questions []string    `json:"questions,omitempty"`
	Title     string      `json:"title,omitempty"`
	Body      string      `json:"body,omitempty"`
	Rationale string      `json:"rationale,omitempty"`
}

func (r AgentResponse) IsReady() bool {
	if r.Status == StatusReady {
		return true
	}
	return r.Status == "" && r.Title != "" && r.Body != ""
}

type Turn struct {
	Question string
	Answer   string
}
