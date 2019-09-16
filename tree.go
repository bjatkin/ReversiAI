package main

type BoardRoot struct {
	OldLeaves []*BoardLeaf
	NewLeaves []*BoardLeaf
	Board     []string
}

type BoardLeaf struct {
	Board   []string
	Utility float64
}
