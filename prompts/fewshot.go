package prompts

func BuildFewShotPrompt(question string) string {
	examples := `Q: What is 2 + 2?
Explain: 2 plus 2 equals 4.
A: 4

Q: Who wrote 'Pride and Prejudice'?
Explain: The author is Jane Austen.
A: Jane Austen
`
	return examples + "\nQ: " + question + "\nExplain: Let's think step by step.\nA:"
}
