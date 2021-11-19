package Helper


type Tree struct {
	LeafHashOfBestPath  string
	Root                GenesisBlock
	BlockHashToBlock    map[string]Block
	LengthOfBestPath    int
	BlocksThatAreWaitingForTheirParent map[string]*SafeArray_Block    // Map containing:        block.prevHash -> block
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

func (t *Tree) FindFinalBlock(rollbackLimit int) (Block, bool) {
	if rollbackLimit >= t.LengthOfBestPath {
		return Block{}, false
	}
	currentBlockHash := t.LeafHashOfBestPath
	for i := 0; i < rollbackLimit; i++ {
		currentBlockHash = t.BlockHashToBlock[currentBlockHash].PrevBlockHash
	}
	return t.BlockHashToBlock[currentBlockHash], true
}

