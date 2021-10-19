package config

// InjectionWords : InjectionWords
var InjectionWords []string

// InitInjectionWords : InjectionWords 초기화
func InitInjectionWords() {
	InjectionWords = []string{
		";", " AND ", " OR ",
	}
}
