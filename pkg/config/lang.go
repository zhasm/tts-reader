package config

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

var supportedLangs []string

type Lang struct {
	Name     string
	NameFUll string
	Reader   string
	Gender   string
	Flag     string
	Regex    string
}

// Define the supported languages
// TODO load from config file
var Langs = []Lang{
	{
		Name:     "fr",
		NameFUll: "fr-FR",
		Reader:   "fr-FR-DeniseNeural",
		Gender:   "Male",
		Flag:     "🇫🇷",
		Regex:    "[a-zA-ZÀ-ÿ]+",
	},
	{
		Name:     "pl",
		NameFUll: "pl-PL",
		Reader:   "pl-PL-AgnieszkaNeural",
		Gender:   "Female",
		Flag:     "🇵🇱",
		Regex:    "[a-zA-ZąćęłńóśźżĄĆĘŁŃÓŚŹŻ]+",
	},
	{
		Name:     "jp",
		NameFUll: "ja-JP",
		Reader:   "ja-JP-MayuNeural",
		Gender:   "Female",
		Flag:     "🇯🇵",
		Regex:    "[ぁ-んァ-ン一-龯]+",
	},
	{
		Name:     "en",
		NameFUll: "en-US",
		Reader:   "en-GB-HollieNeural",
		Gender:   "Female",
		Flag:     "🇺🇸",
		Regex:    "[a-zA-Z]+",
	},
}

// Checks if a language is supported
func IsSupportedLang(name string) bool {
	for _, l := range Langs {
		if l.Name == name {
			return true
		}
	}
	return false
}

// Returns the Lang struct for a given language name, or (zero, false) if not found
func GetLang(name string) (Lang, bool) {
	for _, l := range Langs {
		if l.Name == name {
			return l, true
		}
	}
	return Lang{}, false
}

// GetRegex returns the regex of given language.
func ValidateLangRegex(langName, content string) (bool, error) {
	for _, l := range Langs {
		if l.Name == langName {
			ok, err := regexp.MatchString(l.Regex, content)
			if err != nil {
				return false, fmt.Errorf("error matching regex: %w", err)
			}
			return ok, nil
		}
	}
	return false, fmt.Errorf("not supported language: %s", langName)
}

// GetFlagByName returns the flag emoji for the given language name.
// If the language is not supported, it returns an empty string.
func GetFlagByName(name string) string {
	for _, l := range Langs {
		if l.Name == name {
			return l.Flag
		}
	}
	return ""
}

func GetAllLangShortNames() []string {
	langNames := make([]string, len(Langs))
	for i, l := range Langs {
		langNames[i] = l.Name
	}
	return langNames
}

func GetAllLangShortNamesStr() string {
	names := GetAllLangShortNames()
	slices.Sort(names)
	return strings.Join(names, ", ")
}

// Add this function to initialize supportedLangs
func initSupportedLangs() {
	supportedLangs = make([]string, len(Langs))
	for i, l := range Langs {
		supportedLangs[i] = l.Name
	}
}
