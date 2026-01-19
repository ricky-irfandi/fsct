package advanced

import (
	"math"

	"github.com/ricky-irfandi/fsct/internal/report"
)

type SecurityScore struct {
	Overall     float64 `json:"overall"`
	Android     float64 `json:"android"`
	iOS         float64 `json:"ios"`
	Flutter     float64 `json:"flutter"`
	Security    float64 `json:"security"`
	Policy      float64 `json:"policy"`
	CodeQuality float64 `json:"codeQuality"`
	Grade       string  `json:"grade"`
	Summary     string  `json:"summary"`
}

type ScoreBreakdown struct {
	Category string  `json:"category"`
	Passed   int     `json:"passed"`
	Failed   int     `json:"failed"`
	Score    float64 `json:"score"`
	Weight   float64 `json:"weight"`
	Weighted float64 `json:"weighted"`
}

const (
	androidWeight     = 0.25
	iosWeight         = 0.25
	flutterWeight     = 0.15
	securityWeight    = 0.20
	policyWeight      = 0.10
	codeQualityWeight = 0.05
)

func CalculateScore(findings []report.Finding, totalChecks int) *SecurityScore {
	breakdown := calculateBreakdown(findings)

	overall := 0.0
	for _, b := range breakdown {
		overall += b.Weighted
	}

	grade := calculateGrade(overall)

	return &SecurityScore{
		Overall:     math.Round(overall*100) / 100,
		Android:     math.Round(breakdown[0].Score*100) / 100,
		iOS:         math.Round(breakdown[1].Score*100) / 100,
		Flutter:     math.Round(breakdown[2].Score*100) / 100,
		Security:    math.Round(breakdown[3].Score*100) / 100,
		Policy:      math.Round(breakdown[4].Score*100) / 100,
		CodeQuality: math.Round(breakdown[5].Score*100) / 100,
		Grade:       grade,
		Summary:     generateSummary(grade, overall),
	}
}

func calculateBreakdown(findings []report.Finding) []ScoreBreakdown {
	categoryCounts := make(map[string]int)
	categoryFailed := make(map[string]int)

	categories := []string{"Android", "iOS", "Flutter", "Security", "Policy", "Code Quality"}

	for _, f := range findings {
		cat := categorize(f.ID)
		categoryCounts[cat]++
		categoryFailed[cat]++
	}

	breakdown := make([]ScoreBreakdown, 0)

	for _, cat := range categories {
		passed := categoryCounts[cat] - categoryFailed[cat]
		total := categoryCounts[cat]
		score := 1.0
		if total > 0 {
			score = float64(passed) / float64(getTotalChecksForCategory(cat))
		}

		var weight float64
		switch cat {
		case "Android":
			weight = androidWeight
		case "iOS":
			weight = iosWeight
		case "Flutter":
			weight = flutterWeight
		case "Security":
			weight = securityWeight
		case "Policy":
			weight = policyWeight
		case "Code Quality":
			weight = codeQualityWeight
		}

		breakdown = append(breakdown, ScoreBreakdown{
			Category: cat,
			Passed:   passed,
			Failed:   categoryFailed[cat],
			Score:    score,
			Weight:   weight,
			Weighted: score * weight,
		})
	}

	return breakdown
}

func categorize(checkID string) string {
	switch {
	case len(checkID) >= 4 && checkID[:3] == "AND":
		return "Android"
	case len(checkID) >= 4 && checkID[:3] == "IOS":
		return "iOS"
	case len(checkID) >= 4 && checkID[:3] == "FLT":
		return "Flutter"
	case len(checkID) >= 4 && checkID[:3] == "SEC":
		return "Security"
	case len(checkID) >= 4 && checkID[:3] == "POL":
		return "Policy"
	case len(checkID) >= 4 && checkID[:3] == "COD" ||
		len(checkID) >= 4 && checkID[:3] == "TST" ||
		len(checkID) >= 4 && checkID[:3] == "LIN" ||
		len(checkID) >= 4 && checkID[:3] == "DOC" ||
		len(checkID) >= 4 && checkID[:3] == "PER":
		return "Code Quality"
	default:
		return "Other"
	}
}

func getTotalChecksForCategory(category string) int {
	switch category {
	case "Android":
		return 12
	case "iOS":
		return 12
	case "Flutter":
		return 8
	case "Security":
		return 5
	case "Policy":
		return 5
	case "Code Quality":
		return 33
	default:
		return 1
	}
}

func calculateGrade(score float64) string {
	if score >= 0.95 {
		return "A+"
	} else if score >= 0.90 {
		return "A"
	} else if score >= 0.85 {
		return "A-"
	} else if score >= 0.80 {
		return "B+"
	} else if score >= 0.75 {
		return "B"
	} else if score >= 0.70 {
		return "B-"
	} else if score >= 0.65 {
		return "C+"
	} else if score >= 0.60 {
		return "C"
	} else if score >= 0.55 {
		return "C-"
	} else if score >= 0.50 {
		return "D+"
	} else if score >= 0.45 {
		return "D"
	} else {
		return "F"
	}
}

func generateSummary(grade string, score float64) string {
	switch {
	case grade[0] == 'A':
		return "Excellent! Your app meets high compliance standards. Ready for app store submission."
	case grade[0] == 'B':
		return "Good compliance. A few issues should be addressed before submission."
	case grade[0] == 'C':
		return "Moderate compliance. Several issues need attention before submission."
	case grade[0] == 'D':
		return "Low compliance. Significant issues must be resolved before submission."
	default:
		return "Critical compliance issues. Major changes required before submission."
	}
}
