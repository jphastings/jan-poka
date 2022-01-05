//go:build libnova
package celestial

type Body string

const (
	Moon    Body = "moon"
	Sun     Body = "sun"
	Mercury Body = "mercury"
	Venus   Body = "venus"
	Mars    Body = "mars"
	Jupiter Body = "jupiter"
	Saturn  Body = "saturn"
	Neptune Body = "neptune"
	Uranus  Body = "uranus"
	Pluto   Body = "pluto"
)

var Bodies = map[string]Body{
	string(Moon):    Moon,
	string(Sun):     Sun,
	string(Mercury): Mercury,
	string(Venus):   Venus,
	string(Mars):    Mars,
	string(Jupiter): Jupiter,
	string(Saturn):  Saturn,
	string(Neptune): Neptune,
	string(Uranus):  Uranus,
	string(Pluto):   Pluto,
}
