package prompts

import (
	"fmt"
	"strings"
)

// این تابع پرامپت مرحله‌ی ۱ را می‌سازد (گرفتن مثال از خود مدل)
func BuildFewShotExemplarRequest(question string) string {
	return "" +
		`You are a coding interview coach for LeetCode-style problems in Go.

I will give you a TARGET PROBLEM. Your job is to create ONE ANALOGOUS EXAMPLE
that teaches the same underlying pattern (data structure / algorithm), but with a different story,
different variable names, and different constraints if needed.

STRICT RULES:
- Do NOT solve or restate the target problem. Create a NEW, DIFFERENT example.
- Keep it SHORT and TOKEN-EFFICIENT.
- Provide: (1) Problem (2) 3–5 bullets of Reasoning (3) A short Go code skeleton (10–20 lines) (4) 2–3 small tests.
- The code should be idiomatic Go 1.20+, minimal comments, no external packages.

TARGET PROBLEM:
` + question + `

OUTPUT FORMAT:
### EXAMPLE
Problem: <one-sentence problem statement>
<REASONING>
- bullet 1
- bullet 2
- bullet 3
</REASONING>
<SOLUTION_GO>
" + // short skeleton showing the pattern
package main
func solve(/* signature */) /* returns */ {
    // ...
}
` + "```\n" + `</SOLUTION_GO>
<TESTS>
- Input: ... -> Output: ...
- Input: ... -> Output: ...
</TESTS>`
}

// این تابع پرامپت مرحله‌ی ۲ را می‌سازد (few-shot نهایی با مثال تولیدی)
func BuildFewShotDynamicPrompt(question, exemplar string) string {
	var b strings.Builder
	b.WriteString("You are an expert competitive programmer and Go engineer (Go 1.20+).\n")
	b.WriteString("Use the EXAMPLE as guidance for style/pattern only. Do NOT copy variable names blindly.\n\n")
	b.WriteString("### EXAMPLE (for guidance only)\n")
	b.WriteString(exemplar + "\n\n")
	b.WriteString("### TARGET PROBLEM\n")
	b.WriteString(question + "\n\n")
	b.WriteString(`# OUTPUT FORMAT
<REASONING>
- (3–8 concise bullets)
</REASONING>

<SOLUTION_GO>
` + "```go\n" + `package main
// clean, final solution to the TARGET PROBLEM
func solve(/* signature */) /* returns */ {
    // ...
}
` + "```\n" + `</SOLUTION_GO>

<COMPLEXITY>
Time: O(...), Space: O(...)
</COMPLEXITY>

<TESTS>
1) Input: ... -> Output: ...
2) Input: ... -> Output: ...
</TESTS>
`)
	return b.String()
}

func BuildFewShotPrompt(question string) string {
	repo := SeedExamples()
	k := 2 // می‌تونی از فلگ/کانفیگ بخونیش
	exs := PickTopK(question, repo, k)

	var sb strings.Builder
	sb.WriteString("You are an expert competitive programmer and Go engineer (Go 1.20+).\n")
	sb.WriteString("Solve the target problem using short, structured reasoning, then output a clean Go solution.\n\n")

	// Few-shot exemplars
	for i, ex := range exs {
		sb.WriteString(fmt.Sprintf("### EXAMPLE %d: %s\n", i+1, ex.Title))
		sb.WriteString("Problem:\n")
		sb.WriteString(ex.Description + "\n")
		if len(ex.Reasoning) > 0 {
			sb.WriteString("<REASONING>\n")
			for _, r := range ex.Reasoning {
				sb.WriteString("- " + r + "\n")
			}
			sb.WriteString("</REASONING>\n")
		}
		if ex.CodeGo != "" {
			sb.WriteString("<SOLUTION_GO>\n```go\n")
			sb.WriteString(ex.CodeGo + "\n")
			sb.WriteString("```\n</SOLUTION_GO>\n")
		}
		if len(ex.Tests) > 0 {
			sb.WriteString("<TESTS>\n")
			for _, t := range ex.Tests {
				sb.WriteString("- " + t + "\n")
			}
			sb.WriteString("</TESTS>\n")
		}
		sb.WriteString("\n---\n\n")
	}

	// Target question
	sb.WriteString("### TARGET PROBLEM\n")
	sb.WriteString(question + "\n\n")

	// خروجی استاندارد برای پارس
	sb.WriteString(`# OUTPUT FORMAT
<REASONING>
- (3–8 concise bullets)
</REASONING>

<SOLUTION_GO>

	// brief description
	package main

	import "fmt"

	func solve(/* adapt signature */) /* return types */ {
		// ...
	}

	func main() { /* tiny demo if useful */ }`)
	return sb.String()
}
