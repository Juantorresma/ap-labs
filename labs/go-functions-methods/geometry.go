//Geometry.go Juan Torres A01227885

//I had to change the name of my package to run it
package main

//Here are my imports
import (
	"math"
	"math/rand"
	"os"
	"fmt"
	"strconv"
)

//I had to change to lowe case, given
type Point struct{ x, y float64 }


// traditional function, given
func Distance(p, q Point) float64 {
	return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}


// same thing, but as a method of the Point type, given
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X()-p.X(), q.Y()-p.Y())
}


// A Path is a journey connecting the points with straight lines, given
type Path []Point


// Distance returns the distance traveled along the path, modified to print my data, i did all here and changed return to 0
func (path Path) Distance() float64 {
	sum := 0.0
	fmt.Print("- ")
	for i := range path {
		if i > 0 {
			fmt.Print(path[i-1].Distance(path[i]), " + ")
			sum += path[i-1].Distance(path[i])
		}
	}
	fmt.Print(path[len(path)-1].Distance(path[0]))
	sum += path[len(path)-1].Distance(path[0])
	fmt.Print(" = ")
	fmt.Print(sum)
	return 0
}


//this func is to find the x coord of the point
func (point Point) X() float64 {
	return point.x
}



//this func is to find the y coord of the point
func (point Point) Y() float64 {
	return point.y
}


//this func is to check if it is on a segment
func onSegment(p, q, r Point) bool {
    
    //math function to check if it is on the same segment of a line
	if q.X() <= math.Max(p.X(), r.X()) && q.X() >= math.Min(p.X(), r.X()) && q.Y() <= math.Max(p.Y(), r.Y()) && q.Y() >= math.Min(p.Y(), r.Y()) {
		return true;
	}
	return false;
}


//this checks for the orientations
func orientation(p, q, r Point) int {
	val := (q.Y() - p.Y() * (r.X() - q.X()) - (q.X() - p.X()) * (r.Y() - q.Y()))
	min := -0.0001
	max := 0.0001
	if (val > min && val < max) {
		return 0
	}
	if (val > max) {
		return 1
	} else {
		return 2
	}
}



//this one validates if there are any interceptions, usses onSegment
func intersect(p1, q1, p2, q2 Point) bool {
	o1 := orientation(p1, q1, p2)
	o2 := orientation(p1, q1, q2)
	o3 := orientation(p2, q2, p1)
	o4 := orientation(p2, q2, q1)

	if o1 != o2 && o3 != o4 {
		return true
	}

	if o1 == 0 && onSegment(p1, p2, q1) {
		return true
	}
	if o2 == 0 && onSegment(p1, q2, q1) {
		return true
	}
	if o3 == 0 && onSegment(p2, p1, q2) {
		return true
	}
	if o4 == 0 && onSegment(p2, q1, q2) {
		return true
	}

	return false
}


//this is where we generate our random points and print them
func genRand (Paths Path, numSides int) []Point {
	for i := 0; i < numSides; i++ {
                Paths[i].x = ((rand.Float64() * 200)-100)
                Paths[i].y = ((rand.Float64() * 200)-100)
                fmt.Println("- ( ",Paths[i].X(),", ",Paths[i].Y(),")")
        }
	return Paths
}


//Here we call our functions
func main() {
    if len(os.Args) < 2 {
            fmt.Println("Usage: go run geometry.go <sides>")
    }
	numSides, _ := strconv.Atoi(os.Args[1])
    fmt.Println("- Generating a [",numSides,"] sides figure")
	fmt.Println("- Figure's vertices")
	Paths := make( Path,numSides)
	Paths = genRand(Paths, numSides)
	for intersect(Paths[0], Paths[1], Paths[2], Paths[3]) {
		Paths = genRand(Paths, numSides)
	}
	fmt.Println("- Figure's Perimeter")
	Paths.Distance()
}
