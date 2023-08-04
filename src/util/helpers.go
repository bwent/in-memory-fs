package util

import (
	"fmt"
	"strings"
)

// helper method to split a string into slice of strings separated by "/"
func SplitPath(path string) []string {
	var paths = []string{}
	for _, p := range strings.Split(path, "/") {
		str := strings.TrimSpace(p)
		if str != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

// helper method to check if a file exists in the diven directory
func ExistsInCurrentDir(dir *File, name string, isDir bool) bool {
	return dir.GetChildByName(name) != nil && dir.GetChildByName(name).IsDirectory() == isDir
}

// helper method to recursively traverse the directory tree until we reach the root directory,
// adding directory names to a list as we go
func PwdRecursion(dirs *[]string, curr *File) {
	parent := curr.GetParent()
	if parent == nil {
		// root directory - base case
		*dirs = []string{""}
		return
	}

	PwdRecursion(dirs, parent)
	*dirs = append(*dirs, curr.GetName())
}

func FileSliceToString(data []*File, root *File) []string {
	allMatches := []string{}
	for _, file := range data {
		allMatches = append(allMatches, file.GetFullPathName(root))
	}
	return allMatches
}

func BFS(node *File, target string) []*File {
	if node == nil {
		return nil
	}

	visited := make(map[string]bool)
	queue := queue{node}
	result := []*File{}

	for queue.Size() > 0 {
		next, _ := queue.PopFront()
		if visited[next.GetName()] {
			continue
		}
		visited[next.GetName()] = true

		if next.GetName() == target {
			// Found a match, add it to the result
			result = append(result, next)
		}

		for _, child := range next.GetChildren() {
			queue.PushBack(child)
		}
	}

	return result
}

// recursive helper method to remove files depth-first down to the leaf nodes
func RmRecursion(curr *File) {
	if curr == nil || curr.GetParent() == nil {
		// base case
		return
	}

	delete(curr.GetParent().GetChildren(), curr.GetName())
	for _, c := range curr.GetChildren() {
		// loop through all children nodes and remove subdirectories recursively
		RmRecursion(c)
	}
}

// Helper function to traverse from the current directory to the specified path,
// using an absolute or relative path
func WalkToEndOfPath(pathSplit []string, currentDirectory *File, root *File) (*File, error) {
	wd := currentDirectory

	// If the path name starts with "~", this is an absolute path - start from the root
	// Else start from the current working directory
	if pathSplit[0] == "~" {
		wd = root
		pathSplit = pathSplit[1:]
	}

	for _, name := range pathSplit {
		if name == ".." {
			// If we see ".." we're trying to navigate one directory up in the tree
			// Set the current directory to its parent
			if wd.GetParent() != nil {
				wd = wd.GetParent()
			} else {
				// This means we're already at the root, so we shouldn't need to do anything
			}
		} else if !ExistsInCurrentDir(wd, name, true) {
			return nil, fmt.Errorf("Directory not found: %s", name)
		} else {
			// Advance to the child node by name
			wd = wd.GetChildByName(name)
		}
	}
	return wd, nil
}

func StringSliceToByteSlice(strSlice []string) []byte {
	var byteSlice []byte
	for _, str := range strSlice {
		byteSlice = append(byteSlice, []byte(str)...)
	}
	return byteSlice
}

func ModifyNameToHandleCollisions(name string) string {
	nameSplit := strings.Split(name, ".")
	if len(nameSplit) == 2 {
		name = nameSplit[0] + "1." + nameSplit[1]
	} else {
		name = name + "1"
	}
	return name
}
