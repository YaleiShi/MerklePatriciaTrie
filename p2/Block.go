package BlockChain

import (
	"../p1"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"time"
)

type Block struct {
	header Header
	value  p1.MerklePatriciaTrie `json:"value"`
}

type Header struct {
	height     int32
	timeStamp  int64
	hash       string
	parentHash string
	size       int32
}

type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	MPT        map[string]string `json:"mpt"`
}

func (b *Block) Initial(height int32, parentHash string, value p1.MerklePatriciaTrie) {
	b.header.height = height
	b.header.timeStamp = time.Now().Unix()
	b.header.parentHash = parentHash
	b.value = value

	b.header.size = getMPTLength(value)
	str := string(b.header.height) + string(b.header.timeStamp) + b.header.parentHash + b.value.GetRoot() + string(b.header.size)
	b.header.hash = HashBlock(str)
}

func HashBlock(str string) string {
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

func DecodeFromJson(js string) Block {
	var bj BlockJson
	json.Unmarshal([]byte(js), &bj)
	return BlockJSONToBlock(bj)
}

func BlockJSONToBlock(bj BlockJson) Block {
	var block Block
	var mpt p1.MerklePatriciaTrie
	mpt.Initial()
	mptKV := bj.MPT
	for key, value := range mptKV {
		mpt.Insert(key, value)
	}

	block.value = mpt
	block.header.hash = bj.Hash
	block.header.size = bj.Size
	block.header.timeStamp = bj.Timestamp
	block.header.parentHash = bj.ParentHash
	block.header.height = bj.Height
	return block
}

func (b *Block) EncodeToJSON() string {
	bj := b.ToBlockJson()
	j, _ := json.Marshal(bj)
	return string(j)
}

func (b *Block) ToBlockJson() BlockJson {
	var bj BlockJson
	bj.Height = b.header.height
	bj.Timestamp = b.header.timeStamp
	bj.Hash = b.header.hash
	bj.ParentHash = b.header.parentHash
	bj.Size = b.header.size
	bj.MPT = b.value.GetLeafMap()
	return bj
}

func getMPTLength(data p1.MerklePatriciaTrie) int32 {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return 0
	}
	return int32(len(buf.Bytes()))
}
