package parser

import (
	"os"
	"regexp"
	"strings"
)

type Plist struct {
	CFBundleIdentifier                           string
	CFBundleDisplayName                          string
	CFBundleVersion                              string
	CFBundleShortVersionString                   string
	CFBundleExecutable                           string
	CFBundleName                                 string
	NSPhotoLibraryUsageDescription               string
	NSCameraUsageDescription                     string
	NSMicrophoneUsageDescription                 string
	NSLocationWhenInUseUsageDescription          string
	NSLocationAlwaysAndWhenInUseUsageDescription string
	NSContactsUsageDescription                   string
	NSCalendarsUsageDescription                  string
	NSRemindersUsageDescription                  string
	NSBluetoothAlwaysUsageDescription            string
	NSBluetoothPeripheralUsageDescription        string
	NSAppleMusicUsageDescription                 string
	NSHealthShareUsageDescription                string
	NSHealthUpdateUsageDescription               string
	NSSiriUsageDescription                       string
	ITSAppUsesNonExemptEncryption                *bool
	UIRequiresFullScreen                         *bool
	UIApplicationSceneManifest                   map[string]interface{}
	UIBackgroundModes                            []string
	UILaunchStoryboardName                       string
	UIApplicationSupportsIndirectInputEvents     bool
	ITSAppTransportSecurity                      map[string]interface{}
	UIStatusBarStyle                             string
	UIUserInterfaceStyle                         string
	NSUserTrackingUsageDescription               string
}

func ParseInfoPlist(path string) (*Plist, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(data)
	plist := &Plist{}

	cfbundlePatterns := map[string]*regexp.Regexp{
		"CFBundleIdentifier":                           regexp.MustCompile(`<key>CFBundleIdentifier</key>\s*<string>([^<]+)</string>`),
		"CFBundleDisplayName":                          regexp.MustCompile(`<key>CFBundleDisplayName</key>\s*<string>([^<]+)</string>`),
		"CFBundleVersion":                              regexp.MustCompile(`<key>CFBundleVersion</key>\s*<string>([^<]+)</string>`),
		"CFBundleShortVersionString":                   regexp.MustCompile(`<key>CFBundleShortVersionString</key>\s*<string>([^<]+)</string>`),
		"CFBundleExecutable":                           regexp.MustCompile(`<key>CFBundleExecutable</key>\s*<string>([^<]+)</string>`),
		"CFBundleName":                                 regexp.MustCompile(`<key>CFBundleName</key>\s*<string>([^<]+)</string>`),
		"NSPhotoLibraryUsageDescription":               regexp.MustCompile(`<key>NSPhotoLibraryUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSCameraUsageDescription":                     regexp.MustCompile(`<key>NSCameraUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSMicrophoneUsageDescription":                 regexp.MustCompile(`<key>NSMicrophoneUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSLocationWhenInUseUsageDescription":          regexp.MustCompile(`<key>NSLocationWhenInUseUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSLocationAlwaysAndWhenInUseUsageDescription": regexp.MustCompile(`<key>NSLocationAlwaysAndWhenInUseUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSContactsUsageDescription":                   regexp.MustCompile(`<key>NSContactsUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSCalendarsUsageDescription":                  regexp.MustCompile(`<key>NSCalendarsUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSRemindersUsageDescription":                  regexp.MustCompile(`<key>NSRemindersUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSBluetoothAlwaysUsageDescription":            regexp.MustCompile(`<key>NSBluetoothAlwaysUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSBluetoothPeripheralUsageDescription":        regexp.MustCompile(`<key>NSBluetoothPeripheralUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSAppleMusicUsageDescription":                 regexp.MustCompile(`<key>NSAppleMusicUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSHealthShareUsageDescription":                regexp.MustCompile(`<key>NSHealthShareUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSHealthUpdateUsageDescription":               regexp.MustCompile(`<key>NSHealthUpdateUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSSiriUsageDescription":                       regexp.MustCompile(`<key>NSSiriUsageDescription</key>\s*<string>([^<]*)</string>`),
		"NSUserTrackingUsageDescription":               regexp.MustCompile(`<key>NSUserTrackingUsageDescription</key>\s*<string>([^<]*)</string>`),
	}

	for key, pattern := range cfbundlePatterns {
		if matches := pattern.FindStringSubmatch(content); len(matches) > 1 {
			switch key {
			case "CFBundleIdentifier":
				plist.CFBundleIdentifier = strings.TrimSpace(matches[1])
			case "CFBundleDisplayName":
				plist.CFBundleDisplayName = strings.TrimSpace(matches[1])
			case "CFBundleVersion":
				plist.CFBundleVersion = strings.TrimSpace(matches[1])
			case "CFBundleShortVersionString":
				plist.CFBundleShortVersionString = strings.TrimSpace(matches[1])
			case "CFBundleExecutable":
				plist.CFBundleExecutable = strings.TrimSpace(matches[1])
			case "CFBundleName":
				plist.CFBundleName = strings.TrimSpace(matches[1])
			case "NSPhotoLibraryUsageDescription":
				plist.NSPhotoLibraryUsageDescription = strings.TrimSpace(matches[1])
			case "NSCameraUsageDescription":
				plist.NSCameraUsageDescription = strings.TrimSpace(matches[1])
			case "NSMicrophoneUsageDescription":
				plist.NSMicrophoneUsageDescription = strings.TrimSpace(matches[1])
			case "NSLocationWhenInUseUsageDescription":
				plist.NSLocationWhenInUseUsageDescription = strings.TrimSpace(matches[1])
			case "NSLocationAlwaysAndWhenInUseUsageDescription":
				plist.NSLocationAlwaysAndWhenInUseUsageDescription = strings.TrimSpace(matches[1])
			case "NSContactsUsageDescription":
				plist.NSContactsUsageDescription = strings.TrimSpace(matches[1])
			case "NSCalendarsUsageDescription":
				plist.NSCalendarsUsageDescription = strings.TrimSpace(matches[1])
			case "NSRemindersUsageDescription":
				plist.NSRemindersUsageDescription = strings.TrimSpace(matches[1])
			case "NSBluetoothAlwaysUsageDescription":
				plist.NSBluetoothAlwaysUsageDescription = strings.TrimSpace(matches[1])
			case "NSBluetoothPeripheralUsageDescription":
				plist.NSBluetoothPeripheralUsageDescription = strings.TrimSpace(matches[1])
			case "NSAppleMusicUsageDescription":
				plist.NSAppleMusicUsageDescription = strings.TrimSpace(matches[1])
			case "NSHealthShareUsageDescription":
				plist.NSHealthShareUsageDescription = strings.TrimSpace(matches[1])
			case "NSHealthUpdateUsageDescription":
				plist.NSHealthUpdateUsageDescription = strings.TrimSpace(matches[1])
			case "NSSiriUsageDescription":
				plist.NSSiriUsageDescription = strings.TrimSpace(matches[1])
			case "NSUserTrackingUsageDescription":
				plist.NSUserTrackingUsageDescription = strings.TrimSpace(matches[1])
			}
		}
	}

	itsAppUsesNonExemptPattern := regexp.MustCompile(`<key>ITSAppUsesNonExemptEncryption</key>\s*<([a-z]+)/?>`)
	if matches := itsAppUsesNonExemptPattern.FindStringSubmatch(content); len(matches) > 1 {
		value := matches[1] == "true"
		plist.ITSAppUsesNonExemptEncryption = &value
	}

	uiRequiresFullScreenPattern := regexp.MustCompile(`<key>UIRequiresFullScreen</key>\s*<([a-z]+)/?>`)
	if matches := uiRequiresFullScreenPattern.FindStringSubmatch(content); len(matches) > 1 {
		value := matches[1] == "true"
		plist.UIRequiresFullScreen = &value
	}

	return plist, nil
}

func (p *Plist) HasCameraUsageDescription() bool {
	return p.NSCameraUsageDescription != ""
}

func (p *Plist) HasPhotoLibraryUsageDescription() bool {
	return p.NSPhotoLibraryUsageDescription != ""
}

func (p *Plist) HasLocationUsageDescription() bool {
	return p.NSLocationWhenInUseUsageDescription != "" || p.NSLocationAlwaysAndWhenInUseUsageDescription != ""
}

func (p *Plist) HasMicrophoneUsageDescription() bool {
	return p.NSMicrophoneUsageDescription != ""
}

func (p *Plist) HasContactsUsageDescription() bool {
	return p.NSContactsUsageDescription != ""
}

func (p *Plist) HasCalendarsUsageDescription() bool {
	return p.NSCalendarsUsageDescription != ""
}

func (p *Plist) HasUsageDescription(key string) bool {
	switch key {
	case "NSCameraUsageDescription":
		return p.HasCameraUsageDescription()
	case "NSPhotoLibraryUsageDescription":
		return p.HasPhotoLibraryUsageDescription()
	case "NSLocationWhenInUseUsageDescription", "NSLocationAlwaysAndWhenInUseUsageDescription":
		return p.HasLocationUsageDescription()
	case "NSMicrophoneUsageDescription":
		return p.HasMicrophoneUsageDescription()
	case "NSContactsUsageDescription":
		return p.HasContactsUsageDescription()
	case "NSCalendarsUsageDescription":
		return p.HasCalendarsUsageDescription()
	default:
		return false
	}
}

func (p *Plist) IsUsageDescriptionEmpty(key string) bool {
	switch key {
	case "NSCameraUsageDescription":
		return p.NSCameraUsageDescription == ""
	case "NSPhotoLibraryUsageDescription":
		return p.NSPhotoLibraryUsageDescription == ""
	case "NSLocationWhenInUseUsageDescription":
		return p.NSLocationWhenInUseUsageDescription == ""
	case "NSMicrophoneUsageDescription":
		return p.NSMicrophoneUsageDescription == ""
	case "NSContactsUsageDescription":
		return p.NSContactsUsageDescription == ""
	case "NSCalendarsUsageDescription":
		return p.NSCalendarsUsageDescription == ""
	default:
		return false
	}
}

func (p *Plist) IsEncryptionDeclarationSet() bool {
	return p.ITSAppUsesNonExemptEncryption != nil
}

func (p *Plist) IsEncryptionExempt() bool {
	return p.ITSAppUsesNonExemptEncryption != nil && !*p.ITSAppUsesNonExemptEncryption
}

func (p *Plist) GetFullScreenRequirement() bool {
	return p.UIRequiresFullScreen != nil && *p.UIRequiresFullScreen
}
