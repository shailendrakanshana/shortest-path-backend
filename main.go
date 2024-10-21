package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Coordinates struct {
	Start [2]int `json:"start"` // Use backticks for struct tags
	End   [2]int `json:"end"`   // Use backticks for struct tags
}

type PathResponse struct {
	Path [][2]int `json:"path"` // Use backticks for struct tags
}

// CORS Middleware
func enableCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, or specify allowed domains
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// BFS function to find the shortest path
func bfs(grid [][]int, start [2]int, end [2]int) [][2]int {
	queue := [][2]int{start}
	visited := make(map[[2]int]bool)
	visited[start] = true
	prev := make(map[[2]int][2]int) // To reconstruct the path

	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if current == end {
			break // Path found
		}

		for _, dir := range directions {
			newPos := [2]int{current[0] + dir[0], current[1] + dir[1]}

			// Bounds checking and visit checking
			if newPos[0] >= 0 && newPos[0] < 20 && newPos[1] >= 0 && newPos[1] < 20 && !visited[newPos] {
				visited[newPos] = true
				queue = append(queue, newPos)
				prev[newPos] = current // Track the previous node
			}
		}
	}

	// Reconstruct path from end to start
	var path [][2]int
	for at := end; at != start; at = prev[at] {
		path = append(path, at)
	}
	path = append(path, start) // Add start point
	reverse(path)               // Reverse to get correct order

	return path
}

// Reverse the path
func reverse(path [][2]int) {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
}

// API handler
func findPathHandler(w http.ResponseWriter, r *http.Request) {
	var coords Coordinates
	err := json.NewDecoder(r.Body).Decode(&coords)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	grid := make([][]int, 20)
	for i := range grid {
		grid[i] = make([]int, 20)
	}

	path := bfs(grid, coords.Start, coords.End)

	if len(path) > 0 {
		// Log the path
		fmt.Println("Path found:", path)
		json.NewEncoder(w).Encode(PathResponse{Path: path})
	} else {
		http.Error(w, "No path found", http.StatusNotFound)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/find-path", findPathHandler)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", enableCors(mux)) // Wrap with CORS middleware
}
