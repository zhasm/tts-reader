package main

var supportedLangs []string

type Lang struct {
	Name     string
	NameFUll string
	Reader   string
	Gender   string
	Flag     string
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
	},
	{
		Name:     "pl",
		NameFUll: "pl-PL",
		Reader:   "pl-PL-AgnieszkaNeural",
		Gender:   "Female",
		Flag:     "🇵🇱",
	},
	{
		Name:     "jp",
		NameFUll: "ja-JP",
		Reader:   "ja-JP-MayuNeural",
		Gender:   "Female",
		Flag:     "🇯🇵",
	},
	{
		Name:     "en",
		NameFUll: "en-US",
		Reader:   "en-GB-HollieNeural",
		Gender:   "Female",
		Flag:     "🇺🇸",
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

func GetFlag() string {
	for _, l := range Langs {
		if l.Name == Language {
			return l.Flag
		}
	}
	return ""
}
