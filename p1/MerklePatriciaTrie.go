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
	//fmt.Println("searching path: ", node.String())
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

	if node.node_type == 2 {
		encodedPrefix := node.flag_value.encoded_prefix
		decode := compact_decode(encodedPrefix)
		n := len(decode)
		lenPath := len(path)
		var common int
		for common = 0; common < n && common < lenPath; common++ {
			if path[common] != decode[common] {
				break
			}
		}

		if encodedPrefix[0]/16 < 2 { // extension
			// this is a extension node
			nextHash := node.flag_value.value
			nextNode := mpt.db[nextHash]
			if common == n {
				// common = n go down
				path = path[common:]

				newNextHash := mpt.insertHelper(nextHash, path, value, nextNode)
				node.flag_value.value = newNextHash

				newHash := node.hash_node()
				delete(mpt.db, hashValue)
				mpt.db[newHash] = node
				return newHash
			} else if common == 0 {
				// common = 0 switch to branch
				node.node_type = 1
				node.flag_value.value = ""
				node.flag_value.encoded_prefix = []uint8{}

				if len(decode) == 1 {
					node.branch_value[decode[0]] = nextHash
				} else {
					i := decode[0]
					decode = decode[1:]
					var newExNode Node
					newExNode.node_type = 2
					newExNode.flag_value.encoded_prefix = compact_encode(decode)
					newExNode.flag_value.value = nextHash
					newExNodeHash := newExNode.hash_node()
					mpt.db[newExNodeHash] = newExNode

					node.branch_value[i] = newExNodeHash
				}
				newHash := mpt.insertHelper(hashValue, path, value, node)
				return newHash
			} else {
				// else common = prefix , rest new branch
				leftDecode := decode[:common]
				node.flag_value.encoded_prefix = compact_encode(leftDecode)

				var newBranchNode Node
				newBranchNode.node_type = 1
				rightDecode := decode[common:]

				if len(rightDecode) == 1 {
					newBranchNode.branch_value[rightDecode[0]] = nextHash
				} else {
					var newExNode Node
					newExNode.node_type = 2
					prefix := rightDecode[1:]
					newExNode.flag_value.encoded_prefix = compact_encode(prefix)
					newExNode.flag_value.value = nextHash
					newExNodeHash := newExNode.hash_node()
					mpt.db[newExNodeHash] = newExNode

					newBranchNode.branch_value[rightDecode[0]] = newExNodeHash
				}

				path = path[common:]
				newBrNodeHash := mpt.insertHelper("", path, value, newBranchNode)
				node.flag_value.value = newBrNodeHash
				//delete(mpt.db, hashValue)
				newHash := node.hash_node()
				mpt.db[newHash] = node
				return newHash
			}

			//
		} else { // leaf
			if lenPath == n && common == n { //completely same
				node.flag_value.value = value
				return mpt.updateNode(node, hashValue)
			}
			// not same, first get the common hex and the rest of decode, rest of path
			prefix := decode[:common]
			restDecode := decode[common:]
			restPath := path[common:]
			oldValue := node.flag_value.value

			// make a new branch node, insert the old value and new value to the branch node
			var branchNode Node
			branchNode.node_type = 1
			// update the branch to add the old value
			if len(restDecode) == 0 {
				branchNode.branch_value[16] = oldValue
			} else {
				first := restDecode[0]
				restDecode = restDecode[1:]

				newLeafHash, newLeaf := newLeaf(restDecode, oldValue)
				mpt.db[newLeafHash] = newLeaf
				branchNode.branch_value[first] = newLeafHash

			}

			if len(restPath) == 0 {
				branchNode.branch_value[16] = value
			} else {
				first := restPath[0]
				restPath = restPath[1:]

				newLeafHash, newLeaf := newLeaf(restPath, value)
				mpt.db[newLeafHash] = newLeaf
				branchNode.branch_value[first] = newLeafHash

			}

			if common == 0 {
				// directly change the node to a branch node, insert the old value and new value to the branch node
				node = branchNode
				return mpt.updateNode(node, "")
			}
			// change the node to a extension node
			node.flag_value.encoded_prefix = compact_encode(prefix)
			node.flag_value.value = branchNode.hash_node()
			mpt.db[branchNode.hash_node()] = branchNode
			return mpt.updateNode(node, "")
		}
	}

	// this is a branch node
	if len(path) == 0 {
		node.branch_value[16] = value
		return mpt.updateNode(node, hashValue)
	}
	firstKey := path[0]
	restKey := path[1:]
	if node.branch_value[firstKey] == "" {
		var newLeaf Node
		newLeafHash := mpt.insertHelper("", restKey, value, newLeaf)
		node.branch_value[firstKey] = newLeafHash
		return mpt.updateNode(node, hashValue)
	}
	nextHash := node.branch_value[firstKey]
	nextNode := mpt.db[nextHash]
	newNextHash := mpt.insertHelper(nextHash, restKey, value, nextNode)
	node.branch_value[firstKey] = newNextHash
	return mpt.updateNode(node, hashValue)
}

func (mpt *MerklePatriciaTrie) updateNode(node Node, oldHash string) string {
	//if oldHash != "" {
	//	fmt.Println(oldHash)
	//	delete(mpt.db, oldHash)
	//}
	newHash := node.hash_node()
	mpt.db[newHash] = node

	return newHash
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

func (mpt *MerklePatriciaTrie) Delete(key string) string {
	// TODO
	path := getHexArray(key)
	if mpt.root == "" {
		return "path_not_found"
	}
	rootNode := mpt.db[mpt.root]
	var status int
	mpt.root, status = mpt.deleteHelper(mpt.root, path, rootNode)
	if status == -1 {
		return "path_not_found"
	}
	return ""
}

func (mpt *MerklePatriciaTrie) deleteHelper(oldHash string, path []uint8, node Node) (newHash string, nodeType int) {
	// node type = 0
	if node.node_type == 0 {
		return "", -1
	}

	if node.node_type == 2 {
		encodedPrefix := node.flag_value.encoded_prefix
		decode := compact_decode(encodedPrefix)
		n := len(decode)
		lenPath := len(path)
		var com int
		for com = 0; com < n && com < lenPath; com++ {
			if path[com] != decode[com] {
				break
			}
		}

		if encodedPrefix[0]/16 < 2 {
			// extension node
			if com != n { // end find, report error
				return oldHash, -1
			} else { // can go down
				nextPath := path[com:]
				nextHash := node.flag_value.value
				nextNode := mpt.db[nextHash]
				newNextHash, nextNodeType := mpt.deleteHelper(nextHash, nextPath, nextNode)
				nextNode = mpt.db[newNextHash]

				if nextNodeType == -1 { // still not found below
					return oldHash, -1
				}
				if nextNodeType == 1 { // next is a branch node, OK
					node.flag_value.value = newNextHash
					return mpt.updateNode(node, oldHash), 3
				}
				if nextNodeType == 2 || nextNodeType == 3 { // next is leaf node or extension node, need to be changed
					nextLeafValue := nextNode.flag_value.value
					nextLeafDecode := compact_decode(nextNode.flag_value.encoded_prefix)
					delete(mpt.db, newNextHash)

					node.flag_value.value = nextLeafValue
					decode = append(decode, nextLeafDecode...)
					if nextNodeType == 2 {
						decode = append(decode, 16)
					}
					node.flag_value.encoded_prefix = compact_encode(decode)
					fmt.Println("this node ***** : ", node.String())
					return mpt.updateNode(node, oldHash), nextNodeType
				}

			}
		} else {
			// leaf node
			if lenPath == n && com == n { // completely same, can be deleted
				delete(mpt.db, oldHash)
				return "", 0
			} else {
				return oldHash, -1
			}
		}
	}

	// branch node
	if len(path) == 0 {
		if node.branch_value[16] == "" {
			return oldHash, -1
		}
		node.branch_value[16] = ""
	} else {
		first := path[0]
		nextPath := path[1:]

		if node.branch_value[first] == "" { // search failed, can not go down, nothing need to be changed
			return oldHash, -1
		} else { // can go down
			nextNodeHash := node.branch_value[first]
			nextNode := mpt.db[nextNodeHash]
			newNextHash, status := mpt.deleteHelper(nextNodeHash, nextPath, nextNode)
			if status == -1 {
				return oldHash, -1
			}
			node.branch_value[first] = newNextHash
		}
	}
	// finish update the branch, then check if need to change the branch node
	loc := node.checkNumLeaf()
	if loc == -1 { // there is 2 or more leaf, no need to be changed
		return mpt.updateNode(node, oldHash), 1
	}
	if loc == 16 { // only branch itself saved a value, change to leaf
		node.node_type = 2
		value := node.branch_value[16]
		node.branch_value[16] = ""
		node.flag_value.value = value
		node.flag_value.encoded_prefix = compact_encode([]uint8{16})
		return mpt.updateNode(node, oldHash), 2
	}
	nextNodeHash := node.branch_value[loc]
	nextNode := mpt.db[nextNodeHash]
	nextNodeType := nextNode.node_type

	if nextNodeType == 2 {
		nextNodeEncode := nextNode.flag_value.encoded_prefix
		nextNodeDecode := compact_decode(nextNodeEncode)
		nextNodeValue := nextNode.flag_value.value

		node.node_type = 2
		node.branch_value[loc] = ""
		node.flag_value.value = nextNodeValue
		prefix := append([]uint8{uint8(loc)}, nextNodeDecode...)
		var resType int = 3
		if nextNodeEncode[0]/16 >= 2 { // leaf
			prefix = append(prefix, 16)
			resType = 2
		}
		node.flag_value.encoded_prefix = compact_encode(prefix)
		delete(mpt.db, nextNodeHash)
		return mpt.updateNode(node, oldHash), resType
	}

	if nextNodeType == 1 { // next is a branch node, change this node to extension
		node.node_type = 2
		node.branch_value[loc] = ""
		node.flag_value.value = nextNodeHash
		node.flag_value.encoded_prefix = compact_encode([]uint8{uint8(loc)})
		return mpt.updateNode(node, oldHash), 3
	}

	return oldHash, -1
}

func (node *Node) checkNumLeaf() int {
	if node.node_type != 1 {
		return -1
	}
	var res int = 0
	var cur int = -1
	branch := node.branch_value
	for i := 0; i < 17; i++ {
		if branch[i] != "" {
			res++
			cur = i
		}
	}
	if res == 1 {
		return cur
	}
	return -1
}

//func (mpt *MerklePatriciaTrie) Test(key string) {
//	//fmt.Println(mpt.db[""].node_type)
//	var node Node
//	fmt.Println(mpt.db[node.branch_value[3]].String())
//	//hex1 := getHexArray(key)
//	//fmt.Println(hex1)
//	//
//	//hex2 := compact_encode(hex1)
//	//fmt.Println(hex2)
//
//	//hex3 := compact_encode([]uint8{1,1,2,3,4,5,})
//	//fmt.Println(hex3)
//	//test_compact_encode()
//}

func (mpt *MerklePatriciaTrie) MyTest() {
	mpt.Initial()
	mpt.Insert("p", "apple")
	mpt.Insert("aa", "banana")
}

func getHexArray(key string) []uint8 {
	var res []uint8
	for _, c := range key {
		i := uint8(c)
		first := i / 16
		last := i % 16
		res = append(res, first, last)
	}
	fmt.Println("key: ", key, " hex: ", res)
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
	//fmt.Println(hex_array)
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
	//fmt.Println("test encoded arr: ", encoded_arr)
	flag := encoded_arr[0]
	oddlen := flag / 16 % 2
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
	rs = strings.Replace(rs, "\n", "\r\n", -1)
	return rs
}
