package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sort"
)

type Node struct {
	Name     string
	IsDir    bool
	Size     int64
	SubNodes []*Node
}

func (node *Node) AddSubNodes(path string, withFiles bool) error {
	files, err := ioutil.ReadDir(path)
	sort.Slice(files, func(i int, j int) bool { return files[i].Name() < files[j].Name() })
	if err != nil {
		return err
	}
	for _, file := range files {
		subNode := &Node{
			Name:  file.Name(),
			IsDir: file.IsDir(),
			Size:  file.Size(),
		}
		if file.IsDir() {
			subNode.AddSubNodes(path+string(os.PathSeparator)+file.Name(), withFiles)
		}
		if withFiles || file.IsDir() {
			node.SubNodes = append(node.SubNodes, subNode)
		}
	}
	return nil
}

func (node *Node) PrintDirTree(out io.Writer, prefix string) error {
	indexOfLastNode := len(node.SubNodes) - 1
	for i, subNode := range node.SubNodes {
		var newPrefix string
		var firstCharacter string

		if i == indexOfLastNode {
			firstCharacter = "└"
			if subNode.IsDir {
				newPrefix = prefix + "\t"
			}
		} else {
			firstCharacter = "├"
			newPrefix = prefix + "│\t"
		}
		fmt.Fprintf(out, "%s%s───%s%s\n", prefix, firstCharacter, subNode.Name, subNode.FormatSize())
		if subNode.IsDir {
			subNode.PrintDirTree(out, newPrefix)
		}
	}
	return nil
}

func (node *Node) FormatSize() string {
	if node.IsDir {
		return ""
	}
	if node.Size == 0 {
		return " (empty)"
	}
	return fmt.Sprintf(" (%db)", node.Size)
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	node := &Node{
		Name:  "root",
		IsDir: printFiles,
	}
	err := node.AddSubNodes(path, printFiles)
	if err != nil {
		return err
	}
	node.PrintDirTree(out, "")
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
