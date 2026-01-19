package report

type Severity string

const (
	SeverityInfo    Severity = "INFO"
	SeverityWarning Severity = "WARNING"
	SeverityHigh    Severity = "HIGH"
)

type Finding struct {
	ID         string   `json:"id"`
	Severity   Severity `json:"severity"`
	Title      string   `json:"title"`
	Message    string   `json:"message"`
	File       string   `json:"file,omitempty"`
	Line       int      `json:"line,omitempty"`
	Suggestion string   `json:"suggestion,omitempty"`
}

type Summary struct {
	High    int `json:"high"`
	Warning int `json:"warning"`
	Info    int `json:"info"`
	Passed  int `json:"passed"`
}

type Report struct {
	Version   string    `json:"version"`
	Timestamp string    `json:"timestamp"`
	Project   string    `json:"project"`
	Summary   Summary   `json:"summary"`
	Findings  []Finding `json:"findings"`
}
