package p1

import (
	"encoding/hex"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"reflect"
	"strings"
)

type Flag_value struct {
	encoded_prefix []uint8
	value          string
}

type Node struct {
	node_type    int // 0: Null, 1: Branch, 2: Ext or Leaf
	branch_value [17]string
	flag_value   Flag_value
}

type MerklePatriciaTrie struct {
	db   map[string]Node
	root string
}

func (mpt *MerklePatriciaTrie) Get(key string) (string, error) {
	// TODO
	path := getHexArray(key)
	return mpt.getHelper(path, mpt.db[mpt.root]), errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) getHelper(hexPath []uint8, node Node) string {
	if node.node_type == 0 {
		return ""
	}
	if node.node_type == 2 { // this is leaf or extend
		encodedPrefix := node.flag_value.encoded_prefix
		decode := compact_decode(encodedPrefix)
		if encodedPrefix[0]/16 >= 2 { // this is leaf
			if len(hexPath) != len(decode) {
				return ""
			} else {
				for i := 0; i < len(hexPath); i++ {
					if hexPath[i] != decode[i] {
						return ""
					}
				}
				return node.flag_value.value
			}
		} else { //this is extend
			if len(hexPath) < len(decode) {
				return ""
			}

			for i := 0; i < len(decode); i++ {
				if decode[i] != hexPath[i] {
					return ""
				}
			}
			hexPath = hexPath[len(decode):]
			nextNode := mpt.db[node.flag_value.value]
			return mpt.getHelper(hexPath, nextNode)
		}
	}

	// this is branch
	if len(hexPath) == 0 {
		return node.branch_value[16]
	}
	i := hexPath[0]
	hashValue := node.branch_value[i]
	nextNode := mpt.db[hashValue]
	hexPath = hexPath[1:]
	return mpt.getHelper(hexPath, nextNode)
}

func (mpt *MerklePatriciaTrie) Insert(key string, new_value string) {
	// TODO
	path := getHexArray(key)
	if mpt.root == "" {
		var newRoot Node
		mpt.root = mpt.insertHelper("", path, new_value, newRoot)
		return
	}

	rootNode := mpt.db[mpt.root]
	//delete(mpt.db, mpt.root)
	newRootString := mpt.insertHelper(mpt.root, path, new_value, rootNode)
	mpt.root = newRootString
}

func (mpt *MerklePatriciaTrie) insertHelper(hashValue string, path []uint8, value string, node Node) string {
	// this is a null node
	if node.node_type == 0 {
		var newHashNodeValue string
		newHashNodeValue, node = newLeaf(path, value)
		mpt.db[newHashNodeValue] = node
		return newHashNodeValue
	}
	// this is a branch node
	// this is a extension node
	// this is a leaf node
	return ""
}

func newLeaf(path []uint8, value string) (string, Node) {
	var newNode Node
	newNode.node_type = 2
	newNode.flag_value.value = value

	path = append(path, 16)
	encoded := compact_encode(path)
	newNode.flag_value.encoded_prefix = encoded

	hashNode := newNode.hash_node()
	return hashNode, newNode
}

func (mpt *MerklePatriciaTrie) Delete(key string) (string, error) {
	// TODO
	return "", errors.New("path_not_found")
}

func (mpt *MerklePatriciaTrie) Test(key string) {
	fmt.Println(mpt.db["not exist"].node_type)
	//hex1 := getHexArray(key)
	//fmt.Println(hex1)
	//
	//hex2 := compact_encode(hex1)
	//fmt.Println(hex2)

	//hex3 := compact_encode([]uint8{1,1,2,3,4,5,})
	//fmt.Println(hex3)
	//test_compact_encode()
}

func getHexArray(key string) []uint8 {
	var res []uint8
	for _, c := range key {
		i := uint8(c)
		first := i / 16
		last := i % 16
		res = append(res, first, last)
	}
	return res
}

func compact_encode(hex_array []uint8) []uint8 {
	// TODO
	term := 0
	if hex_array[len(hex_array)-1] == 16 {
		term = 1
	}
	if term == 1 {
		hex_array = hex_array[:len(hex_array)-1]
	}
	oddlen := len(hex_array) % 2
	flags := 2*term + oddlen
	if oddlen == 1 {
		hex_array = append([]uint8{uint8(flags)}, hex_array...)
	} else {
		hex_array = append([]uint8{uint8(flags), 0}, hex_array...)
	}
	fmt.Println(hex_array)
	var res []uint8
	for i := 0; i < len(hex_array); i += 2 {
		res = append(res, 16*hex_array[i]+hex_array[i+1])
	}
	return res
}

// If Leaf, ignore 16 at the end
func compact_decode(encoded_arr []uint8) []uint8 {
	// TODO
	var res []uint8
	flag := encoded_arr[0]
	oddlen := flag % 2
	for _, i := range encoded_arr {
		first := i / 16
		last := i % 16
		res = append(res, first, last)
	}

	if oddlen == 0 {
		res = res[2:]
	} else {
		res = res[1:]
	}
	return res
}

func test_compact_encode() {
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{1, 2, 3, 4, 5})), []uint8{1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 1, 2, 3, 4, 5})), []uint8{0, 1, 2, 3, 4, 5}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{0, 15, 1, 12, 11, 8, 16})), []uint8{0, 15, 1, 12, 11, 8}))
	fmt.Println(reflect.DeepEqual(compact_decode(compact_encode([]uint8{15, 1, 12, 11, 8, 16})), []uint8{15, 1, 12, 11, 8}))
}

func (node *Node) hash_node() string {
	var str string
	switch node.node_type {
	case 0:
		str = ""
	case 1:
		str = "branch_"
		for _, v := range node.branch_value {
			str += v
		}
	case 2:
		str = node.flag_value.value
	}

	sum := sha3.Sum256([]byte(str))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}

func (node *Node) String() string {
	str := "empty string"
	switch node.node_type {
	case 0:
		str = "[Null Node]"
	case 1:
		str = "Branch["
		for i, v := range node.branch_value[:16] {
			str += fmt.Sprintf("%d=\"%s\", ", i, v)
		}
		str += fmt.Sprintf("value=%s]", node.branch_value[16])
	case 2:
		encoded_prefix := node.flag_value.encoded_prefix
		node_name := "Leaf"
		if is_ext_node(encoded_prefix) {
			node_name = "Ext"
		}
		ori_prefix := strings.Replace(fmt.Sprint(compact_decode(encoded_prefix)), " ", ", ", -1)
		str = fmt.Sprintf("%s<%v, value=\"%s\">", node_name, ori_prefix, node.flag_value.value)
	}
	return str
}

func node_to_string(node Node) string {
	return node.String()
}

func (mpt *MerklePatriciaTrie) Initial() {
	mpt.db = make(map[string]Node)
}

func is_ext_node(encoded_arr []uint8) bool {
	return encoded_arr[0]/16 < 2
}

func TestCompact() {
	test_compact_encode()
}

func (mpt *MerklePatriciaTrie) String() string {
	content := fmt.Sprintf("ROOT=%s\n", mpt.root)
	for hash := range mpt.db {
		content += fmt.Sprintf("%s: %s\n", hash, node_to_string(mpt.db[hash]))
	}
	return content
}

func (mpt *MerklePatriciaTrie) Order_nodes() string {
	raw_content := mpt.String()
	content := strings.Split(raw_content, "\n")
	root_hash := strings.Split(strings.Split(content[0], "HashStart")[1], "HashEnd")[0]
	queue := []string{root_hash}
	i := -1
	rs := ""
	cur_hash := ""
	for len(queue) != 0 {
		last_index := len(queue) - 1
		cur_hash, queue = queue[last_index], queue[:last_index]
		i += 1
		line := ""
		for _, each := range content {
			if strings.HasPrefix(each, "HashStart"+cur_hash+"HashEnd") {
				line = strings.Split(each, "HashEnd: ")[1]
				rs += each + "\n"
				rs = strings.Replace(rs, "HashStart"+cur_hash+"HashEnd", fmt.Sprintf("Hash%v", i), -1)
			}
		}
		temp2 := strings.Split(line, "HashStart")
		flag := true
		for _, each := range temp2 {
			if flag {
				flag = false
				continue
			}
			queue = append(queue, strings.Split(each, "HashEnd")[0])
		}
	}
	return rs
}
