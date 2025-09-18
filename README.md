# Prompt Engineering Algorithms in Go

This project implements and compares multiple **prompt engineering strategies** for Large Language Models (LLMs), tailored to **LeetCode-style coding problems**.  
It allows users to send a question to an LLM with a chosen prompting method and returns both the model's response and usage statistics (tokens, latency, etc.).

---

## ✨ Supported Prompting Methods
- **Zero-Shot** – Direct question answering without examples.
- **CoT (Chain-of-Thought)** – Step-by-step reasoning before the final answer.
- **Few-Shot** – Adds a few worked examples before the actual question.
- **AHDIQ** – Advanced prompting technique (to be implemented according to paper definitions).

---

## 📂 Project Structure
# Prompt-Engineering-algorithm
/prompt-system
/llm
client.go # LLM client using OpenAI API
/prompts
zeroshot.go
zeroshot_cot.go
fewshot.go
ahdiq.go
main.go # CLI entry point
go.mod
README.md
.env.example # API key placeholder


---

## 🔧 Setup

1. **Clone repository**
   ```bash
   git clone https://github.com/yourname/prompt-system.git
   cd prompt-system

Install dependencies

go mod tidy


Set your OpenAI API key
Create a .env file in the project root:

OPENAI_API_KEY=sk-xxxxxxxxxxxxxxxx


Or set it directly in your shell:

export OPENAI_API_KEY="sk-xxxxxxxxxxxxxxxx"

▶️ Usage

Run with defaults (zero-shot, gpt-4o-mini, temperature=0.2, question from stdin):

go run main.go


Provide all flags in a single command:

go run main.go \
--method=cot \
--model=gpt-4o-mini \
--temperature=0.2 \
--question="Given an array of integers, return the two numbers that sum to a target."

⚙️ CLI Flags
Flag	Type	Default	Description
--method	string	zero-shot	Prompting method: zero-shot, cot, few-shot, ahdiq
--model	string	gpt-4o-mini	LLM model name
--temperature	float	0.2	Sampling temperature (0.0 = deterministic, 2.0 = very creative)
--question	string	""	Question text (if omitted, program will ask via stdin)
📊 Output Example
🔹 پاسخ مدل:
<the LLM-generated solution>

📊 اطلاعات پرامپت:
مدل: gpt-4o-mini | روش: cot | دما: 0.20
توکن ورودی: 128
توکن خروجی: 256
کل توکن: 384
زمان پاسخ: 1.23s
