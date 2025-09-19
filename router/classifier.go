package router

import (
	"encoding/json"
)

type LLMClassifierRequest struct {
	Question string
}

type LLMClassifierResponse struct {
	Method string  `json:"method"` // "zero-shot" | "cot" | "few-shot" | "adihq"
	Mode   string  `json:"mode"`   // "" | "repo" | "dynamic"
	Reason string  `json:"reason"`
	Score  float64 `json:"score"` // 0..1
}

func BuildClassifierPrompt(q string) string {
	return "" +
		`You are a routing classifier for prompting strategies. 
Given a LeetCode-style problem, choose the BEST prompting method for solving it with a Golang solution.

Allowed methods:
- "zero-shot": direct concise solution, when the pattern is straightforward (two-pointers, sliding window).
- "cot": chain-of-thought reasoning, when multi-step logic or tricky math/DP counting is needed.
- "few-shot": when showing 1–2 analogous examples helps (classic patterns like parentheses, intervals, anagrams). 
  Prefer "mode":"repo" unless the repo is empty, then "dynamic".
- "adihq": structured hybrid (Analyze, Design, Implement, Handle, Quality) when edge cases/constraints matter.

Return a single-line strict JSON, no extra text.

Input:
` + q + `

Output JSON schema:
{"method":"zero-shot|cot|few-shot|adihq","mode":"","reason":"short","score":0.0}`
}

func ParseClassifierJSON(s string) (LLMClassifierResponse, error) {
	var r LLMClassifierResponse
	err := json.Unmarshal([]byte(firstJSONLine(s)), &r)
	if r.Mode == "" && r.Method == "few-shot" {
		r.Mode = "repo"
	}
	return r, err
}

// extract first {...} JSON object from response (defensive)
func firstJSONLine(s string) string {
	start := -1
	end := -1
	for i, c := range s {
		if c == '{' {
			start = i
			break
		}
	}
	for j := len(s) - 1; j >= 0; j-- {
		if s[j] == '}' {
			end = j
			break
		}
	}
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}
