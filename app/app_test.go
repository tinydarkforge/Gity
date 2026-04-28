package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func rootWith(screen Screen) Root {
	r := New()
	m, _ := r.Update(NavigateMsg{Screen: screen})
	return m.(Root)
}

func update(r Root, msg tea.Msg) Root {
	m, _ := r.Update(msg)
	return m.(Root)
}

func TestNavigateToCreate(t *testing.T) {
	r := update(New(), NavigateMsg{Screen: ScreenCreate})
	if r.ActiveScreen() != ScreenCreate {
		t.Fatalf("want ScreenCreate, got %v", r.ActiveScreen())
	}
}

func TestNavigateToDetail(t *testing.T) {
	r := update(New(), NavigateMsg{Screen: ScreenDetail})
	if r.ActiveScreen() != ScreenDetail {
		t.Fatalf("want ScreenDetail, got %v", r.ActiveScreen())
	}
}

func TestNavigateToSettings(t *testing.T) {
	r := update(New(), NavigateMsg{Screen: ScreenSettings})
	if r.ActiveScreen() != ScreenSettings {
		t.Fatalf("want ScreenSettings, got %v", r.ActiveScreen())
	}
}

func TestBackMsgReturnsToNC(t *testing.T) {
	r := rootWith(ScreenCreate)
	r = update(r, BackMsg{})
	if r.ActiveScreen() != ScreenNC {
		t.Fatalf("want ScreenNC after BackMsg, got %v", r.ActiveScreen())
	}
}

func TestOpenIssueMsgNavigatesToDetail(t *testing.T) {
	r := update(New(), OpenIssueMsg{Number: 42})
	if r.ActiveScreen() != ScreenDetail {
		t.Fatalf("want ScreenDetail after OpenIssueMsg, got %v", r.ActiveScreen())
	}
}

func TestNewStartsAtNC(t *testing.T) {
	r := New()
	if r.ActiveScreen() != ScreenNC {
		t.Fatalf("want ScreenNC on init, got %v", r.ActiveScreen())
	}
}
