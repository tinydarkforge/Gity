package types

import "time"

type Issue struct {
	Number    int       `json:"number"`
	Title     string    `json:"title"`
	State     string    `json:"state"`
	URL       string    `json:"url"`
	Body      string    `json:"body"`
	Labels    []Label   `json:"labels"`
	Assignees []User    `json:"assignees"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Label struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type User struct {
	Login string `json:"login"`
}

type Comment struct {
	Author    User      `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"createdAt"`
}

type Draft struct {
	Title    string   `json:"title"`
	Body     string   `json:"body"`
	Labels   []string `json:"labels"`
	Template string   `json:"template"`
}
