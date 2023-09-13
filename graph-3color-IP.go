package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Message struct {
	Sender    string
	Recipient string
	Data      Transcript
}

// Create a struct to store the communication transcript
// Some of the fields might not be necessary
type Transcript struct {
	Msg      string // Message exchanged between the prover & verifier
	Vertex1  int    // To store the vertices that the verifier requests the prover to show their coloring
	Vertex2  int
	Verified bool // Check if the communication is verified (and accepted) by the verifier when checking the vertices
	Color1   int
	Color2   int
	ErrorMsg error // Store the error messages that occur while communicating (from the prover or verifier)
	//vertexColor

}

// -------------------------------------------- Part 1: Graph Coloring Algorithm ------------------------------------

// Assume we have a set of integers (to represent the colors) in the range 1 to k, say, 1 to 10
// We need the minimum set of positive numbers (colors) to assign the vertices such that no adjacent vertex share the same-
// -integer (color); in this situation it's called a clash
// c(v) belongs to {1,2,....,k} such that c(u) != c(v) for all {u,v} in E and K is minimal
// Small example: if we have three colors as the minimum and the vertex v is assigned to the third color, then c(v) = 3

// If all the vertices are assigned colors, then the coloring is called complete. Otherwise, the coloring is partial.
// A proper coloring is when every pair of adjacent coloring do not clash (have the same coloring)

// n choices of colors, each n vertex will have n choices --> n^(n) candidate solution for the exhaustive approach
// We'll use a backtracking solution

// Backtracking implementation:

// Function that checks if it's OK to assign a color to a particular vertex
func isProperColor(g [][]int, n int, color int, intColors []int) bool {
	for i := 0; i < len(g); i++ { // Iterate through all the vertices & check if they have the same color
		if g[n][i] == 1 && intColors[i] == color {
			return false // return false if the vertices 'clash' (adjacent vertices with the same color)
		}
	}
	return true // Return true to indicate that it's a proper coloring (no adjacent vertices have the same color)
}

// Function that recursively explores possible colorings of the graph
func graphColor(g [][]int, intColors []int, n int, k int) bool {
	if n == len(g) { // current vertex n equals to the number of the graph
		return true // All vertices are (validly) colored in this case
	}
	// For each color in the range 1 to k
	for c := 1; c <= k; c++ {
		if isProperColor(g, n, c, intColors) { // check if it's alright to color the current vertex n
			intColors[n] = c                      //
			if graphColor(g, intColors, n+1, k) { // Recursivly call the graphColor on the next vertex (n+1)
				return true // Returns true in the case we've found a proper/valid coloring
			}
			intColors[n] = 0 // If false was returned (not a valid color), backtrack by setting intColors to 0 and try the next color
		}
	}
	return false // If all colors are not valid for the vertex, return false to backtrack furthur
}

func graphInput(g [][]int, intColors []int, k int, permutation []int) (properColor map[int]int, err error) {

	for i, v := range intColors {
		intColors[i] = v % k // Assigning colors within the range [1,2,3]. In other words, using only three colors.
	}

	if graphColor(g, intColors, 0, k) {
		// Create a map (dictionary) to store the vertex and its corresponding proper color
		properColor = make(map[int]int) //map[vertex:color]
		for i, color := range intColors {
			properColor[i+1] = permutation[color%k]
		}
		return properColor, nil
	} else {
		return nil, fmt.Errorf("No proper coloring possible.\n")
	}
}

// -------------------------------------Part 2: Prover-Verifier --------------------------------------------------------

// Function for the Prover to permute the colors
func coloringPermutation(n int) []int {
	perm := rand.Perm(n)

	for p := range perm {
		perm[p]++ // Shift indices to start the count from 1 rather than 0
	}
	return perm
}

// Prover Go-routine function
func Prover(g [][]int, ptv, vtp chan Message, wg *sync.WaitGroup) {
	// Defer control mechanism - delays the execution of the 'wg.Done()' statement until the rest of the function is complete
	defer wg.Done()
	// Choose at random the coloring assignment & color the graph
	k := 3 // k-color, where k is minimal; in the case of graph 3-color,  we require k = 3.
	intColors := make([]int, len(g))

	perm := coloringPermutation(k)

	// The certificate is a map that contains the vertices as keys and coloring as the respective values.
	certificate, err := graphInput(g, intColors, k, perm)
	if err != nil {
		fmt.Println("Error: ", err) // Error will print in case that the graph is not three-colorable
		return
	}

	fmt.Println("Prover's certificate (only known to him):", certificate)
	fmt.Println()
	// Iterate through the certificate map and print the proof
	for vertex, color := range certificate {
		fmt.Printf("Vertex: %d, Color: %d\n", vertex, color)
	}
	fmt.Println("\n")

	// pi(coloring) --> Commit to them (lock them) --> Send it to the verifier -- this is not implemented yet
	// If the edge is in E, then the Prover reveals the colors of the endpoints belonging to that edge

	response := <-vtp                                                                    // receiving Verifier's request of the edge and storing it in a struct variable called response.
	fmt.Printf("%s to %s: %v\n", response.Sender, response.Recipient, response.Data.Msg) // Access struct fields and print

	// Check if the theirs an edge in the graph between the vertices requested by the verifier
	// Accesss the edge (u,v) from the Transcript sent by the verifier
	var colorProof Message
	if g[response.Data.Vertex1][response.Data.Vertex2] == 1 { // of there's an edge

		fmt.Printf("%s to %s: Yep, (%d, %d) ∈ E. %s\n", response.Recipient, response.Sender, response.Data.Vertex1+1, response.Data.Vertex2+1, colorProof.Data.Msg)
		// If the edge in in E, then send the verifier the color of the chosen vertices through the channel
		//range through the proof map and get the verifier-selected vertex and send the coloring

		var p1, p2 int
		// Checking if the key (vertex) exist in the map (certificate), then store its value (color)
		if color1, ok1 := certificate[response.Data.Vertex1+1]; ok1 {
			p1 = color1
		}
		if color2, ok2 := certificate[response.Data.Vertex2+1]; ok2 {
			p2 = color2
		}
		colorProof = Message{
			Sender:    "Prover",
			Recipient: "Verifier",
			Data: Transcript{
				Msg:      "Coloring of your selected vertices: ",
				Vertex1:  response.Data.Vertex1, // The verifier's equested edge (vertices)
				Vertex2:  response.Data.Vertex2,
				Verified: true,
				Color1:   p1, // colors of the vertices
				Color2:   p2,
				ErrorMsg: nil,
			},
		}

		ptv <- colorProof // Send the colors that correspond to the vertex

	} else { // if there's no edge

		colorProof = Message{
			Sender:    "Prover",
			Recipient: "Verifier",
			Data: Transcript{
				Msg:      fmt.Sprintf("%s to %s: Sorry, (%d, %d) ∉ E. No coloring for you!\n", response.Recipient, response.Sender, response.Data.Vertex1+1, response.Data.Vertex2+1),
				Vertex1:  response.Data.Vertex1,
				Vertex2:  response.Data.Vertex2,
				Verified: false,
				Color1:   0,
				Color2:   0,
				ErrorMsg: nil,
			},
		}
		ptv <- colorProof
	}

}

// Verifier Go-routine function
func Verifier(g [][]int, ptv, vtp chan Message, wg *sync.WaitGroup) {
	defer wg.Done()

	// Verifier chooses a random edge that belongs in the set E
	// Accessing edges from the graph g (randomly)
	// Initialize the RNG with a seed for that^
	rand.Seed(time.Now().UnixNano())

	verticesN := len(g)       // get the number of vertices (which is 10)
	u := rand.Intn(verticesN) // randomly choose the vertices within the graph's length
	v := rand.Intn(verticesN)

	transcript := Message{
		Sender:    "Verifier",
		Recipient: "Prover",
		Data: Transcript{
			Msg:      fmt.Sprintf("Show me the endpoints of the edge (%d, %d)\n", u+1, v+1),
			Vertex1:  u,
			Vertex2:  v,
			Verified: true,
			Color1:   0,
			Color2:   0,
			ErrorMsg: nil,
		},
	}

	vtp <- transcript // Verifier sends the randomly selected edge e=(u,v) to the prover

	//Verifier receives the coloring from the prover after it's validated
	Vcolors := <-ptv

	if Vcolors.Data.Verified != false {

		fmt.Printf("%s to %s: %s (V: %d --> C: %d), (V: %d --> C: %d).\n", Vcolors.Sender, Vcolors.Recipient, Vcolors.Data.Msg, Vcolors.Data.Vertex1+1, Vcolors.Data.Color1, Vcolors.Data.Vertex2+1, Vcolors.Data.Color2)
		// Verifier checks if the colors of the edge are the same
		if Vcolors.Data.Color1 != Vcolors.Data.Color2 {
			fmt.Println("Verifier: OK. I am one step closer to being convinced that, indeed, this graph is three-colorable...\n")
		} else {
			fmt.Println("Verifier: Nope, you're not going to decieve me!\n")
		}
	} else {
		fmt.Println(Vcolors.Data.Msg)
	}

}

func main() {

	// Graph G = (V,E); V = |n| nodes -> n=10 and E = {u,v}; |m| edges
	// V = {v1, v2, ..., v10} and
	// E = {{v1,v2}, {v1,v3}, {v1,v4}, {v1,v6}, {v1,v7}, {v2,v5}, {v3,v4}, {v3,v6}, {v3,v7}, {v4,v5}, {v4,v6}, {v4,v7}, {v4,v8}, {v5,v7}, {v5,v8}, {v5,v10}, {v6,v7}, {v6,v9}, {v7,v9}, {v8,v10}, {v9,v10}}
	// Manually representing the graph in an adjacency matrix form
	// Graph G is 3-Colorable; n = 10 nodes; edge(row, col)
	fmt.Println()
	G := [][]int{{0, 1, 0, 0, 1, 1, 0, 0, 0, 1},
		{1, 0, 1, 0, 0, 1, 1, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 1, 1, 0, 0},
		{0, 0, 1, 0, 1, 0, 0, 1, 1, 0},
		{1, 0, 0, 1, 0, 0, 0, 0, 1, 1},
		{1, 1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 1, 1, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 0, 0, 0, 0}}

	// Print the graph
	for i := 0; i < len(G); i++ {
		fmt.Printf("%d", G[i])
		fmt.Println()
	}

	fmt.Println()
	// Creating two communication channels for the prover and verifier
	ProverToVerifier := make(chan Message, 10)
	VerifierToProver := make(chan Message, 10)

	wg.Add(1)
	go Verifier(G, ProverToVerifier, VerifierToProver, &wg)

	wg.Add(1)
	go Prover(G, ProverToVerifier, VerifierToProver, &wg)

	// Wait for both goroutines to finish
	wg.Wait()

	// Close the channels when done
	close(ProverToVerifier)
	close(VerifierToProver)

	fmt.Println("Protocol complete. (not really)\n")
}
