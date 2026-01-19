package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Match struct {
	File    string
	Line    int
	Content string
	Pattern string
}

type DartScanner struct {
	basePath string
}

func NewDartScanner(basePath string) *DartScanner {
	return &DartScanner{basePath: basePath}
}

func (s *DartScanner) ScanDir(pattern string) ([]Match, error) {
	var matches []Match

	err := filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if !strings.HasSuffix(info.Name(), ".dart") {
			return nil
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}

		fileMatches, err := s.scanFile(path, re)
		if err != nil {
			return err
		}

		matches = append(matches, fileMatches...)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return matches, nil
}

func (s *DartScanner) ScanDirMulti(patterns []string) ([]Match, error) {
	var allMatches []Match

	for _, pattern := range patterns {
		matches, err := s.ScanDir(pattern)
		if err != nil {
			return nil, err
		}
		allMatches = append(allMatches, matches...)
	}

	return allMatches, nil
}

func (s *DartScanner) scanFile(path string, re *regexp.Regexp) ([]Match, error) {
	var matches []Match

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if re.MatchString(line) {
			relPath, err := filepath.Rel(s.basePath, path)
			if err != nil {
				relPath = path
			}
			matches = append(matches, Match{
				File:    relPath,
				Line:    i + 1,
				Content: strings.TrimSpace(line),
				Pattern: re.String(),
			})
		}
	}

	return matches, nil
}

func (s *DartScanner) FindPatternInFile(path string, pattern string) ([]Match, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return s.scanFile(path, re)
}

func ScanDartFiles(basePath string, pattern string) ([]Match, error) {
	scanner := NewDartScanner(basePath)
	return scanner.ScanDir(pattern)
}

func ScanDartFilesMulti(basePath string, patterns []string) ([]Match, error) {
	scanner := NewDartScanner(basePath)
	return scanner.ScanDirMulti(patterns)
}

func FindAPIKeys(basePath string) ([]Match, error) {
	patterns := []string{
		`api[_]?key\s*[:=]\s*["']?[a-zA-Z0-9_-]+`,
		`apiKey\s*[:=]\s*["']?[a-zA-Z0-9_-]+`,
		`API[_]?KEY\s*[:=]\s*["']?[a-zA-Z0-9_-]+`,
		`sk-[a-zA-Z0-9]{20,}`,
		`pk_live_[a-zA-Z0-9]{20,}`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindHardcodedPasswords(basePath string) ([]Match, error) {
	patterns := []string{
		`password\s*[:=]\s*["'][^"']+`,
		`"password123"`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindPrintStatements(basePath string) ([]Match, error) {
	patterns := []string{
		`\bprint\s*\(`,
		`\bdebugPrint\s*\(`,
		`\bdeveloper\.log\s*\(`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindHTTPURLs(basePath string) ([]Match, error) {
	pattern := `http://[a-zA-Z0-9][-a-zA-Z0-9]*(\.[a-zA-Z0-9][-a-zA-Z0-9]*)+`
	return ScanDartFiles(basePath, pattern)
}

func FindPrivacyPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`privacy`,
		`privacyPolicy`,
		`privacy_policy`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindTermsPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`terms`,
		`termsOfService`,
		`tos`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindLoginPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`login`,
		`signIn`,
		`sign_in`,
		`authenticate`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindDeleteAccountPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`deleteAccount`,
		`delete_account`,
		`removeAccount`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindLogoutPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`signOut`,
		`logout`,
		`logOut`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindPasswordRecoveryPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`forgotPassword`,
		`resetPassword`,
		`forgot_password`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}

func FindTestEmailPatterns(basePath string) ([]Match, error) {
	patterns := []string{
		`test@`,
		`example.com`,
		`@example.org`,
	}
	return ScanDartFilesMulti(basePath, patterns)
}
