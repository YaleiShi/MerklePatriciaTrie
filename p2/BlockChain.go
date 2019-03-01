package BlockChain

import (
	"encoding/json"
	"fmt"
)

type BlockChain struct {
	Chain  map[int32][]Block
	Length int32
}

func (bc *BlockChain) Get(height int32) []Block {
	if blocks, ok := bc.Chain[height]; ok {
		return blocks
	}
	return nil
}

func (bc *BlockChain) Insert(b Block) {
	height := b.header.height
	if blocks, ok := bc.Chain[height]; ok {
		fmt.Println("append it !!!")
		hash := b.header.hash
		for _, v := range blocks {
			if v.header.hash == hash {
				fmt.Println("there is a same!!!")
				return
			}
		}
		bc.Chain[height] = append(bc.Chain[height], b)
	} else {
		fmt.Println("new it !!!")
		bc.Chain[height] = []Block{b}
		if height > bc.Length {
			bc.Length = height
		}
	}
}

func (bc *BlockChain) EncodeToJson() (string, error) {
	chain := bc.Chain
	var jsonBlocks []BlockJson
	for _, blocks := range chain {
		for _, block := range blocks {
			jsonBlocks = append(jsonBlocks, block.ToBlockJson())
		}
	}
	res, err := json.Marshal(jsonBlocks)
	return string(res), err
}

func (bc *BlockChain) DecodeFromJSON(js string) (BlockChain, error) {
	var bjs []BlockJson
	err := json.Unmarshal([]byte(js), &bjs)
	fmt.Println("length: ", len(bjs))
	for _, bj := range bjs {
		block := BlockJSONToBlock(bj)
		bc.Insert(block)
	}
	fmt.Println("height 1: ", len(bc.Chain[1]))
	fmt.Println("height 2: ", len(bc.Chain[2]))
	fmt.Println("height 3: ", len(bc.Chain[3]))
	return *bc, err
}

func (bc *BlockChain) Initial() {
	bc.Chain = make(map[int32][]Block)
}
