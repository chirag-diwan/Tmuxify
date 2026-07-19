package tmuxify

import (
	"fmt"
	"strings"
)

type RadixNode struct{
	children []*RadixNode
	segment string
}

func AddPath(root * RadixNode ,path string){
	segments := strings.Split(path, "/")
	var curr *RadixNode = root
	for _ , seg := range segments {
		contains := false
		for _ , child := range curr.children {
			if child.segment == seg{
				contains = true
				curr = child
				break
			}
		}
		if !contains {
			curr.children = append(curr.children, &RadixNode{children: []*RadixNode{} , segment: seg})
		}
	}
}


func getMatchRecurse(curr *RadixNode, keyWord string, parentPath string) []string {
	path := parentPath;
	if curr.segment != "/" && len(curr.segment) != 0 {
		path = parentPath + "/" + curr.segment
	}
	var ret []string

	if strings.Contains(curr.segment, keyWord) {
		ret = append(ret, path)
	}

	for _, child := range curr.children {
		ret = append(ret, getMatchRecurse(child, keyWord, path)...)
	}

	return ret
}


func GetMatch(root *RadixNode , key_word string) []string{
	return getMatchRecurse(root , key_word , "")
}

func PrintTree(root *RadixNode) {
	if root == nil {
		return
	}

	fmt.Println(root.segment)

	for i, child := range root.children {
		printTree(child, "", i == len(root.children)-1)
	}
}

func printTree(node *RadixNode, prefix string, isLast bool) {
	if node == nil {
		return
	}

	connector := "├── "
	nextPrefix := prefix + "│   "

	if isLast {
		connector = "└── "
		nextPrefix = prefix + "    "
	}

	fmt.Println(prefix + connector + node.segment)

	for i, child := range node.children {
		printTree(child, nextPrefix, i == len(node.children)-1)
	}
}
