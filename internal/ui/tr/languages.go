package tr

// Language is a language the UI is translated into.
type Language struct {
	// Code is the locale, as the message file is named.
	Code string
	// Name is the language's name in itself.
	Name string
}

// Languages lists the translations, in the order the language picker
// shows them.
var Languages = []Language{
	{"en", "English"},
	{"bg", "Български"},
	{"cs", "Čeština"},
	{"da", "Dansk"},
	{"de", "Deutsch"},
	{"et", "Eesti"},
	{"el", "Ελληνικά"},
	{"es", "Español"},
	{"fr", "Français"},
	{"hr", "Hrvatski"},
	{"it", "Italiano"},
	{"lv", "Latviešu"},
	{"lt", "Lietuvių"},
	{"hu", "Magyar"},
	{"nl", "Nederlands"},
	{"nb", "Norsk bokmål"},
	{"pl", "Polski"},
	{"pt", "Português"},
	{"ro", "Română"},
	{"ru", "Русский"},
	{"sk", "Slovenčina"},
	{"sl", "Slovenščina"},
	{"fi", "Suomi"},
	{"sv", "Svenska"},
	{"uk", "Українська"},
}
