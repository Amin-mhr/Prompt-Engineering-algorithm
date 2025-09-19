package prompts

import (
	"regexp"
	"sort"
	"strings"
)

type Scored struct {
	Ex    Example
	Score float64
}

// یک شباهت خیلی ساده: اشتراک کیورد/تگ‌ها (Jaccard)
func SimilarityScore(question string, ex Example) float64 {
	qTokens := tokenize(question)
	exTokens := map[string]struct{}{}
	for _, t := range ex.Tags {
		exTokens[strings.ToLower(t)] = struct{}{}
	}
	//for _, w := range tokenize(ex.Description) {
	//	exTokens[w] = struct{}{}
	//}

	inter := 0
	for w := range qTokens {
		if _, ok := exTokens[w]; ok {
			inter++
		}
	}
	union := len(qTokens) + len(exTokens) - inter
	if union == 0 {
		return 0
	}
	return float64(inter) / float64(union)
}

func tokenize(s string) map[string]struct{} {
	s = strings.ToLower(s)
	re := regexp.MustCompile(`[a-zA-Z]+`)
	words := re.FindAllString(s, -1)
	set := make(map[string]struct{}, len(words))
	for _, w := range words {
		set[w] = struct{}{}
	}
	return set
}

func PickTopK(question string, repo []Example, k int) []Example {
	scored := make([]Scored, 0, len(repo))
	for _, ex := range repo {
		scored = append(scored, Scored{Ex: ex, Score: SimilarityScore(question, ex)})
	}
	sort.Slice(scored, func(i, j int) bool { return scored[i].Score > scored[j].Score })
	if k > len(scored) {
		k = len(scored)
	}
	out := make([]Example, 0, k)
	for i := 0; i < k; i++ {
		if scored[i].Score <= 0 {
			break
		}
		out = append(out, scored[i].Ex)
	}
	return out
}
