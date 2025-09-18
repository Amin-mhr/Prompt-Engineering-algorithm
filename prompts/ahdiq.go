package prompts

func BuildAhdiqPrompt(question string) string {
	// TODO: جایگزین با دستورالعمل دقیق AHDIQ طبق مقاله/طراحی شما
	return "Apply the AHDIQ strategy to answer the question.\nQuestion: " + question + "\nReason step by step, then provide the final answer."
}
