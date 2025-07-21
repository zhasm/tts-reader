package config

import "testing"

func TestIsSupportedLang(t *testing.T) {
	if !IsSupportedLang("fr") {
		t.Error("Expected 'fr' to be supported")
	}
	if IsSupportedLang("xx") {
		t.Error("Expected 'xx' to be unsupported")
	}
}

func TestGetLang(t *testing.T) {
	lang, found := GetLang("fr")
	if !found || lang.Name != "fr" {
		t.Error("Expected to find 'fr' language")
	}
	_, found = GetLang("xx")
	if found {
		t.Error("Expected not to find 'xx' language")
	}
}

func TestGetFlag(t *testing.T) {
	Language = "fr"
	flag := GetFlag()
	if flag == "" {
		t.Error("Expected flag for 'fr'")
	}
	Language = "xx"
	flag = GetFlag()
	if flag != "" {
		t.Error("Expected empty flag for unsupported language")
	}
}
