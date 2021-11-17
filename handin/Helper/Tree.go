package Helper


type Tree struct {
	LeafHashOfBestPath  string
	Root                GenesisBlock
	BlockHashToBlock    map[string]Block
	LengthOfBestPath    int
	BlocksWhoseParentIsNotInTree map[string]Block    // Map containing:        block.prevHash -> block
}

func (t *Tree) Insert(block Block) {
	blockHash := block.Hash()

	t.BlockHashToBlock[blockHash] = block

	var lenToRoot int
	if block.PrevBlockHash == t.Root.Seed {
		lenToRoot = 1
	} else {
		lenToRoot = t.BlockHashToBlock[block.PrevBlockHash].LengthToRoot + 1
	}

	block.LengthToRoot = lenToRoot
	if (lenToRoot == t.LengthOfBestPath && blockHash > t.LeafHashOfBestPath) ||
		(lenToRoot > t.LengthOfBestPath) {

		t.LengthOfBestPath = block.LengthToRoot
		t.LeafHashOfBestPath = blockHash
	}
}




