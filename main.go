package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"Prompt_Engineering_algorithm/llm"
	"Prompt_Engineering_algorithm/prompts"
)

func main() {
	// Flags
	method := flag.String("method", "zero-shot", "prompting method: zero-shot | zero-shot-cot | few-shot | ahdiq")
	model := flag.String("model", "gpt-4o-mini", "LLM model name")
	temperature := flag.Float64("temperature", 0.2, "sampling temperature (0.0 - 2.0)")
	questionFlag := flag.String("question", "", "question text (if empty, read from stdin)")
	flag.Parse()

	// Read question
	var question string
	if q := strings.TrimSpace(*questionFlag); q != "" {
		question = q
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("❓ سوال خود را وارد کنید: ")
		q, _ := reader.ReadString('\n')
		question = strings.TrimSpace(q)
	}

	// Build prompt
	var prompt string
	switch *method {
	case "zero-shot":
		prompt = prompts.BuildZeroShotPrompt(question)
	case "zero-shot-cot":
		prompt = prompts.BuildZeroShotCotPrompt(question)
	case "few-shot":
		prompt = prompts.BuildFewShotPrompt(question)
	case "ahdiq":
		prompt = prompts.BuildAhdiqPrompt(question)
	default:
		fmt.Println("روش ناشناخته! از بین: zero-shot | zero-shot-cot | few-shot | ahdiq")
		return
	}

	// Call LLM
	start := time.Now()
	response, usage, err := llm.CallLLM(prompt, *model, float32(*temperature))
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("❌ خطا: %v\n", err)
		return
	}

	// Output
	fmt.Println("\n🔹 پاسخ مدل:")
	fmt.Println(response)

	fmt.Println("\n📊 اطلاعات پرامپت:")
	fmt.Printf("مدل: %s | روش: %s | دما: %.2f\n", *model, *method, *temperature)
	fmt.Printf("توکن ورودی: %d\n", usage.PromptTokens)
	fmt.Printf("توکن خروجی: %d\n", usage.CompletionTokens)
	fmt.Printf("کل توکن: %d\n", usage.TotalTokens)
	fmt.Printf("زمان پاسخ: %s\n", elapsed)
}
