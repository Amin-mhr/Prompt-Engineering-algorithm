package prompts

// BuildAhdiqPrompt builds an ADIHQ-style prompt tailored for LeetCode/algorithmic coding in Go.
// ADIHQ = Analyze, Design, Implement, Handle, Quality, Redundancy Check
// هدف: استدلال مختصر و قانون‌مند مثل CoT ولی کم‌هزینه‌تر (توکن کمتر) با تمرکز بر کدنویسی تمیز Go.
func BuildAdihqPrompt(question string) string {
	return "" +
		`You are an expert competitive programmer and senior Go engineer (Go 1.20+).
Answer LeetCode-style coding problems using the **ADIHQ** framework.
Return a concise reasoning and a single clean, correct **Go** solution. Avoid any duplicate functions or repeated code.

# TASK
` + question + `

# ADIHQ STEPS (keep each section brief)
- Analyze: Extract core operation, inputs, outputs, constraints, and important edge cases.
- Design: If multiple solutions exist, pick the most efficient (time/space). Mention the chosen approach (e.g., two-pointers, heap, DP, greedy, graph).
- Implement: Write idiomatic Go with clear names; avoid unnecessary packages; keep it production-quality and readable.
- Handle: Consider edge cases (empty input, bounds, overflow, duplicates, invalid states if applicable).
- Quality: Minimal but useful comments; deterministic behavior; no global state; prefer iterative unless recursion is natural.
- Redundancy Check: Ensure you output only **one** final function (and a minimal ` + "`main`" + ` for demo) without repeating the solution.

# OUTPUT FORMAT (strictly follow)
<REASONING>
- (3–8 bullets with brief, stepwise reasoning aligned to ADIHQ)
</REASONING>

<SOLUTION_GO>

	// brief description of the approach
	package main

	import "fmt"

	// implement solution here
	func solve(/* adapt signature to the problem */) /* return types */ {
		// ...
	}

	func main()
	{
		// optional tiny demo using one representative test
		fmt.Println( /* call solve(...) */)
	}
</SOLUTION_GO>`
}
