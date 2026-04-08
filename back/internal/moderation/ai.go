// Package moderation provides a stub AI moderation check.
// Replace RunAICheck with a real API call when available.
package moderation

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

const aiPassThreshold = 0.40 // score below this → rejected

// Result is returned by RunAICheck.
type Result struct {
	Score     float64  `json:"score"`      // 0.0 (bad) – 1.0 (clean)
	Passed    bool     `json:"passed"`
	Flags     []string `json:"flags"`      // reasons for rejection
	CheckedAt time.Time
}

// bannedKeywords — слова, которые гарантированно срабатывают как флаги.
var bannedKeywords = map[string]string{
	"казино":    "gambling",
	"casino":    "gambling",
	"наркотик":  "drugs",
	"drug":      "drugs",
	"оружие":    "weapons",
	"weapon":    "weapons",
	"порно":     "adult_content",
	"porn":      "adult_content",
	"ставки":    "gambling",
	"betting":   "gambling",
	"обман":     "fraud",
	"scam":      "fraud",
	"отмывание": "fraud",
	"laundering": "fraud",
}

// RunAICheck imitates an AI moderation check.
// title and description are checked for banned keywords;
// a random component simulates model uncertainty.
//
// To integrate a real model: replace this function body with an HTTP call
// to your inference endpoint and map the response to Result.
func RunAICheck(title, description string) Result {
	text := strings.ToLower(title + " " + description)

	flags := []string{}
	seen := map[string]bool{}
	for kw, flag := range bannedKeywords {
		if strings.Contains(text, kw) && !seen[flag] {
			flags = append(flags, flag)
			seen[flag] = true
		}
	}

	// Base score: 1.0 minus a penalty per flag, plus small random noise.
	score := 1.0 - float64(len(flags))*0.35
	noise := (rand.Float64() - 0.5) * 0.10 // ±5%
	score += noise
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	score = roundFloat(score, 3)

	return Result{
		Score:     score,
		Passed:    score >= aiPassThreshold && len(flags) == 0,
		Flags:     flags,
		CheckedAt: time.Now(),
	}
}

// FlagsToJSON serialises flags slice to a JSON string for storage.
func FlagsToJSON(flags []string) string {
	b, _ := json.Marshal(flags)
	return string(b)
}

func roundFloat(v float64, decimals int) float64 {
	p := 1.0
	for i := 0; i < decimals; i++ {
		p *= 10
	}
	return float64(int(v*p+0.5)) / p
}
