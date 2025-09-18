package prompts

func BuildZeroShotPrompt(question string) string {
	// بدون راهنمایی زنجیره‌تفکر؛ فقط سؤال، یا با نقش/قیود کوچک
	return "Answer the following question concisely and correctly:\n\nQ: " + question + "\nA:"
}
