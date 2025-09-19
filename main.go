package main

import (
	"Prompt_Engineering_algorithm/router"
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"

	"Prompt_Engineering_algorithm/llm"
	"Prompt_Engineering_algorithm/prompts"
)

func main() {
	// Load .env if present (fallback to OS env if not)
	_ = godotenv.Load(".env")

	// Flags
	method := flag.String("method", "zero-shot", "prompting method: auto | zero-shot | cot | few-shot | adihq")
	routerMode := flag.String("router", "hybrid", "auto routing: heuristic | llm | hybrid")
	routerThresh := flag.Float64("router-threshold", 0.7, "min confidence to accept heuristic decision")
	model := flag.String("model", "gpt-4o-mini", "LLM model name")
	temperature := flag.Float64("temperature", 0.2, "sampling temperature for main solve (0.0 - 2.0)")
	exemplarTemp := flag.Float64("exemplar-temperature", 0.6, "temperature for exemplar generation (few-shot dynamic only)")
	fewshotMode := flag.String("fewshot", "repo", "few-shot mode when --method=few-shot: repo | dynamic")
	k := flag.Int("k", 2, "number of repository examples for few-shot repo mode")
	questionFlag := flag.String("question", "", "question text (if empty, read from stdin)")
	flag.Parse()

	// Read question
	var question string
	if q := strings.TrimSpace(*questionFlag); q != "" {
		question = q
	} else {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("❓ Enter your LeetCode-style problem: ")
		q, _ := reader.ReadString('\n')
		question = strings.TrimSpace(q)
	}

	// Echo config
	fmt.Printf("⚙️  Config -> method=%s | model=%s | temp=%.2f", *method, *model, *temperature)
	if *method == "few-shot" && *fewshotMode == "dynamic" {
		fmt.Printf(" | exemplar_temp=%.2f", *exemplarTemp)
	}
	if *method == "few-shot" && *fewshotMode == "repo" {
		fmt.Printf(" | k=%d", *k)
	}
	fmt.Println()

	var prompt string
	var totalElapsed time.Duration
	var totalPromptTokens, totalCompletionTokens, totalTokens int

	var picked router.Decision
	if *method == "auto" {
		// 1) heuristic
		h := router.HeuristicPick(question)
		picked = h

		// 2) if needed, ask LLM classifier
		needLLM := (*routerMode == "llm") || (*routerMode == "hybrid" && h.Confidence < *routerThresh)
		if needLLM {
			fmt.Println("🧭 Router: invoking LLM classifier...")
			req := router.BuildClassifierPrompt(question)
			respText, _, err := llm.CallLLM(req, *model, 0.0) // deterministic
			if err == nil {
				cls, perr := router.ParseClassifierJSON(respText)
				if perr == nil && (cls.Method == "zero-shot" || cls.Method == "cot" || cls.Method == "few-shot" || cls.Method == "adihq") {
					picked.Method = router.Method(cls.Method)
					if cls.Method == "few-shot" && (cls.Mode == "repo" || cls.Mode == "dynamic") {
						picked.Mode = cls.Mode
					}
					picked.Reason = "LLM: " + cls.Reason
					picked.Confidence = cls.Score
				}
			}
		}
		fmt.Printf("🔎 Auto-Selected -> method=%s", picked.Method)
		if picked.Method == router.MethodFewShot && picked.Mode != "" {
			fmt.Printf(" | mode=%s", picked.Mode)
		}
		fmt.Printf(" | reason=%s | conf=%.2f\n", picked.Reason, picked.Confidence)
	}
	finalMethod := *method
	if finalMethod == "auto" {
		finalMethod = string(picked.Method)
	}

	switch finalMethod {
	case "zero-shot":
		prompt = prompts.BuildZeroShotPrompt(question)

	case "cot":
		prompt = prompts.BuildCotPrompt(question)

	case "adihq":
		prompt = prompts.BuildAdihqPrompt(question)

	case "few-shot":
		switch *fewshotMode {
		case "repo":
			// Optionally adjust K used by repo builder (simple version uses internal k=2).
			// If you want k to be wired through, expose it inside prompts.BuildFewShotPrompt or add a new function.
			// For now we call the current builder (which uses an internal default k).
			prompt = prompts.BuildFewShotPrompt(question)

		case "dynamic":
			// --- Stage 1: Ask the model to synthesize ONE analogous exemplar (on-the-fly) ---
			exemplarReq := prompts.BuildFewShotExemplarRequest(question)
			fmt.Println("Generating exemplar (stage 1)...")
			stage1Start := time.Now()
			exemplar, usage1, err := llm.CallLLM(exemplarReq, *model, float32(*exemplarTemp))
			stage1Elapsed := time.Since(stage1Start)
			if err != nil {
				fmt.Printf("❌ Stage 1 error: %v\n", err)
				return
			}
			fmt.Println("✅ Exemplar generated.\n")

			// accumulate metrics
			totalElapsed += stage1Elapsed
			totalPromptTokens += usage1.PromptTokens
			totalCompletionTokens += usage1.CompletionTokens
			totalTokens += usage1.TotalTokens

			// --- Stage 2: Build final few-shot prompt by injecting the exemplar + solve target problem ---
			prompt = prompts.BuildFewShotDynamicPrompt(question, exemplar)

		default:
			fmt.Println("❌ unknown few-shot mode; use: repo | dynamic")
			return
		}

	default:
		fmt.Println("❌ unknown method; use: zero-shot | cot | few-shot | adihq")
		return
	}

	// --- Final solve call (for all methods, and stage 2 of dynamic few-shot) ---
	fmt.Println("🧠 Solving target problem (final call)...")
	start := time.Now()
	response, usage2, err := llm.CallLLM(prompt, *model, float32(*temperature))
	elapsed := time.Since(start)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	// aggregate metrics
	totalElapsed += elapsed
	totalPromptTokens += usage2.PromptTokens
	totalCompletionTokens += usage2.CompletionTokens
	totalTokens += usage2.TotalTokens

	// Output
	fmt.Println("\n🔹 Model Answer:")
	fmt.Println(response)

	fmt.Println("\n📊 Usage & Timing:")
	fmt.Printf("Model: %s | Method: %s", *model, *method)
	if *method == "few-shot" {
		fmt.Printf(" | Mode: %s", *fewshotMode)
	}
	fmt.Printf(" | Temp: %.2f\n", *temperature)
	fmt.Printf("Prompt tokens: %d\n", totalPromptTokens)
	fmt.Printf("Completion tokens: %d\n", totalCompletionTokens)
	fmt.Printf("Total tokens: %d\n", totalTokens)
	fmt.Printf("Total elapsed: %s\n", totalElapsed)
}
