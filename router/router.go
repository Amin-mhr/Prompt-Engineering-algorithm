package router

import (
	"regexp"
	"strings"
)

type Method string

const (
	MethodZeroShot Method = "zero-shot"
	MethodCoT      Method = "cot"
	MethodFewShot  Method = "few-shot" // mode decided by caller: repo|dynamic
	MethodADIHQ    Method = "adihq"
)

type Decision struct {
	Method     Method
	Mode       string // only for few-shot: "repo" or "dynamic"
	Reason     string
	Confidence float64 // 0..1 (heuristic-based)
}

// quick lowercase helper
func lc(s string) string { return strings.ToLower(s) }

// very light heuristics for LeetCode-style prompts
func HeuristicPick(question string) Decision {
	q := lc(question)

	// Signals for CoT: "explain", "prove", multi-step math, "minimum steps", "derive", "proof"
	if reAny(q, `explain|prove|why|derive|show that|minimum steps|number of ways|count paths|dp|dynamic programming|probability`) {
		return Decision{Method: MethodCoT, Reason: "reasoning-heavy / DP/counting keywords", Confidence: 0.8}
	}

	// Signals for Few-Shot (benefits from pattern anchoring): parentheses, intervals merge, anagram variants, parsing formats
	if reAny(q, `valid parentheses|parentheses|merge intervals|intervals|anagram|decode|parse|format`) {
		return Decision{Method: MethodFewShot, Mode: "repo", Reason: "classic patterns benefit from exemplars", Confidence: 0.7}
	}

	// Graph words -> CoT or ADIHQ depending on phrasing
	if reAny(q, `graph|bfs|dfs|shortest path|dijkstra|topo|topological|cycle`) {
		return Decision{Method: MethodADIHQ, Reason: "algorithmic with edge-cases; ADIHQ balances guidance and tokens", Confidence: 0.65}
	}

	// Sliding window / two-pointers often solvable succinctly
	if reAny(q, `sliding window|two pointers|two-pointers|window|substring without repeating`) {
		return Decision{Method: MethodZeroShot, Reason: "succinct optimal pattern", Confidence: 0.65}
	}

	// If constraints/edge-cases emphasized -> ADIHQ
	if reAny(q, `constraints|edge cases|overflow|large input|time limit|memory limit`) {
		return Decision{Method: MethodADIHQ, Reason: "quality & edge-case handling requested", Confidence: 0.7}
	}

	// Fallback defaults
	return Decision{Method: MethodADIHQ, Reason: "default safe choice", Confidence: 0.5}
}

func reAny(s string, pat string) bool {
	re := regexp.MustCompile(pat)
	return re.FindStringIndex(s) != nil
}
