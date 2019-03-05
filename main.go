package main

import (
	"./p2"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	MPT        map[string]string `json:"mpt"`
}

func main() {
	var bj BlockJson
	bj.Height = 32
	bj.Timestamp = 123123123
	bj.Hash = "UseAbc"
	bj.ParentHash = "fight!"
	bj.Size = 10086

	m := make(map[string]string)
	m["first"] = "1"
	m["second"] = "2"

	//for key, value := range m {
	//	fmt.Println(key, ": ", value)
	//}
	bj.MPT = m

	var bjs []BlockJson = []BlockJson{bj, bj, bj}
	j, _ := json.Marshal(bjs)
	fmt.Println(string(j))
}

func TestBlockChainBasic(t *testing.T) {
	var bc BlockChain.BlockChain
	bc.Initial()
	jsonBlockChain := "[{\"hash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"timeStamp\": 1234567890, \"height\": 1, \"parentHash\": \"genesis\", \"size\": 1174, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}, {\"hash\": \"24cf2c336f02ccd526a03683b522bfca8c3c19aed8a1bed1bbc23c33cd8d1159\", \"timeStamp\": 1234567890, \"height\": 2, \"parentHash\": \"3ff3b4efe9177f705550231079c2459ba54a22d340a517e84ec5261a0d74ca48\", \"size\": 1231, \"mpt\": {\"hello\": \"world\", \"charles\": \"ge\"}}]"
	bc, err := bc.DecodeFromJSON(jsonBlockChain)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	jsonNew, err := bc.EncodeToJson()
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	var realValue []BlockJson
	var expectedValue []BlockJson
	err = json.Unmarshal([]byte(jsonNew), &realValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	err = json.Unmarshal([]byte(jsonBlockChain), &expectedValue)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	if !reflect.DeepEqual(realValue, expectedValue) {
		fmt.Println("=========Real=========")
		fmt.Println(realValue)
		fmt.Println("=========Expcected=========")
		fmt.Println(expectedValue)
		t.Fail()
	}
}

func getInt(key string) []int {
	var res []int
	for _, c := range key {
		res = append(res, int(c)-96)
	}
	return res
}
