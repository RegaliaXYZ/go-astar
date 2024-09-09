package main

import (
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
)

type Point struct {
	X, Y int
}

type Node struct {
	Point    Point
	Cost     float64 // Cost from the start node
	Priority float64 // f = g + h
	Parent   *Node   // Parent node to trace the path
}

type HeuristicFunc func(a, b Point) float64

func ManhattanDistance(start, end Point) float64 {
	return math.Abs(float64(start.X-end.X)) + math.Abs(float64(start.Y-end.Y))
}

// EuclideanDistance calculates the straight-line distance between two points.
func EuclideanDistance(start, end Point) float64 {
	return math.Sqrt(math.Pow(float64(start.X-end.X), 2) + math.Pow(float64(start.Y-end.Y), 2))
}

type AStar struct {
	board         [][]int
	heuristicFunc HeuristicFunc
	start         Point
	end           Point
}

func (a *AStar) SetStart(start Point) {
	a.start = start
}

func (a *AStar) SetEnd(end Point) {
	a.end = end
}

func (a *AStar) SetBoard(board [][]int) {
	a.board = board
}

func (a *AStar) SetHeuristic(heuristic HeuristicFunc) {
	a.heuristicFunc = heuristic
}

func (a *AStar) Randomize(width, height int) error {
	grid := make([][]int, height)
	zeroLocations := make([][2]int, 0)

	for i := range grid {
		grid[i] = make([]int, width)
		for j := range grid[i] {
			if rand.Float64() < 0.2 {
				grid[i][j] = 1
			} else {
				grid[i][j] = 0
				zeroLocations = append(zeroLocations, [2]int{i, j})
			}
		}
	}
	a.SetBoard(grid)
	if len(zeroLocations) < 2 {
		return fmt.Errorf("not enough 0 values found on the board")
	}
	randomIndex1 := rand.IntN(len(zeroLocations))

	loc1 := zeroLocations[randomIndex1]
	zeroLocations = append(zeroLocations[:randomIndex1], zeroLocations[randomIndex1+1:]...)
	a.SetStart(Point{
		X: loc1[1],
		Y: loc1[0],
	})

	randomIndex2 := rand.IntN(len(zeroLocations))
	loc2 := zeroLocations[randomIndex2]
	a.SetEnd(Point{
		X: loc2[1],
		Y: loc2[0],
	})
	a.SetHeuristic(ManhattanDistance)
	return nil
}

func (a *AStar) PrettyPrint(path []Point) {

	fmt.Println(path)
	for i, row := range a.board {
		// i = y
		for j := range row {

			// j = x
			if a.start.X == j && a.start.Y == i {
				fmt.Print("S ")
			} else if a.end.X == j && a.end.Y == i {
				fmt.Print("E ")
			} else if slices.Contains(path, Point{j, i}) {
				fmt.Print("X ")
			} else {
				fmt.Printf("%d ", a.board[i][j])
			}
		}
		fmt.Println()
	}
}

func (a *AStar) RawPath(path []Point) {
	if len(path) == 0 {
		fmt.Println("No path.")
		return
	}
	for _, p := range path {
		fmt.Printf("(%d, %d) -> ", p.X, p.Y)
	}
	fmt.Println("Goal")
}
func (a *AStar) FindPath() ([]Point, bool) {
	openList := []*Node{}
	closedSet := make(map[Point]bool)

	startNode := &Node{Point: a.start, Cost: 0, Priority: ManhattanDistance(a.start, a.end)}
	openList = append(openList, startNode)

	for len(openList) > 0 {
		// find the node with the lowest priority (f = g + h)
		currentIndex := 0
		for i, node := range openList {
			if node.Priority < openList[currentIndex].Priority {
				currentIndex = i
			}
		}
		current := openList[currentIndex]

		if current.Point == a.end {
			return a.reconstructPath(current), true
		}
		openList = append(openList[:currentIndex], openList[currentIndex+1:]...)
		closedSet[current.Point] = true

		// neighbors (up, down, left, right)
		neighbors := []Point{
			{current.Point.X + 1, current.Point.Y}, {current.Point.X - 1, current.Point.Y},
			{current.Point.X, current.Point.Y + 1}, {current.Point.X, current.Point.Y - 1},
		}
		for _, neighborPos := range neighbors {
			if !a.inBounds(neighborPos, a.board) || a.board[neighborPos.Y][neighborPos.X] == 1 {
				continue
			}

			if closedSet[neighborPos] {
				continue
			}

			tentativeCost := current.Cost + 1

			var neighborNode *Node
			exists := false
			for _, node := range openList {
				if node.Point == neighborPos {
					neighborNode = node
					exists = true
					break
				}
			}

			if !exists {
				neighborNode = &Node{Point: neighborPos}
				openList = append(openList, neighborNode)
			} else if tentativeCost >= neighborNode.Cost {
				continue
			}

			neighborNode.Cost = tentativeCost
			neighborNode.Priority = tentativeCost + ManhattanDistance(neighborPos, a.end)
			neighborNode.Parent = current
		}
	}
	return nil, false
}

func (a *AStar) reconstructPath(node *Node) []Point {
	path := []Point{}
	for n := node; n != nil; n = n.Parent {
		path = append([]Point{n.Point}, path...)
	}
	return path
}

func (a *AStar) inBounds(p Point, grid [][]int) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < len(grid[0]) && p.Y < len(grid)
}

func main() {
	// Example grid (0 = walkable, 1 = obstacle)
	a := AStar{}
	var r_width, r_height int
	for {
		r_width = rand.IntN(20)
		r_height = rand.IntN(20)
		if r_width > 1 && r_height > 1 {
			break
		}
	}

	a.Randomize(r_width, r_height)
	fmt.Println("------ RANDOMIZED BOARD ------")
	a.PrettyPrint(nil)
	fmt.Println("------ END OF BOARD ------")
	fmt.Println("Searching for path...")
	path, found := a.FindPath()
	if found {
		fmt.Println("Path found!")
	} else {
		fmt.Println("No path found...")
	}
	//fmt.Println(a.start, a.end)
	fmt.Println("------ PRINTING FINAL BOARD ------ ")
	a.PrettyPrint(path)
	fmt.Println("------ END OF BOARD ------")
	fmt.Println("Printing raw path taken.")
	a.RawPath(path)
}
