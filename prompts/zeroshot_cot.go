package prompts

func BuildZeroShotCotPrompt(question string) string {
	// CoT: تشویق به استدلال گام‌به‌گام
	return "You are a careful reasoner. Solve step by step, then give the final answer.\n" +
		"Question: " + question + "\n" +
		"Let's think step by step.\n" +
		"Finally, provide a short final answer after the reasoning."
}
