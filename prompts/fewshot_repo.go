package prompts

type Example struct {
	ID          string
	Title       string
	Tags        []string
	Description string   // توضیح مسأله (خیلی کوتاه)
	Reasoning   []string // 2-5 بولت خلاصه
	CodeGo      string   // راه‌حل Go کوتاه
	Tests       []string // چند نمونه ورودی/خروجی
}

func SeedExamples() []Example {
	return []Example{
		{
			ID:          "two-sum",
			Title:       "Two Sum",
			Tags:        []string{"array", "hashmap"},
			Description: "Given nums and target, return indices of two numbers such that they add up to target.",
			Reasoning: []string{
				"use hashmap to store value->index",
				"for each x, check target-x in map",
			},
			CodeGo: `package main
func twoSum(nums []int, target int) []int {
    m := make(map[int]int)
    for i, x := range nums {
        if j, ok := m[target-x]; ok { return []int{j, i} }
        m[x] = i
    }
    return nil
}`,
			Tests: []string{
				"nums=[2,7,11,15], target=9 -> [0,1]",
				"nums=[3,2,4], target=6 -> [1,2]",
			},
		},
		{
			ID:          "valid-parentheses",
			Title:       "Valid Parentheses",
			Tags:        []string{"string", "stack"},
			Description: "Check if parentheses/brackets are valid.",
			Reasoning: []string{
				"use stack; push opens",
				"on close, top must match",
			},
			CodeGo: `package main
func isValid(s string) bool {
    st := []rune{}
    pair := map[rune]rune{')':'(',']':'[','}':'{'}
    for _, c := range s {
        if c=='(' || c=='[' || c=='{' { st = append(st, c); continue }
        if p, ok := pair[c]; !ok || len(st)==0 || st[len(st)-1]!=p { return false }
        st = st[:len(st)-1]
    }
    return len(st)==0
}`,
			Tests: []string{
				"() -> true", "()[]{} -> true", "(] -> false",
			},
		},
		// ... مثال‌های بیشتر با تنوع تگ
	}
}
