package services_test

import (
	"testing"

	"github.com/tinydarkforge/gity/types"
)

// agentResponse helpers — unit-test the AgentResponse.IsReady logic without
// hitting Ollama.

func TestAgentResponseIsReady(t *testing.T) {
	cases := []struct {
		name  string
		resp  types.AgentResponse
		ready bool
	}{
		{
			name:  "explicit ready",
			resp:  types.AgentResponse{Status: types.StatusReady, Title: "T", Body: "B"},
			ready: true,
		},
		{
			name:  "needs_info",
			resp:  types.AgentResponse{Status: types.StatusNeedsInfo, Questions: []string{"q1"}},
			ready: false,
		},
		{
			name:  "implicit ready — status omitted but title+body present",
			resp:  types.AgentResponse{Title: "T", Body: "B"},
			ready: true,
		},
		{
			name:  "empty response — not ready",
			resp:  types.AgentResponse{},
			ready: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.resp.IsReady(); got != tc.ready {
				t.Errorf("IsReady() = %v, want %v", got, tc.ready)
			}
		})
	}
}
