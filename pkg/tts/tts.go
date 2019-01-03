package tts

type Engine interface {
	Speak(string) error
}
