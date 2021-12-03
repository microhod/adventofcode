// Package christmas spreads festive cheer!!!
package christmas

import (
	"strings"

	"github.com/mgutz/ansi"
	"github.com/moul/sapin"
)

var (
	blinkRed   = ansi.ColorFunc("red+Bbh")
	blinkGreen = ansi.ColorFunc("green+Bbh")
	blinkGold  = ansi.ColorFunc("yellow+Bbh")
)

func Tree() string {
	sapin := sapin.NewSapin(3)
	sapin.AddStar()
	sapin.AddBalls(20)
	sapin.AddGarlands(5)
	sapin.Colorize()

	return sapin.String()
}

func Lights() string {
	lights := []string{}
	colours := []func(string)string {
		blinkRed,
		blinkGreen,
		blinkGold,
	}
	for i := 0; i < 28; i++ {
		lights = append(lights, colours[i % 3]("•"))
	}

	return strings.Join(lights, " ")
}
