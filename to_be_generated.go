package main

type LightBlockHeaderProof struct {

	// The index of the light block header in the vector commitment tree
	Index uint64 `json:"index"`

	// The encoded proof.
	Proof []byte `json:"proof"`

	// Represents the depth of the tree that is being proven, i.e. the number of edges from a leaf to the root.
	Treedepth uint64 `json:"treedepth"`
}
