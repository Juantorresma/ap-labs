package main

import (
	"fmt"
	"math/rand"
	"time"
	"sync"
	"math"
	)

/*Constants*/
const (
	PrintColor = "\033[38;5;%dm%s\033[39;49m\n"
	startingVal = 88
	yOffset = 2
	acceleration = -9.8
)

/*Global Variables*/
var voxelMap [][][]int
var balls = [][]int{}
var states = []string{}
/*Balls Data*/
var numBalls = 5
var iterations = 10
var maxHeight = 100
var n = 15
/*Ball Counter*/
var north = 0
var west = 0
var south = 0
var east = 0
var island = 0
/*Layout Helpers*/
var tOffset int

func main() {
	print("\033[H\033[2J")
	enterInputs()
	print("\033[H\033[2J")
	rand.Seed(time.Now().UnixNano())
	voxelMap = generateVoxelMap()
	initBalls()
	finalOffset()
	displayVoxelMap()
	time.Sleep(time.Second)
	go updateEverything()
	startRain()
}

func enterInputs() {
	/*Maps Size*/
    fmt.Println("Size of map (nxn):")
    _, err := fmt.Scan(&n)
    if err != nil {
    	n = 5
    }
    /*Max Height*/
    fmt.Println("Max height::")
    _, err = fmt.Scan(&maxHeight)
    if err != nil {
    	maxHeight = 10
    }
    /*Iterations*/
    fmt.Println("Iterations:")
    _, err = fmt.Scan(&iterations)
    if err != nil {
    	iterations = 1
    }
    /*Number of Balls per Iteration*/
    fmt.Println("Balls per iteration::")
    _, err = fmt.Scan(&numBalls)
    if err != nil {
    	numBalls = 1
    }
}

func generateVoxelMap() [][][]int {
	newVoxelMap := [][][]int{}
	for i := 0; i < n; i++ {
		row := [][]int{}
		for j := 0; j < n; j++ {
			if i == 0 || i == n - 1 || j == 0 || j == n - 1 {
				s := []int{0, 0}
				row = append(row, s)
			} else {
				s := []int{rand.Intn(maxHeight), 0}
				row = append(row, s)
			}
		}
		newVoxelMap = append(newVoxelMap, row)
	}
	return newVoxelMap
}

func displayVoxelMap() {
	for i := 0; i < len(voxelMap); i++ {
		for j := 0; j < len(voxelMap); j++ {
			fmt.Printf("\033[%d;%dH", i  + yOffset, (j * 3) + 70)
			if voxelMap[i][j][1] == 0 {
				fmt.Printf(PrintColor, startingVal + (5 - (voxelMap[i][j][0]) % 5), "*")
			} else {
				fmt.Printf(PrintColor, 3, "o")
			}
		}
	}
}

func startRain() {
	/*
	This function creates the goroutines for each ball.
	There is a sleep of 3 seconds between every iteration.
	Every iteration creates numBalls goroutines.
	sync.WaitGroup is used to make sure every goroutine finishes
	before ending the program.
	*/
	var wg sync.WaitGroup
	k := 0
	for i := 0; i < iterations; i++ {
		for j := 0; j < numBalls; j++ {
			wg.Add(1)
			go rain(k, &wg)
			k++
		}
		time.Sleep(time.Second * 3)
	}
	wg.Wait()
	updateScreen()
}

func rain(index int, wg *sync.WaitGroup) {
	/*
	This function handles the behavior of every ball.
	When this function is done it calls sync.WaitGroup.Done(),
	which is the same sync.WaitGroup used in startRain(), to decrease the counter
	on the sync.WaitGroup.
	*/
	x, y := ballsPosition(index)
	fall(index, x, y)

	for true {
		balls[index][1] = y
		balls[index][2] = x
		balls[index][3] = voxelMap[x][y][0]
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(1000)))
		if x == 0 || y == 0 || x == (n - 1) || y == (n - 1) {
			states[index] = seaBalls(index, x, y)
			voxelMap[x][y][1] = 0
			break
		}else if y < (n - 1) && voxelMap[x][y + 1][0] < voxelMap[x][y][0] && voxelMap[x][y + 1][1] == 0 {
			voxelMap[x][y][1] = 0
			voxelMap[x][y+1][1] = 1
			calculateVelocity(index, voxelMap[x][y + 1][0], "s")
			y += 1
		} else if x < (n - 1) && voxelMap[x + 1][y][0] < voxelMap[x][y][0] && voxelMap[x + 1][y][1] == 0 {
			voxelMap[x][y][1] = 0
			voxelMap[x+1][y][1] = 1
			calculateVelocity(index, voxelMap[x + 1][y][0], "e")
			x += 1
		} else if y > 0 && voxelMap[x][y - 1][0] < voxelMap[x][y][0] && voxelMap[x][y - 1][1] == 0 {
			voxelMap[x][y][1] = 0
			voxelMap[x][y-1][1] = 1
			calculateVelocity(index, voxelMap[x][y - 1][0], "n")
			y -= 1
		} else if x > 0 && voxelMap[x - 1][y][0] < voxelMap[x][y][0] && voxelMap[x - 1][y][1] == 0 {
			voxelMap[x][y][1] = 0
			voxelMap[x-1][y][1] = 1
			calculateVelocity(index, voxelMap[x - 1][y][0], "w")
			x -= 1
		} else {
			states[index] = "Lake"
			island += 1
			break
		}

	}
	balls[index][1] = y
	balls[index][2] = x
	calculateVelocity(index, voxelMap[x][y][0], "_")
	wg.Done()
}

func ballsPosition(index int) (int, int) {
	/*
	Gives an (X, Y) position to the ball in the map, if that
	location is occupied by another ball, then it generates new
	coordinates until it finds an empty location.
	*/
	x := 1 + rand.Intn(n - 2)
	y := 1 + rand.Intn(n - 2)
	for voxelMap[x][y][1] == 1 {
		x = 1 + rand.Intn(n - 2)
		y = 1 + rand.Intn(n - 2)
	}
	voxelMap[x][y][1] = 1
	balls[index][1] = x
	balls[index][2] = y
	return x, y
}

func fall(index int, x int, y int) {
	/*
	Controlls the falling velocity and Z position of each ball.
	*/
	states[index] = "Falling"
	t := 500
	totalTime := 0
	for balls[index][3] >= voxelMap[x][y][0] {
		balls[index][3] = int(balls[index][6] * (totalTime / 1000)) + balls[index][3]
		balls[index][6] = balls[index][6] + int(acceleration * float64(totalTime / 1000))
		totalTime += t
		time.Sleep(time.Millisecond * time.Duration(t))
	}
	balls[index][3] = voxelMap[x][y][0]
	states[index] = "Moving"
}

func seaBalls(index int, x int, y int)string {
	/*
	Tracks the number of balls that go to the sea
	from the north, west, south, east.
	*/
	if x == 0 { north += 1; return "North" }
	if y == 0 { west += 1; return "West" }
	if x == (n - 1) { south += 1; return "South" }
	if y == (n - 1) {east += 1; return "East"}
	return ""
}

func calculateVelocity(index int, zPos int, direction string) {
	/*
	Calculates x and y velocities for every ball.
	*/
	angle := math.Atan((float64(balls[index][3]) / float64(zPos)))
	if direction == "n" {
		balls[index][5] = int(float64(balls[index][6]) * math.Cos(angle))
		balls[index][6] = int(float64(balls[index][6]) * math.Sin(angle))
	} else if direction == "w" {
		balls[index][4] = int(float64(balls[index][6]) * math.Cos(angle))
		balls[index][6] = int(float64(balls[index][6]) * math.Sin(angle))
	} else if direction == "s" {
		balls[index][4] = -int(float64(balls[index][6]) * math.Cos(angle))
		balls[index][6] = int(float64(balls[index][6]) * math.Sin(angle))
	} else if direction == "e" {
		balls[index][5] = -int(float64(balls[index][6]) * math.Cos(angle))
		balls[index][6] = int(float64(balls[index][6]) * math.Sin(angle))
	} else {
		balls[index][4] = 0
		balls[index][5] = 0
		balls[index][6] = 0
	}
}

func initBalls() {
	for i := 0; i < (numBalls * iterations); i++ {
		ball := []int{i, 0, 0, rand.Intn(maxHeight * 20) + maxHeight, 0, 0, 0}
		balls = append(balls, ball)
		states = append(states, "Waiting")
	}
}

func displayBallInfo() {
	fmt.Printf("\033[2;0H")
	fmt.Printf("ID | Xp | Yp |  Zp  | Vx | Vy |  Vz  |  State  |")
	fmt.Printf("\033[3;0H")
	fmt.Printf("---+----+----+------+----+----+------+---------+")
	for i := 0; i < len(balls); i++ {
		/*Print Ball ID*/
		fmt.Printf("\033[%d;0H", i+4)
		fmt.Printf("%d", balls[i][0])
		/*Separator*/
		fmt.Printf("\033[%d;4H", i+4)
		fmt.Printf("|")
		/*Print Ball X Position*/
		fmt.Printf("\033[%d;6H", i+4)
		fmt.Printf("%d", balls[i][1])
		/*Separator*/
		fmt.Printf("\033[%d;9H", i+4)
		fmt.Printf("|")
		/*Print Ball Y Position*/
		fmt.Printf("\033[%d;11H", i+4)
		fmt.Printf("%d", balls[i][2])
		/*Separator*/
		fmt.Printf("\033[%d;14H", i+4)
		fmt.Printf("|")
		/*Print Ball Z Position*/
		fmt.Printf("\033[%d;16H", i+4)
		fmt.Printf("%d", balls[i][3])
		/*Separator*/
		fmt.Printf("\033[%d;21H", i+4)
		fmt.Printf("|")
		/*Print Ball X Velocity*/
		fmt.Printf("\033[%d;23H", i+4)
		fmt.Printf("%d", balls[i][4])
		/*Separator*/
		fmt.Printf("\033[%d;26H", i+4)
		fmt.Printf("|")
		/*Print Ball Y Velocity*/
		fmt.Printf("\033[%d;28H", i+4)
		fmt.Printf("%d", balls[i][5])
		/*Separator*/
		fmt.Printf("\033[%d;31H", i+4)
		fmt.Printf("|")
		/*Print Ball Z Velocity*/
		fmt.Printf("\033[%d;33H", i+4)
		fmt.Printf("%d", balls[i][6])
		/*Separator*/
		fmt.Printf("\033[%d;38H", i+4)
		fmt.Printf("|")
		/*Print Ball State*/
		fmt.Printf("\033[%d;40H", i+4)
		fmt.Printf(states[i])
		/*Separator*/
		fmt.Printf("\033[%d;48H", i+4)
		fmt.Printf("|")
	}
	fmt.Printf("\033[%d;50H", 2)
	fmt.Printf("Island:")
	fmt.Printf("\033[%d;58H", 2)
	fmt.Printf("%d", island)

	fmt.Printf("\033[%d;50H", 3)
	fmt.Printf("North:")
	fmt.Printf("\033[%d;58H", 3)
	fmt.Printf("%d", north)

	fmt.Printf("\033[%d;50H", 4)
	fmt.Printf("West:")
	fmt.Printf("\033[%d;58H", 4)
	fmt.Printf("%d", west)

	fmt.Printf("\033[%d;50H", 5)
	fmt.Printf("South:")
	fmt.Printf("\033[%d;58H", 5)
	fmt.Printf("%d", south)

	fmt.Printf("\033[%d;50H", 6)
	fmt.Printf("East:")
	fmt.Printf("\033[%d;58H", 6)
	fmt.Printf("%d", east)
}



func updateScreen() {
	print("\033[H\033[2J")
  	displayVoxelMap()
  	displayBallInfo()
  	fmt.Printf("\033[%d;0H", tOffset + 5)
}

func updateEverything() {
  for {
    <-time.After(200 * time.Millisecond)
    go updateScreen()
  }
}

func finalOffset() {
	if (numBalls * iterations) > n { tOffset = (numBalls * iterations) } else { tOffset = n }
}









