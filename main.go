package main

import "fmt"

// Point 结构体保持不变
type Point struct {
	Row int
	Col int
}

// 定义类型约束
type MazeElement comparable

// 1. 在 MazeSolver 中增加 wall 字段
type MazeSolver[T MazeElement] struct {
	maze    [][]T
	rows    int
	cols    int
	parents map[Point]Point
	wall    T // 告诉求解器，“墙”长什么样
}

func (s *MazeSolver[T]) reconstructPath(start, end Point) []Point {
	path := []Point{}
	for current := end; current != start; current = s.parents[current] {
		path = append(path, current)
	}
	path = append(path, start)

	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

// 2. NewMazeSolver 接收 wall 作为参数
func NewMazeSolver[T MazeElement](maze [][]T, wall T) *MazeSolver[T] {
	return &MazeSolver[T]{
		maze:    maze,
		rows:    len(maze),
		cols:    len(maze[0]),
		parents: make(map[Point]Point),
		wall:    wall, // 初始化 wall
	}
}

func (s *MazeSolver[T]) findEndNode(endValue T) (Point, bool) {
	for r := 0; r < s.rows; r++ {
		for c := 0; c < s.cols; c++ {
			if s.maze[r][c] == endValue {
				return Point{r, c}, true
			}
		}
	}
	return Point{}, false
}

func (s *MazeSolver[T]) SolveBFS(start Point, endValue T) ([]Point, bool) {
	end, found := s.findEndNode(endValue)
	if !found {
		fmt.Printf("错误：在迷宫中未找到值为 %v 的终点\n", endValue)
		return nil, false
	}

	queue := []Point{start}
	visited := make(map[Point]bool)
	visited[start] = true

	directions := []Point{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
	}

	for len(queue) > 0 {
		currentNode := queue[0]
		queue = queue[1:]
		if currentNode == end {
			return s.reconstructPath(start, end), true
		}
		for _, dir := range directions {
			next := Point{Row: currentNode.Row + dir.Row, Col: currentNode.Col + dir.Col}
			if next.Row < 0 || next.Row >= s.rows || next.Col < 0 || next.Col >= s.cols {
				continue
			}
			// 3. 使用 s.wall 进行比较，这是 T == T 的安全比较
			if s.maze[next.Row][next.Col] == s.wall {
				continue
			}
			if visited[next] {
				continue
			}
			visited[next] = true
			s.parents[next] = currentNode
			queue = append(queue, next)
		}
	}
	return nil, false
}

func main() {
	maze := [][]rune{
		{'0', '0', '1', '1', '1', '1', '1', '0', '0', '0'},
		{'1', '0', '1', '0', '0', '0', '1', '0', '1', '0'},
		{'1', '0', '1', '0', '1', '0', '1', '0', '1', '0'},
		{'1', '0', '1', '0', '1', '0', '1', '0', '1', '0'},
		{'1', '0', '1', '2', '1', '2', '1', '0', '1', '0'},
		{'1', '0', '1', '0', '1', '0', '0', '0', '1', '0'},
		{'1', '0', '1', '0', '1', '0', '1', '0', '1', '0'},
		{'1', '0', '1', '0', '1', '0', '1', '0', '1', '0'},
		{'1', '0', '1', '0', '1', '0', '1', '0', '1', '0'},
		{'1', '0', '0', '0', '1', '0', '#', '0', '1', '0'}, // 目标 '#'
	}

	// 4. 创建 Solver 时，告诉它墙是 '1'
	solver := NewMazeSolver(maze, '1')
	startNode := Point{Row: 0, Col: 0}
	endValue := '#' // rune 类型

	fmt.Printf("正在从起点 %v 寻找通往终点 (值为 '%c') 的最短路径...\n", startNode, endValue)
	path, found := solver.SolveBFS(startNode, endValue)

	if found {
		fmt.Println("成功找到最短路径！路径如下:")
		// 保留你需要的箭头打印逻辑
		for i, p := range path {
			if i > 0 {
				fmt.Print(" -> ")
			}
			fmt.Printf("{%d, %d}", p.Row, p.Col)

			// 对于起点(i=0)，它没有父节点，需要特殊处理避免panic
			if i == 0 {
				continue
			}

			parent := solver.parents[p]
			if parent.Col == p.Col && parent.Row+1 == p.Row {
				fmt.Printf("⬇️")
			}
			if parent.Col+1 == p.Col && parent.Row == p.Row {
				fmt.Printf("➡️")
			}
			if parent.Col-1 == p.Col && parent.Row == p.Row {
				fmt.Printf("⬅️")
			}
			if parent.Col == p.Col && parent.Row-1 == p.Row {
				fmt.Printf("⬆️")
			}
		}
		fmt.Println()
	} else {
		fmt.Println("未找到从起点到终点的有效路径。")
	}
}
