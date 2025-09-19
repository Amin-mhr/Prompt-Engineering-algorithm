package prompts

func BuildZeroShotPrompt(question string) string {
	return "" +
		`You are an expert Golang (Go 1.20+) engineer and competitive programmer.
Solve the following LeetCode-style problem and return a single, clean **Go** solution.

REQUIREMENTS
- Output must be Go code only, inside one fenced code block.
- Use idiomatic Go. No external packages.
- Implement a function named solve(/* adapt signature to problem */) and, optionally, a tiny main() demo.
- Handle edge cases safely. Keep comments minimal but useful.
- Do NOT add any text outside the code block. No explanations.

PROBLEM
` + question + `

OUTPUT (code block only):
	package main

	import "fmt"

	// solve implements the solution to the target problem.
	func solve(/* adapt signature */) /* return types */ {
		// ...
	}

	func main() {
		// optional: tiny demo
		fmt.Println(/* call solve(...) */)
	}`
}
