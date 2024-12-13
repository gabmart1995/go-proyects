package main

type BlockChain struct {
	blocks []*Block
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)

	bc.blocks = append(bc.blocks, newBlock)
}

// crea una cadena de bloques implementando un
// bloque genesis
func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*Block{NewGenesisBlock()},
	}
}
