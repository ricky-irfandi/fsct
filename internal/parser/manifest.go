package parser

import (
	"encoding/xml"
	"os"
	"regexp"
	"strings"
)

type AndroidManifest struct {
	Package         string
	VersionCode     string
	VersionName     string
	Debuggable      string
	AllowBackup     string
	UsesPermissions []UsesPermission
	Activities      []Activity
	Queries         []Queries
}

type UsesPermission struct {
	Name string
}

type Activity struct {
	Name          string
	Exported      string
	IntentFilters []IntentFilter
}

type IntentFilter struct {
	Actions []Action
}

type Action struct {
	Name string
}

type Queries struct {
	Packages []Package
}

type Package struct {
	Name string
}

func getAttrValue(attrs []xml.Attr, name string) string {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value
		}
		if strings.HasPrefix(attr.Name.Local, "android:") && strings.TrimPrefix(attr.Name.Local, "android:") == name {
			return attr.Value
		}
	}
	return ""
}

func ParseAndroidManifest(path string) (*AndroidManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	manifest := &AndroidManifest{
		UsesPermissions: make([]UsesPermission, 0),
		Activities:      make([]Activity, 0),
		Queries:         make([]Queries, 0),
	}

	xmlDecoder := xml.NewDecoder(strings.NewReader(content))

	for {
		token, err := xmlDecoder.Token()
		if err != nil {
			break
		}

		switch elem := token.(type) {
		case xml.StartElement:
			name := getAttrValue(elem.Attr, "name")
			exported := getAttrValue(elem.Attr, "exported")
			debuggable := getAttrValue(elem.Attr, "debuggable")
			allowBackup := getAttrValue(elem.Attr, "allowBackup")

			switch elem.Name.Local {
			case "manifest":
				manifest.Package = getAttrValue(elem.Attr, "package")
				manifest.VersionCode = getAttrValue(elem.Attr, "versionCode")
				manifest.VersionName = getAttrValue(elem.Attr, "versionName")
			case "uses-permission":
				manifest.UsesPermissions = append(manifest.UsesPermissions, UsesPermission{
					Name: name,
				})
			case "activity":
				manifest.Activities = append(manifest.Activities, Activity{
					Name:     name,
					Exported: exported,
				})
			case "queries":
				manifest.Queries = append(manifest.Queries, Queries{Packages: make([]Package, 0)})
			case "package":
				if len(manifest.Queries) > 0 {
					lastQuery := &manifest.Queries[len(manifest.Queries)-1]
					lastQuery.Packages = append(lastQuery.Packages, Package{Name: name})
				}
			case "application":
				manifest.Debuggable = debuggable
				manifest.AllowBackup = allowBackup
			}
		}
	}

	return manifest, nil
}

func (m *AndroidManifest) GetPackageName() string {
	return m.Package
}

func (m *AndroidManifest) GetDebuggable() bool {
	return strings.ToLower(m.Debuggable) == "true"
}

func (m *AndroidManifest) GetAllowBackup() bool {
	return strings.ToLower(m.AllowBackup) != "false"
}

func (m *AndroidManifest) HasPermission(permission string) bool {
	for _, p := range m.UsesPermissions {
		if strings.Contains(p.Name, permission) {
			return true
		}
	}
	return false
}

func (m *AndroidManifest) GetActivitiesWithIntentFilters() []Activity {
	return m.Activities
}

func (m *AndroidManifest) ActivityNeedsExported() []Activity {
	var activities []Activity
	for _, a := range m.Activities {
		if a.Exported == "" {
			activities = append(activities, a)
		}
	}
	return activities
}

func (m *AndroidManifest) HasInternetPermission() bool {
	return m.HasPermission("INTERNET")
}

func (m *AndroidManifest) HasCameraPermission() bool {
	return m.HasPermission("CAMERA")
}

func (m *AndroidManifest) HasLocationPermission() bool {
	return m.HasPermission("LOCATION")
}

func (m *AndroidManifest) HasQueryPackage(pkg string) bool {
	for _, q := range m.Queries {
		for _, p := range q.Packages {
			if strings.Contains(p.Name, pkg) {
				return true
			}
		}
	}
	return false
}

type GradleConfig struct {
	ApplicationID    string
	MinSDKVersion    string
	TargetSDKVersion string
	VersionCode      string
	VersionName      string
	NdkVersion       string
	KotlinVersion    string
	FlutterVersion   string
}

func ParseGradleFile(path string) (*GradleConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	config := &GradleConfig{}

	patterns := map[string]*regexp.Regexp{
		"applicationId":    regexp.MustCompile(`applicationId\s+["']([^"']+)["']`),
		"minSdkVersion":    regexp.MustCompile(`minSdkVersion\s+(\d+)`),
		"targetSdkVersion": regexp.MustCompile(`targetSdkVersion\s+(\d+)`),
		"versionCode":      regexp.MustCompile(`versionCode\s+(\d+)`),
		"versionName":      regexp.MustCompile(`versionName\s+["']([^"']+)["']`),
		"ndkVersion":       regexp.MustCompile(`ndkVersion\s+["']([^"']+)["']`),
	}

	for key, pattern := range patterns {
		if matches := pattern.FindStringSubmatch(content); len(matches) > 1 {
			switch key {
			case "applicationId":
				config.ApplicationID = matches[1]
			case "minSdkVersion":
				config.MinSDKVersion = matches[1]
			case "targetSdkVersion":
				config.TargetSDKVersion = matches[1]
			case "versionCode":
				config.VersionCode = matches[1]
			case "versionName":
				config.VersionName = matches[1]
			case "ndkVersion":
				config.NdkVersion = matches[1]
			}
		}
	}

	return config, nil
}

type BuildGradleKts struct {
	ApplicationID    string
	MinSDKVersion    string
	TargetSDKVersion string
	VersionCode      string
	VersionName      string
}

func ParseGradleKtsFile(path string) (*BuildGradleKts, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	config := &BuildGradleKts{}

	patterns := map[string]*regexp.Regexp{
		"applicationId":    regexp.MustCompile(`applicationId\s*=\s*["']([^"']+)["']`),
		"minSdkVersion":    regexp.MustCompile(`minSdkVersion\s*=\s*["']?(\d+)["']?`),
		"targetSdkVersion": regexp.MustCompile(`targetSdkVersion\s*=\s*["']?(\d+)["']?`),
		"versionCode":      regexp.MustCompile(`versionCode\s*=\s*["']?(\d+)["']?`),
		"versionName":      regexp.MustCompile(`versionName\s*=\s*["']([^"']+)["']`),
	}

	for key, pattern := range patterns {
		if matches := pattern.FindStringSubmatch(content); len(matches) > 1 {
			switch key {
			case "applicationId":
				config.ApplicationID = matches[1]
			case "minSdkVersion":
				config.MinSDKVersion = matches[1]
			case "targetSdkVersion":
				config.TargetSDKVersion = matches[1]
			case "versionCode":
				config.VersionCode = matches[1]
			case "versionName":
				config.VersionName = matches[1]
			}
		}
	}

	return config, nil
}
