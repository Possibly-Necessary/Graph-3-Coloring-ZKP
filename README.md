# Zero Knowldege Interactive Proof for Graph 3-Colorability in Go

A (horrible/hideous, messy and incomplete) Go program that attempts to visualize/simulate the concept of Graph 3-Coloring Zero-Knowledge Interactive Proof (G3C ZKIP) discussed in this paper:

  - "Proofs the Yeild Nothing but their Validity, or All languages in NP have Zero-Knowledge Proofs" by Goldriech, Micali & Wigderson.
    link: https://dl.acm.org/doi/10.1145/116825.116852

Missing parts include: the commitment scheme, the protocol should be done in m^2 rounds where m is the number of edges in the graph G, etc.

Reference book(s):

  - "A Guide to Graph Colouring Algorithms and Applications" by R.M.R. Lewis
  - "Mastering Concurrency in Go" by Nathan Kozyra

Testing the program: the following graph, represented in an adjacency matrix, is the common input for both the prover and verifier. This particular graph is three-colorable.

 ![Graph-COLOR(1)](https://github.com/Possibly-Necessary/Graph-3-Coloring-ZKP/assets/109365947/224efc8c-6649-4fd0-bbe8-b2b427f6e316.jpg) 

                      [0 1 0 0 1 1 0 0 0 1]
                      [1 0 1 0 0 1 1 0 0 0]
                      [0 1 0 1 0 0 1 1 0 0]
                      [0 0 1 0 1 0 0 1 1 0]
                      [1 0 0 1 0 0 0 0 1 1]
                      [1 1 0 0 0 0 0 0 0 0]
                      [0 1 1 0 0 0 0 0 0 0]
                      [0 0 1 1 0 0 0 0 0 0]
                      [0 0 0 1 1 0 0 0 0 0]
                      [1 0 0 0 1 0 0 0 0 0]

The prover and verifier are implemented as Go-routine functions to communicate concurrently through a channel. The prover function invokes the 3-graph coloring (implemented using a backtracking algorithm) and colors the graph. The algorithm generates a certificate mapping of the vertices as keys and their respective colors as values:

                    map[1:1 2:3 3:1 4:3 5:2 6:2 7:2 8:2 9:1 10:3] 

Note the function that solves the 3-graph coloring generates a distinct (permuted) coloring for the same graph each time it's invoked. For instance, as apposed to the above output, rerunning the program again yields the following vertex coloring:

                    map[1:2 2:3 3:2 4:3 5:1 6:1 7:1 8:1 9:2 10:3]

Here's another instance:

                    map[1:3 2:1 3:3 4:1 5:2 6:2 7:2 8:2 9:3 10:1] 

The coloring permutation is required for the prover to alter the color of the vertices each time the verifier asks for a randomly chosen edge. Below is the case where the edge chosen by the verifier does not belong to the set of the graph's edges, which the prover will verify upon receiving it via the channel:

![Not -in-edge](https://github.com/Possibly-Necessary/Graph-3-Coloring-ZKP/assets/109365947/b359c8e4-dc69-4be4-a3e7-5fa97cae4665.jpg)

Below is the alternative case, where the prover validates the verifier's edge to be a part of the graph's edge set and "reveals", or simply, sends forth the coloring of the edge's endpoints to the verifier. The verifier's next step is to check if the adjacent vertex colors (connected by his selected edge) match:

![In Edge](https://github.com/Possibly-Necessary/Graph-3-Coloring-ZKP/assets/109365947/0f924310-844a-4b2a-86f3-765a68e8946a.jpg)





