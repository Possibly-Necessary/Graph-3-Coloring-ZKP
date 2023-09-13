# Zero Knowldege Interactive Proof for Graph 3-Colorability in Go

A (horrible/hideous, messy and incomplete) Go program that attempts to visulize/simulate the concept of Graph 3-Coloring Zero-Knowledge Interactive Proof (G3C ZKIP) discussed in this paper:

  - "Proofs the Yeild Nothing but their Validity, or All languages in NP have Zero-Knowledge Proofs" by Goldriech, Micali & Wigderson.
    link: https://dl.acm.org/doi/10.1145/116825.116852

Missing parts include: the commitment scheme, the protocol should be done in m^2 rounds where m is the number of edges in the graph G, etc.

Reference book(s):

  - "A Guide to Graph Colouring Algorithms and Applications" by R.M.R. Lewis
  - "Mastering Concurrency in Go" by Nathan Kozyra

Testing the program: the following graph, represented in an adjacency matrix, is the common input for both the prover and verifier. This particular graph is three-colorable.
![Graph-COLOR](https://github.com/Possibly-Necessary/Graph-3-Coloring-ZKP/assets/109365947/7b7ae6bf-060e-4fa8-bc40-f8ed1c455a27)

        
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

The prover and verifier are implemented as Go-routine functions to communicate with eachother (through a channel) concurrently. The prover function will invoke the 3-graph coloring, implemented using a backtracking algorithm, and colors the graph. The algorithm returns a certificate mapping of the vertices as the keys and their colors as the values:

                    map[1:1 2:3 3:1 4:3 5:2 6:2 7:2 8:2 9:1 10:3] 

The function that solves the 3-graph coloring returns a different (permuted) coloring for the same graph each time it's invoked. For instance, compared to the above output, running the program again yields the following coloring for the vertices:

                    map[1:2 2:3 3:2 4:3 5:1 6:1 7:1 8:1 9:2 10:3]

Another example:

                    map[1:3 2:1 3:3 4:1 5:2 6:2 7:2 8:2 9:3 10:1] 

The coloring permutation is required for the prover to permute the color of the vertices each time the verifier asks for a randomly chosen edge, which happens next. Below is the case if the edge selected by the verifier does not belong to the set of edges in the graph which will be checked by the prover after it receives it through the channel:

