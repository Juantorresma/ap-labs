package snake

import (
	"math/rand"
	"time"
)

type arena struct {
	food       []food
	snake      *snake
	//enemies    []*snake
	enemy 		 *snake
	hasFood    func(*arena, coord) int
	height     int
	width      int
	pointsChan chan (int)
	livesChan  chan (int)
}

func newArena(s, e *snake, p chan (int), l chan (int), h, w, f int) *arena {
	rand.Seed(time.Now().UnixNano())

	a := &arena{
		snake:      s,
		enemy:			e,
		height:     h,
		width:      w,
		pointsChan: p,
		livesChan: l,
		hasFood:    hasFood,
		food: make([]food,f),
		//enemies: make([]*snake,nE),
	}

	a.placeFood(f)
	//a.placeEnemies(nE)

	return a
}

func (a *arena) moveSnake() error {
	if err := a.snake.move(); err != nil {
		return err
	}

	if a.sLeftArena() {
		return a.snake.kill()
	}

	if a.sTouchEnemy(a.snake,a.enemy) {
		go a.enemyHit(1)
	}

	if f:=a.hasFood(a, a.snake.head()); f>=0{
		go a.addPoints(a.food[f].points)
		a.snake.length++
		a.food = append(a.food[:f], a.food[f+1:]...)
	}

	return nil
}


func (a *arena) moveEnemy() error {
	if err := a.enemy.move(); err != nil {
		return nil //Change to to err to kill enemy when touches itself
	}

	if a.sLeftArena() {
		return a.enemy.kill()
	}
	if a.sTouchEnemy(a.enemy,a.snake) {
		go a.enemyHit(1)
	}

	if f:=a.hasFood(a, a.enemy.head()); f>=0{
		go a.addPoints(-a.food[f].points)
		a.enemy.length++
		a.food = append(a.food[:f], a.food[f+1:]...)
	}

	return nil
}

func (a *arena) sLeftArena() bool {
	h := a.snake.head()
	return h.x > a.width || h.y > a.height || h.x < 0 || h.y < 0
}

func (a *arena) sTouchEnemy(s1, s2 *snake) bool {
	h := s1.head()

	return s2.isOnPos(h)
}

func (a *arena) eLeftArena() bool {
	h := a.enemy.head()
	return h.x > a.width || h.y > a.height || h.x < 0 || h.y < 0
}

func (a *arena) addPoints(p int) {
	a.pointsChan <- p
}

func (a *arena) enemyHit( l int) {
	a.livesChan <- l
}

func (a *arena) placeFood(n int) {
	var x, y int
	for i:=0;i<n;i++{

		for {
			x = rand.Intn(a.width)
			y = rand.Intn(a.height)

			if !a.isFull(coord{x: x, y: y}) {
				break
			}
		}

		a.food[i] = newFood(x, y)
	}
}

func hasFood(a *arena, c coord) int {
	for i:=0;i<len(a.food);i++{
		if c.x == a.food[i].x && c.y == a.food[i].y{
			return i
		}
	}
	return -1
}

func (a *arena) isFull(c coord) bool {
	return a.snake.isOnPos(c) || a.enemy.isOnPos(c)
}
