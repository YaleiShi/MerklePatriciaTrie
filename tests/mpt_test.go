package tests

import (
	"../p2"
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

func TestBlockChainBasic(t *testing.T) {
	var bc BlockChain.BlockChain
	bc.Initial()
	jsonBlockChain := "[{\"height\":1,\"timeStamp\":1551025401,\"hash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"parentHash\":\"genesis\",\"size\":2089,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}},{\"height\":2,\"timeStamp\":1551025401,\"hash\":\"944eb943b05caba08e89a613097ac5ac7d373d863224d17b1958541088dc20e2\",\"parentHash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"size\":2146,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}},{\"height\":2,\"timeStamp\":1551025401,\"hash\":\"f8af68feadf25a635bc6e81c08f81c6740bbe1fb2514c1b4c56fe1d957c7448d\",\"parentHash\":\"6c9aad47a370269746f172a464fa6745fb3891194da65e3ad508ccc79e9a771b\",\"size\":707,\"mpt\":{\"ge\":\"Charles\"}},{\"height\":3,\"timeStamp\":1551025401,\"hash\":\"f367b7f59c651e69be7e756298aad62fb82fddbfeda26cb06bfd8adf9c8aa094\",\"parentHash\":\"f8af68feadf25a635bc6e81c08f81c6740bbe1fb2514c1b4c56fe1d957c7448d\",\"size\":707,\"mpt\":{\"ge\":\"Charles\"}},{\"height\":3,\"timeStamp\":1551025401,\"hash\":\"05ac44dd82b6cc398a5e9664add21856ae19d107d9035af5fc54c9b0ffdef336\",\"parentHash\":\"944eb943b05caba08e89a613097ac5ac7d373d863224d17b1958541088dc20e2\",\"size\":2146,\"mpt\":{\"CS686\":\"BlockChain\",\"test1\":\"value1\",\"test2\":\"value2\",\"test3\":\"value3\",\"test4\":\"value4\"}}]"
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

//func TestExt(t *testing.T) {
//	mpt := new(p1.MerklePatriciaTrie)
//	mpt.Initial()
//	mpt.Insert("p", "apple")
//	mpt.Insert("aa", "banana")
//	mpt.Insert("ap", "orange")
//	inserted_trie := mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestExt 011", mpt.Order_nodes(), "./mpt_tests/ext_011.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 011", mpt.Order_nodes(), "./mpt_tests/ext_011.txt", t)
//	mpt.Delete("b")
//	check_eq("TestExt 011", mpt.Order_nodes(), inserted_trie, t)
//	v, _ := mpt.Get("aa")
//	check_eq("TestExt 011", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("p", "apple")
//	mpt.Insert("aa", "banana")
//	mpt.Insert("ap", "orange")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("ba", "new")
//	check_mpt("TestExt 013", mpt.Order_nodes(), "./mpt_tests/ext_013.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 013", mpt.Order_nodes(), "./mpt_tests/ext_013.txt", t)
//	mpt.Delete("ba")
//	check_eq("TestExt 013", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("aa")
//	check_eq("TestExt 013", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("aaa", "apple")
//	mpt.Insert("aap", "banana")
//	mpt.Insert("bb", "right leaf")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("aa", "new")
//	check_mpt("TestExt 030", mpt.Order_nodes(), "./mpt_tests/ext_030.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 030", mpt.Order_nodes(), "./mpt_tests/ext_030.txt", t)
//	mpt.Delete("aa")
//	check_eq("TestExt 030", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("aaa")
//	check_eq("TestExt 030", v, "apple", t)
//
//	mpt.Initial()
//	mpt.Insert("p", "apple")
//	mpt.Insert("aaa", "banana")
//	mpt.Insert("aap", "orange")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestExt 031", mpt.Order_nodes(), "./mpt_tests/ext_031.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 031", mpt.Order_nodes(), "./mpt_tests/ext_031.txt", t)
//	mpt.Delete("b")
//	check_eq("TestExt 031", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("aaa")
//	check_eq("TestExt 031", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("p", "apple")
//	mpt.Insert("aaa", "banana")
//	mpt.Insert("aap", "orange")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("ba", "new")
//	check_mpt("TestExt 033", mpt.Order_nodes(), "./mpt_tests/ext_033.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 033", mpt.Order_nodes(), "./mpt_tests/ext_033.txt", t)
//	mpt.Delete("ba")
//	check_eq("TestExt 033", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("aaa")
//	check_eq("TestExt 033", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ap", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestExt 111", mpt.Order_nodes(), "./mpt_tests/ext_111.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 111", mpt.Order_nodes(), "./mpt_tests/ext_111.txt", t)
//	mpt.Delete("b")
//	check_eq("TestExt 111", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("ap")
//	check_eq("TestExt 111", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ap", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("bc", "new")
//	check_mpt("TestExt 113", mpt.Order_nodes(), "./mpt_tests/ext_113.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 113", mpt.Order_nodes(), "./mpt_tests/ext_113.txt", t)
//	mpt.Delete("bc")
//	check_eq("TestExt 113", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("ap")
//	check_eq("TestExt 113", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("p", "apple")
//	mpt.Insert("aaaa", "banana")
//	mpt.Insert("aaap", "orange")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("a", "new")
//	check_mpt("TestExt 140", mpt.Order_nodes(), "./mpt_tests/ext_140.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 140", mpt.Order_nodes(), "./mpt_tests/ext_140.txt", t)
//	mpt.Delete("a")
//	check_eq("TestExt 140", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("aaaa")
//	check_eq("TestExt 140", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("aaa", "apple")
//	mpt.Insert("aap", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestExt 131", mpt.Order_nodes(), "./mpt_tests/ext_131.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 131", mpt.Order_nodes(), "./mpt_tests/ext_131.txt", t)
//	mpt.Delete("b")
//	check_eq("TestExt 131", mpt.Order_nodes(), inserted_trie, t)
//	v, _ = mpt.Get("aap")
//	check_eq("TestExt 131", v, "banana", t)
//
//	mpt.Initial()
//	mpt.Insert("aaa", "apple")
//	mpt.Insert("aap", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("bc", "new")
//	check_mpt("TestExt 133", mpt.Order_nodes(), "./mpt_tests/ext_133.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestExt 133", mpt.Order_nodes(), "./mpt_tests/ext_133.txt", t)
//}
//
//func TestLeaf(t *testing.T) {
//	mpt := new(p1.MerklePatriciaTrie)
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	mpt.Insert("a", "new")
//	check_mpt("TestLeaf 000", mpt.Order_nodes(), "./mpt_tests/leaf_000.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 000", mpt.Order_nodes(), "./mpt_tests/leaf_000.txt", t)
//	mpt.Delete("a")
//	check_mpt("TestLeaf 000", mpt.Order_nodes(), "./mpt_tests/delete_basic_1.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	mpt.Insert("ab", "new")
//	check_mpt("TestLeaf 002", mpt.Order_nodes(), "./mpt_tests/leaf_002.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 002", mpt.Order_nodes(), "./mpt_tests/leaf_002.txt", t)
//	mpt.Delete("ab")
//	check_mpt("TestLeaf 002", mpt.Order_nodes(), "./mpt_tests/basic_0.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("p", "banana")
//	inserted_trie := mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestLeaf 011", mpt.Order_nodes(), "./mpt_tests/leaf_011.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 011", mpt.Order_nodes(), "./mpt_tests/leaf_011.txt", t)
//	mpt.Delete("b")
//	check_eq("TestLeaf 011", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("p", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("bc", "new")
//	check_mpt("TestLeaf 013", mpt.Order_nodes(), "./mpt_tests/leaf_013.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 013", mpt.Order_nodes(), "./mpt_tests/leaf_013.txt", t)
//	mpt.Delete("bc")
//	check_eq("TestLeaf 013", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("bab", "apple")
//	mpt.Insert("aa", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestLeaf 040", mpt.Order_nodes(), "./mpt_tests/leaf_040.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 040", mpt.Order_nodes(), "./mpt_tests/leaf_040.txt", t)
//	mpt.Delete("b")
//	check_eq("TestLeaf 040", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("aab", "apple")
//	mpt.Insert("app", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("ac", "new")
//	check_mpt("TestLeaf 031", mpt.Order_nodes(), "./mpt_tests/leaf_031.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 031", mpt.Order_nodes(), "./mpt_tests/leaf_031.txt", t)
//	mpt.Delete("ac")
//	check_eq("TestLeaf 031", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("aab", "apple")
//	mpt.Insert("app", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("ace", "new")
//	check_mpt("TestLeaf 033", mpt.Order_nodes(), "./mpt_tests/leaf_033.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 033", mpt.Order_nodes(), "./mpt_tests/leaf_033.txt", t)
//	mpt.Delete("ace")
//	check_eq("TestLeaf 033", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("p", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("a", "apple")
//	mpt.Insert("a", "new")
//	check_mpt("TestLeaf 100", mpt.Order_nodes(), "./mpt_tests/leaf_100.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 100", mpt.Order_nodes(), "./mpt_tests/leaf_100.txt", t)
//	mpt.Delete("a")
//	check_eq("TestLeaf 100", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("p", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("abc", "new")
//	check_mpt("TestLeaf 104", mpt.Order_nodes(), "./mpt_tests/leaf_104.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 104", mpt.Order_nodes(), "./mpt_tests/leaf_104.txt", t)
//	mpt.Delete("abc")
//	check_eq("TestLeaf 104", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestLeaf 111", mpt.Order_nodes(), "./mpt_tests/leaf_111.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 111", mpt.Order_nodes(), "./mpt_tests/leaf_111.txt", t)
//	mpt.Delete("b")
//	check_eq("TestLeaf 111", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("bc", "new")
//	check_mpt("TestLeaf 113", mpt.Order_nodes(), "./mpt_tests/leaf_113.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 113", mpt.Order_nodes(), "./mpt_tests/leaf_113.txt", t)
//	mpt.Delete("bc")
//	check_eq("TestLeaf 113", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("ap", "apple")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("b", "new")
//	check_mpt("TestLeaf 131", mpt.Order_nodes(), "./mpt_tests/leaf_131.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 131", mpt.Order_nodes(), "./mpt_tests/leaf_131.txt", t)
//	mpt.Delete("b")
//	check_eq("TestLeaf 131", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("ap", "apple")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("bp", "new")
//	check_mpt("TestLeaf 133", mpt.Order_nodes(), "./mpt_tests/leaf_133.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestLeaf 133", mpt.Order_nodes(), "./mpt_tests/leaf_133.txt", t)
//	mpt.Delete("bp")
//	check_eq("TestLeaf 133", mpt.Order_nodes(), inserted_trie, t)
//}
//
//func TestDeleteBasic(t *testing.T) {
//	mpt := new(p1.MerklePatriciaTrie)
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	mpt.Delete("a")
//	check_mpt("TestDeleteBasic 1", mpt.Order_nodes(), "./mpt_tests/delete_basic_1.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	mpt.Delete("c")
//	check_mpt("TestDeleteBasic 0", mpt.Order_nodes(), "./mpt_tests/delete_basic_0.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("abb", "banana")
//	mpt.Delete("a")
//	check_mpt("TestDeleteBasic 2", mpt.Order_nodes(), "./mpt_tests/delete_basic_2.txt", t)
//}
//
//func TestBranch(t *testing.T) {
//	mpt := new(p1.MerklePatriciaTrie)
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ap", "banana")
//	inserted_trie := mpt.Order_nodes()
//	mpt.Insert("a", "new")
//	check_mpt("TestBranch nv_np", mpt.Order_nodes(), "./mpt_tests/branch_nv_np.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestBranch nv_np", mpt.Order_nodes(), "./mpt_tests/branch_nv_np.txt", t)
//	mpt.Delete("a")
//	check_eq("TestBranch nv_np", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("a", "old")
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ap", "banana")
//	mpt.Insert("a", "new")
//	check_mpt("TestBranch v_np", mpt.Order_nodes(), "./mpt_tests/branch_v_np.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestBranch v_np", mpt.Order_nodes(), "./mpt_tests/branch_v_np.txt", t)
//	mpt.Delete("a")
//	expected_mpt := new(p1.MerklePatriciaTrie)
//	expected_mpt.Initial()
//	expected_mpt.Insert("aa", "apple")
//	expected_mpt.Insert("ap", "banana")
//	check_eq("TestBranch v_np", mpt.Order_nodes(), expected_mpt.Order_nodes(), t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("c", "new")
//	check_mpt("TestBranch nv_p", mpt.Order_nodes(), "./mpt_tests/branch_nv_p.txt", t)
//	mpt.Delete("cc")
//	check_mpt("TestBranch nv_p", mpt.Order_nodes(), "./mpt_tests/branch_nv_p.txt", t)
//	mpt.Delete("c")
//	check_eq("TestBranch nv_p", mpt.Order_nodes(), inserted_trie, t)
//
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ap", "banana")
//	mpt.Insert("a", "old")
//	inserted_trie = mpt.Order_nodes()
//	mpt.Insert("aA", "new")
//	check_mpt("TestBranch v_p", mpt.Order_nodes(), "./mpt_tests/branch_v_p.txt", t)
//	mpt.Delete("c")
//	check_mpt("TestBranch v_p", mpt.Order_nodes(), "./mpt_tests/branch_v_p.txt", t)
//	mpt.Delete("aA")
//	check_eq("TestBranch v_p", mpt.Order_nodes(), inserted_trie, t)
//}
//
//func TestLeafBasic(t *testing.T) {
//	mpt := new(p1.MerklePatriciaTrie)
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	check_mpt("TestLeafBasic 0", mpt.Order_nodes(), "./mpt_tests/basic_0.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("p", "banana")
//	check_mpt("TestLeafBasic 1", mpt.Order_nodes(), "./mpt_tests/basic_1.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	check_mpt("TestLeafBasic 2", mpt.Order_nodes(), "./mpt_tests/basic_2.txt", t)
//}
//
//func TestExtBasic(t *testing.T) {
//	mpt := new(p1.MerklePatriciaTrie)
//	mpt.Initial()
//	mpt.Insert("a", "apple")
//	mpt.Insert("b", "banana")
//	check_mpt("TestExtBasic 1", mpt.Order_nodes(), "./mpt_tests/ext_basic_1.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ap", "banana")
//	check_mpt("TestExtBasic 2", mpt.Order_nodes(), "./mpt_tests/ext_basic_2.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("aa", "apple")
//	mpt.Insert("ab", "banana")
//	check_mpt("TestExtBasic 3", mpt.Order_nodes(), "./mpt_tests/ext_basic_3.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("aaa", "apple")
//	mpt.Insert("aap", "banana")
//	check_mpt("TestExtBasic 4", mpt.Order_nodes(), "./mpt_tests/ext_basic_4.txt", t)
//
//	mpt.Initial()
//	mpt.Insert("p", "apple")
//	mpt.Insert("aa", "banana")
//	mpt.Insert("ap", "orange")
//	check_mpt("TestExtBasic 5", mpt.Order_nodes(), "./mpt_tests/ext_basic_5.txt", t)
//}
//
//func check_mpt(id string, mpt_string string, file_path string, t *testing.T) {
//	content, _ := ioutil.ReadFile(file_path)
//	if string(content) != mpt_string {
//		for i := 0; i < len(mpt_string); i++ {
//			if mpt_string[i] == string(content)[i] {
//				fmt.Print(string(mpt_string[i]))
//			}
//			if mpt_string[i] != string(content)[i] {
//				fmt.Print("failed piece: *", mpt_string[i], "*")
//				fmt.Print("right piece: &", string(content)[i], "&")
//				break
//			}
//		}
//		fmt.Println()
//		fmt.Println("=========" + id + "============")
//		fmt.Println("=========Real============")
//		fmt.Println(mpt_string)
//		fmt.Println("=========Expcected============")
//		fmt.Println(string(content))
//		fmt.Println("=====================")
//		t.Fail()
//	}
//}
//
//func check_eq(id string, real string, expected string, t *testing.T) {
//	if real != expected {
//		for i := 0; i < len(real); i++ {
//			if real[i] == expected[i] {
//				fmt.Print(real[i])
//			}
//		}
//		fmt.Println("=========" + id + "============")
//		fmt.Println("=========Real============")
//		fmt.Println(real)
//		fmt.Println("=========Expcected============")
//		fmt.Println(expected)
//		fmt.Println("=====================")
//		t.Fail()
//	}
//}
//
//func TestCompactEncode(t *testing.T) {
//	p1.TestCompact()
//}
