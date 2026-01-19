package formatter

import (
	"encoding/json"
	"fmt"

	"github.com/ricky-irfandi/fsct/internal/report"
)

type SARIFFormatter struct{}

type SARIFReport struct {
	Schema  string   `json:"$schema"`
	Version string   `json:"version"`
	Run     SARIFRun `json:"runs"`
}

type SARIFRun struct {
	Tool    SARIFTool     `json:"tool"`
	Results []SARIFResult `json:"results,omitempty"`
}

type SARIFTool struct {
	Driver SARIFDriver `json:"driver"`
}

type SARIFDriver struct {
	Name            string      `json:"name"`
	Version         string      `json:"version"`
	SemanticVersion string      `json:"semanticVersion,omitempty"`
	Rules           []SARIFRule `json:"rules,omitempty"`
}

type SARIFRule struct {
	ID               string       `json:"id"`
	Name             string       `json:"name"`
	ShortDescription SARIFMessage `json:"shortDescription"`
	FullDescription  SARIFMessage `json:"fullDescription,omitempty"`
	HelpURI          string       `json:"helpUri,omitempty"`
	SeverityLevel    string       `json:"defaultConfiguration,omitempty"`
}

type SARIFMessage struct {
	Text string `json:"text"`
}

type SARIFResult struct {
	RuleID    string          `json:"ruleId"`
	RuleIndex int             `json:"ruleIndex,omitempty"`
	Level     string          `json:"level"`
	Message   SARIFMessage    `json:"message"`
	Locations []SARIFLocation `json:"locations,omitempty"`
	Artifacts []SARIFArtifact `json:"artifacts,omitempty"`
}

type SARIFLocation struct {
	PhysicalLocation SARIFPhysicalLocation `json:"physicalLocation"`
}

type SARIFPhysicalLocation struct {
	ArtifactLocation SARIFArtifactLocation `json:"artifactLocation"`
	Region           SARIFRegion           `json:"region,omitempty"`
}

type SARIFArtifactLocation struct {
	URI string `json:"uri"`
}

type SARIFRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
	EndLine     int `json:"endLine,omitempty"`
	EndColumn   int `json:"endColumn,omitempty"`
}

type SARIFArtifact struct {
	Location SARIFArtifactLocation `json:"location"`
	Contents string                `json:"contents,omitempty"`
	Encoding string                `json:"encoding,omitempty"`
}

func (f *SARIFFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	sarifReport := SARIFReport{
		Schema:  "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json",
		Version: "2.1.0",
		Run: SARIFRun{
			Tool: SARIFTool{
				Driver: SARIFDriver{
					Name:            "FSCT",
					Version:         "1.0.0",
					SemanticVersion: "1.0.0",
				},
			},
		},
	}

	for _, finding := range results {
		result := SARIFResult{
			RuleID: finding.ID,
			Level:  mapSeverity(finding.Severity),
			Message: SARIFMessage{
				Text: finding.Message,
			},
			Locations: []SARIFLocation{
				{
					PhysicalLocation: SARIFPhysicalLocation{
						ArtifactLocation: SARIFArtifactLocation{
							URI: finding.File,
						},
						Region: SARIFRegion{
							StartLine: finding.Line,
						},
					},
				},
			},
		}
		sarifReport.Run.Results = append(sarifReport.Run.Results, result)
	}

	return json.MarshalIndent(sarifReport, "", "  ")
}

func (f *SARIFFormatter) GetExtension() string {
	return "sarif"
}

func mapSeverity(severity report.Severity) string {
	switch severity {
	case report.SeverityHigh:
		return "error"
	case report.SeverityWarning:
		return "warning"
	default:
		return "note"
	}
}

type GitHubSummaryFormatter struct{}

type GitHubSummaryOutput struct {
	Title       string             `json:"title"`
	Summary     string             `json:"summary"`
	Conclusions []GitHubConclusion `json:"conclusions,omitempty"`
}

type GitHubConclusion struct {
	Conclusion string `json:"conclusion"`
	Summary    string `json:"summary"`
}

func (f *GitHubSummaryFormatter) Format(results []report.Finding, summary report.Summary) ([]byte, error) {
	output := GitHubSummaryOutput{
		Title: "FSCT Compliance Report",
	}

	if summary.High > 0 {
		output.Summary = fmt.Sprintf("## Compliance Issues Found\n\n**High Severity:** %d\n**Warning:** %d\n**Info:** %d\n\n### Action Required\n\n%d high severity issues need to be addressed before submitting to the app store.", summary.High, summary.Warning, summary.Info, summary.High)
	} else if summary.Warning > 0 {
		output.Summary = fmt.Sprintf("## Compliance Report\n\n**High Severity:** 0\n**Warning:** %d\n**Info:** %d\n\n### Recommendations\n\n%d warning issues should be reviewed.", summary.Warning, summary.Info, summary.Warning)
	} else {
		output.Summary = fmt.Sprintf("## Compliance Passed\n\nAll checks passed! Your app is ready for submission.\n\n**Passed:** %d\n**High Severity:** 0\n**Warning:** 0\n**Info:** 0", summary.Passed)
	}

	return json.MarshalIndent(output, "", "  ")
}

func (f *GitHubSummaryFormatter) GetExtension() string {
	return "json"
}
