package Helper

import "sync"

type Tree struct {
	LeafHashOfBestPath  string
	Root                GenesisBlock
	BlockHashToBlock    SafeMap_string_to_Block
	LengthOfBestPath    int
	BlocksThatAreWaitingForTheirParent map[string]*SafeArray_Block    // Map containing:        block.prevHash -> block
	WaitingForParentMapLock sync.Mutex
}

func (t *Tree) Insert(block Block) {
	var lenToRoot int
	if block.PrevBlockHash == t.Root.Seed {
		lenToRoot = 1
	} else {
		prevBlock, _ := t.BlockHashToBlock.Get(block.PrevBlockHash)
		lenToRoot = prevBlock.LengthToRoot + 1
		//lenToRoot = t.BlockHashToBlock[block.PrevBlockHash].LengthToRoot + 1
	}

	block.LengthToRoot = lenToRoot
	blockHash := block.Hash()
	t.BlockHashToBlock.Put(blockHash, block)
	//t.BlockHashToBlock[blockHash] = block

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
		currentBlock, _ := t.BlockHashToBlock.Get(currentBlockHash)
		currentBlockHash = currentBlock.PrevBlockHash
	}
	block, ok := t.BlockHashToBlock.Get(currentBlockHash)
	if block.HasBeenApplied || !ok {
		return Block{}, false
	}
	return block, true
	//return t.BlockHashToBlock[currentBlockHash], true
}

