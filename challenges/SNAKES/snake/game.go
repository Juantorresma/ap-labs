package snake

import (
	"time"
	"github.com/nsf/termbox-go"

)

var (
	pointsChan         = make(chan int)
	keyboardEventsChan = make(chan keyboardEvent)
	enemyMoveChan			 = make(chan direction)
	livesChan					 = make(chan int)
)

// Game type
type Game struct {
	arena  *arena
	pScore  int
	eScore  int
	isOver  bool
	toWin   int
	lives	  int
	win 		bool
	fN 			int
}

func initialSnake() *snake {
	return newSnake(RIGHT, []coord{
		coord{x: 1, y: 1},
		coord{x: 1, y: 2},
		coord{x: 1, y: 3},
		coord{x: 1, y: 4},
	})
}

func initialEnemy() *snake {
	return newSnake(LEFT, []coord{
		coord{x: 40, y: 1},
		coord{x: 40, y: 2},
		coord{x: 40, y: 3},
		coord{x: 40, y: 4},
	})
}

func initialArena(f int) *arena {
	return newArena(initialSnake(),initialEnemy(), pointsChan, livesChan, 20, 50, f)
}

func (g *Game) end() {
	g.isOver = true
}

func (g *Game) moveInterval() time.Duration {
	ms := 100 - g.pScore
	return time.Duration(ms) * time.Millisecond
}

func (g *Game) retry() {
	g.arena = initialArena(g.fN)
	g.pScore = 0
	g.eScore = 0
	g.toWin = len(g.arena.food)
	g.lives = 3
	g.isOver = false
	g.win = false
}

// NewGame creates new Game object
func NewGame(f int) *Game {
	a := initialArena(f)
	return &Game{arena: a, pScore: 0, eScore: 0, toWin: len(a.food), lives: 3, win: false, fN: f}
}

// Start starts the game
func (g *Game) Start() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	defer termbox.Close()

	go listenToKeyboard(keyboardEventsChan)
	go enemyMove(enemyMoveChan, g)

	if err := g.render(); err != nil {
		panic(err)
	}

mainloop:
	for {
		select {
		case l := <-livesChan:
			g.lives -= l
		case p := <-pointsChan:
			if(p>0){
				g.pScore += p
			} else{
				g.eScore -= p
			}
			g.toWin -= 1
			//fmt.Prinln("PUNTO")
			//fmt.Print(g.arena.enemies[0].length)

		case m := <- enemyMoveChan:
			g.arena.enemy.changeDir(m)

		case e := <-keyboardEventsChan:
			switch e.eventType {
			case MOVE:
				d := keyToDirection(e.key)
				g.arena.snake.changeDir(d)
			case RETRY:
				g.retry()
			case END:
				break mainloop
			}
		default:
			if !g.isOver {
				if err := g.arena.moveSnake(); err != nil {
					g.end()
				}
			}

			if !g.isOver {
				if err := g.arena.moveEnemy(); err != nil {
					g.end()
				}
			}

			if g.toWin==0{
				if g.pScore > g.eScore{
					g.win = true
				}
				g.end()
			}

			if g.lives<=0{
				g.end()
			}

			if err := g.render(); err != nil {
				panic(err)
			}

			time.Sleep(g.moveInterval())
		}
	}
}
