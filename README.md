This is a server that runs Reversi and an AI bot built to play it

TODO:
  - [x] Change []string for the board to []int so that we're comparing int's which is faster
  - [] Change the validMoves function to check 4 of the adjcent squares and flip depending on if we're at P or 0
  - [] Try trimming the network as we move down to get better results
  - [x] alpha beta pruning
  - [] check the alpha beta pruning and minmax fuction to make sure they are working correctly
  - [] we need a better heuristic function