package jsonx

import (
	"encoding/json"
	"errors"
	"strings"
)

// Extract finds a JSON object inside arbitrary model output. It tolerates
// fenced code blocks, leading/trailing prose, and picks the largest
// balanced {...} region when multiple candidates exist.
func Extract(raw string, v any) error {
	s := strings.TrimSpace(raw)

	if err := json.Unmarshal([]byte(s), v); err == nil {
		return nil
	}

	if stripped, ok := stripFence(s); ok {
		if err := json.Unmarshal([]byte(stripped), v); err == nil {
			return nil
		}
		s = stripped
	}

	if obj, ok := largestObject(s); ok {
		if err := json.Unmarshal([]byte(obj), v); err == nil {
			return nil
		}
	}

	return errors.New("no parseable JSON object found")
}

func stripFence(s string) (string, bool) {
	if !strings.HasPrefix(s, "```") {
		return s, false
	}
	nl := strings.IndexByte(s, '\n')
	if nl < 0 {
		return s, false
	}
	end := strings.LastIndex(s, "```")
	if end <= nl {
		return s, false
	}
	return strings.TrimSpace(s[nl+1 : end]), true
}

func largestObject(s string) (string, bool) {
	best, bestLen := "", 0
	for i := 0; i < len(s); i++ {
		if s[i] != '{' {
			continue
		}
		depth := 0
		inStr := false
		esc := false
		for j := i; j < len(s); j++ {
			c := s[j]
			if esc {
				esc = false
				continue
			}
			if c == '\\' && inStr {
				esc = true
				continue
			}
			if c == '"' {
				inStr = !inStr
				continue
			}
			if inStr {
				continue
			}
			switch c {
			case '{':
				depth++
			case '}':
				depth--
				if depth == 0 {
					cand := s[i : j+1]
					if len(cand) > bestLen {
						best, bestLen = cand, len(cand)
					}
					break
				}
			}
			if depth == 0 && c == '}' {
				break
			}
		}
	}
	return best, bestLen > 0
}
