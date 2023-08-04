package util

import (
	"fmt"
	"strings"
)

// Splits a string into slice of strings separated by "/"
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

// Check if a file exists in the diven directory. "isDir" is used to specify whether we should
// check if it's a file or directory
func ExistsInCurrentDir(dir *File, name string, isDir bool) bool {
	return dir.GetChildByName(name) != nil && dir.GetChildByName(name).IsDirectory() == isDir
}

// Recursively traverse the directory tree until we reach the root directory,
// adding the current directory names to a list as we go
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

// Convert a slice of Files into a string slice, using the filename
func FileSliceToString(data []*File, root *File) []string {
	allMatches := []string{}
	for _, file := range data {
		allMatches = append(allMatches, file.GetFullPathName(root))
	}
	return allMatches
}

// Breadth-first serach implementation used for searching files within the filesystem
// Uses a map
func BFS(node *File, target string) []*File {
	if node == nil {
		return nil
	}

	// Keep track of all nodes we've already visited (optimization)
	visited := make(map[string]bool)

	// Use a queue for inspecting nodes
	queue := queue{node}
	result := []*File{}

	for queue.Size() > 0 {
		// Take the next node off the queue
		next, _ := queue.PopFront()
		// If we've already seen it, skip
		if visited[next.GetName()] {
			continue
		}
		visited[next.GetName()] = true

		if next.GetName() == target {
			// Found a match, so add it to the result
			result = append(result, next)
		}

		// Add all the child nodes to the queue for inspection
		for _, child := range next.GetChildren() {
			queue.PushBack(child)
		}
	}

	// Empty result indicates none found
	return result
}

// Recursively remove files depth-first down to the leaf nodes
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

// Traverse from the current directory to the specified path, using an absolute or relative path
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

// Convert a slice of strings to a byte slice
func StringSliceToByteSlice(strSlice []string) []byte {
	var byteSlice []byte
	for i, str := range strSlice {
		if i > 0 {
			byteSlice = append(byteSlice, ' ')
		}
		byteSlice = append(byteSlice, []byte(str)...)
	}
	return byteSlice
}

// Add a special extension in case we're attempting to create a duplicate file
func ModifyNameToHandleCollisions(name string) string {
	nameSplit := strings.Split(name, ".")
	if len(nameSplit) == 2 {
		name = nameSplit[0] + "1." + nameSplit[1]
	} else {
		name = name + "1"
	}
	return name
}
