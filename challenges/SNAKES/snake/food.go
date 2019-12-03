package snake

import (
	"os"
	"strings"
)

type food struct {
	emoji        rune
	points, x, y int
}

func newFood(x, y int) food {
	if strings.Contains(os.Getenv("LANG"), "UTF-8") {
		return food{
			points: 10,
			emoji:  '🌟',
			x:      x,
			y:      y,
		}
	} else {
		return food{
			points: 10,
			emoji:  '@',
			x:      x,
			y:      y,
		}
	}
}
