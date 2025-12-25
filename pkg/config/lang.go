package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/zhasm/tts-reader/pkg/logger"
	"gopkg.in/yaml.v3"
)

var supportedLangs []string

type Lang struct {
	Name     string `yaml:"name"`
	NameFUll string `yaml:"full_name"`
	Reader   string `yaml:"reader"`
	Gender   string `yaml:"gender"`
	Flag     string `yaml:"flag"`
	Regex    string `yaml:"regex"`
}

type LangConfig struct {
	Langs []Lang `yaml:"langs"`
}

// Define the supported languages
var Langs = DefaultLangs

var DefaultLangs = []Lang{
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

func LoadConfig() {
	var configPath string

	// 1. Check command line flag
	if ConfigFile != "" {
		configPath = ConfigFile
	} else {
		// 2. Check current directory
		if _, err := os.Stat("./tts-langs.yml"); err == nil {
			configPath = "./tts-langs.yml"
		} else {
			// 3. Check home directory
			home, err := os.UserHomeDir()
			if err == nil {
				if _, err := os.Stat(filepath.Join(home, ".tts-langs.yml")); err == nil {
					configPath = filepath.Join(home, ".tts-langs.yml")
				} else if _, err := os.Stat(filepath.Join(home, ".config", "tts-langs.yml")); err == nil {
					// 4. Check ~/.config/tts-langs.yml
					configPath = filepath.Join(home, ".config", "tts-langs.yml")
				}
			}
		}
	}

	if configPath == "" {
		logger.LogWarn("No config file found, using default languages")
		Langs = DefaultLangs
	} else {
		data, err := os.ReadFile(configPath)
		if err != nil {
			logger.LogWarn("Error reading config file %s: %v. Using defaults.", configPath, err)
			Langs = DefaultLangs
		} else {
			var config LangConfig
			if err := yaml.Unmarshal(data, &config); err != nil {
				logger.LogWarn("Error parsing config file %s: %v. Using defaults.", configPath, err)
				Langs = DefaultLangs
			} else {
				if len(config.Langs) == 0 {
					logger.LogWarn("Config file %s has no languages. Using defaults.", configPath)
					Langs = DefaultLangs
				} else {
					Langs = config.Langs
					logger.LogInfo("Loaded configuration from %s", configPath)
				}
			}
		}
	}

	initSupportedLangs()
}

func GenerateConfigFile() {
	config := LangConfig{
		Langs: DefaultLangs,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		logger.LogError("Error marshaling default config: %v", err)
		os.Exit(1)
	}

	filename := "tts-langs.yml"
	if err := os.WriteFile(filename, data, 0644); err != nil {
		logger.LogError("Error writing config file %s: %v", filename, err)
		os.Exit(1)
	}
	logger.LogInfo("Generated default config file: %s", filename)
}
