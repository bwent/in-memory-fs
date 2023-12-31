package src

import (
	"errors"
	"fmt"
	"in-memory-fs/src/util"
	"strings"
)

type Filesystem struct {
	root             *util.File
	currentDirectory *util.File
}

// Creates a new filesystem and sets the current directory to the root ()
func NewFileSystem() *Filesystem {
	rootDir := util.NewFile("/", true, nil)
	return &Filesystem{
		root:             rootDir,
		currentDirectory: rootDir,
	}
}

// Returns the current working directory, e.g. "/Users/bwent/home"
//
// Parameters: N/A
// Returns:
//
//	string - the current working directory
func (fs *Filesystem) Pwd() string {
	if fs.currentDirectory == fs.root {
		// If we're at the root, simply return "/"
		return "/"
	}
	dirs := []string{}
	// Recursively iterate from the current directory to the root, adding each parent to a list of strings
	util.PwdRecursion(&dirs, fs.currentDirectory)
	return strings.Join(dirs, "/")
}

// Creates a new directory specified by "path" within the current working directory.
// Does NOT currently support using ".." or "~" to create a new directory using either the
// parent or the root directory
//
// Parameters:
//
//	path (string) - can be either a name (e.g. /bwent) or a full path (e.g. /bwent/home/test), as long as each
//	                path element before the final one is an existing directory
//
// Returns:
//
//	string - the newly-created directory name
//	error  - an error if we were unable to successfully create the directory
func (fs *Filesystem) MkDir(path string) (string, error) {
	// Get the current working directory
	wd := fs.currentDirectory

	// Split the path into individual directory names
	pathSplit := util.SplitPath(path)
	length := len(pathSplit)

	// The name of the new directory
	var name string

	if length == 0 {
		return "", errors.New("Must provide at least one directory name")
	} else if length == 1 {
		// If there's only one element, the new dir name is the first element
		name = pathSplit[0]
	} else {
		pathToTraverse := pathSplit[:len(pathSplit)-1]
		leafNode, err := util.WalkToEndOfPath(pathToTraverse, fs.currentDirectory, fs.root)
		if err != nil {
			return "", err
		}
		wd = leafNode
		// Set the dir name to the last element
		name = pathSplit[len(pathSplit)-1]
	}

	// Take the last element and add the new directory
	newDir := util.NewFile(name, true, wd)
	wd.UpsertChild(name, newDir)

	return name, nil
}

// Changes the current working directory to the specified path
//
// Parameters:
//
//	   path (string) - the path we want to navigate to. If prefixed with "~" we will
//						  start from the root. If prefixed with ".." we'll navigate one directory up
//						  in the tree.
//
// Returns:
//
//	string - the current working directory name
//	error  - an error if the path provided is invalid
func (fs *Filesystem) Cd(path string) (string, error) {
	// Traverse to the end of the path specified
	leafNode, err := util.WalkToEndOfPath(util.SplitPath(path), fs.currentDirectory, fs.root)
	if err != nil {
		return "", err
	}
	// Set the current working directory to the last node in the tree
	fs.currentDirectory = leafNode
	return leafNode.GetName(), nil
}

// Lists the contents (files and subdirectories) of the specified path or current directory.
//
// Parameters:
//
//	paths (string) - 0 or 1 paths. If 0 provided, we'll list the contents of the current directory,
//	                 else we'll list the contents of the specified (valid) path
//
// Returns:
//
//	string - the children/contents of the directory, separated by a space
//	error - an error if the specified path is invalid
func (fs *Filesystem) Ls(path ...string) (string, error) {
	var wd *util.File

	if len(path) == 1 {
		splitPath := util.SplitPath(path[0])

		// Traverse to the end of the path
		leafNode, err := util.WalkToEndOfPath(splitPath, fs.currentDirectory, fs.root)
		if err != nil {
			return "", err
		}
		wd = leafNode
	} else {
		wd = fs.currentDirectory
	}

	// Return all the child directory names
	return strings.Join(wd.GetChildrenNames(), " "), nil
}

// Removes a file or directory from the current directory. If a directory is provided, the removal must be recursive unless
// the directory has no children.
// Parameters:
//
//	path (string) -  the path of the file/directory to remove
//	recusrive (bool) - if the removal should be done recursively to remove all sub-directories
//
// Returns:
//
//	string - the removed path name
//	error - an error if the removal was unsuccessful
func (fs *Filesystem) Rm(path string, recursive bool) (string, error) {
	// Sanitize the string
	path = strings.Trim(path, "/")

	wd := fs.currentDirectory

	// Get the file or directory to remove
	toRemove := wd.GetChildByName(path)
	if toRemove == nil {
		return "", fmt.Errorf("Directory not found: %s", path)
	}

	if !recursive {
		// Can only remove non-recursively if this is a non-empty directory
		if toRemove.IsDirectory() && len(toRemove.GetChildren()) > 0 {
			return "", errors.New("Method does not support removing non-empty directories. Use the recursive option")
		}
		// If not recursive, simply remove the path from the children of the current directory
		wd.RemoveChild(path)
	} else {
		// Don't try recursion if the path provided is a file, not a directory
		if !toRemove.IsDirectory() {
			return "", errors.New("Method does not support removing files recursively")
		}
		// Remove the directory and all subdirectories recursively
		util.RmRecursion(toRemove)
	}

	return toRemove.GetName(), nil
}

// Creates a new empty file in the current directory. If the filename already exists, we'll simply append a "1"
// to the end.
// Parameters:
//
//	name (string) - the name of the file to create
//
// Returns:
//
//	string - the newly created file name
//	error - an error if the file was not able to be created
func (fs *Filesystem) MkFile(name string) (string, error) {
	// Set the current working directory
	wd := fs.currentDirectory

	// Check if the name contains the '/' character, which is not supported in filenames
	if strings.ContainsRune(name, '/') {
		return "", errors.New("/ character not supported in filenames")
	}

	// If a file with the same name already exists in the current directory, modify the name to handle collisions
	if util.ExistsInCurrentDir(wd, name, false) {
		name = util.ModifyNameToHandleCollisions(name)
	}

	// Create the new file and set the parent to the working directory
	newFile := util.NewFile(name, false, wd)

	// Add the new file to the children of the current directory
	wd.UpsertChild(name, newFile)

	return name, nil
}

// Writes a string of data to the specified file in the current directory. The max amount of data any
// file can have is 2000000MB or 2GB.
// Parameters:
//
//	name (string) - the name of the file to write
//	data (...string) - the text to write to the file
//
// Returns:
//
//	string - the name of the file we just wrote to
//	error - an error if the file doesn't exist or we've exceeded the max data size (defined in `file.go`)
func (fs *Filesystem) WriteFile(name string, data ...string) (string, error) {
	wd := fs.currentDirectory
	file := wd.GetChildByName(name)

	if file == nil {
		return "", fmt.Errorf("File %s does not exist", name)
	}

	return name, file.WriteFileData(util.StringSliceToByteSlice(data))
}

// Reads the contents of the filename specified. Must be in the curernt directory
//
// Parameters:
//
//	name (string) - the  name of the file to read in
//
// Returns:
//
//	string - the contents of the file, up to 2000 chars (see limit in `util/file.go`)
//	error - an error if the file does not exist
func (fs *Filesystem) ReadFile(name string) (string, error) {
	wd := fs.currentDirectory
	file := wd.GetChildByName(name)

	if file == nil {
		return "", fmt.Errorf("File %s does not exist!", name)
	}

	return file.ReadFileContents(), nil
}

// Moves the specified file (within the current directory) to the specified target directory.
// TODO: Support relative/absolute paths for the source file in the future
//
// Paramters:
//
//	name (string) -   the name of the file to move
//	target (string) - the name of the target directory
//
// Returns:
//
//	string - the name of the target directory if the move was successful
//	error  - an error if the move was unsuccessful
func (fs *Filesystem) MvFile(name string, target string) (string, error) {
	// Sanitize the strings
	name = strings.Trim(name, "/")
	target = strings.Trim(target, "/")

	wd := fs.currentDirectory
	file := wd.GetChildByName(name)

	splitPath := util.SplitPath(target)
	if len(splitPath) == 0 {
		return "", fmt.Errorf("Invalid target path: %s", target)
	}

	// Walk to the end of the path
	targetDir, err := util.WalkToEndOfPath(util.SplitPath(target), fs.currentDirectory, fs.root)
	if err != nil {
		return "", err
	}

	// Validation
	if file == nil {
		return "", fmt.Errorf("File %s does not exist", name)
	}

	if file.IsDirectory() {
		return "", fmt.Errorf("File %s is a directory; cannot move", name)
	}

	if targetDir == nil {
		return "", fmt.Errorf("Target directory %s does not exist", target)
	}

	if !targetDir.IsDirectory() {
		return "", fmt.Errorf("Target path %s is not a directory", target)
	}

	wd.RemoveChild(name)

	if util.ExistsInCurrentDir(targetDir, name, false) {
		// If a file with the same name already exists in the target directory, add a "1" extension
		name = util.ModifyNameToHandleCollisions(name)
		// Update the file name too
		file.SetName(name)
	}

	targetDir.UpsertChild(name, file)
	file.SetParent(targetDir)

	return target, nil
}

// Attempts to find a file or directory within the current working directory (and/or its children)
//
// Parameters:
//
//	target (string) - the name of the file/directory to find
//	searchSubtrees (bool) - whether or not we should search the subdirectories of the current directory
//
// Returns:
//
//	[]string - all matching results represented as a full path
func (fs *Filesystem) FindFileOrDir(target string, searchSubtrees bool) []string {
	if searchSubtrees {
		return util.FileSliceToString(util.BFS(fs.root, target), fs.root)
	}

	result := []string{}
	for key := range fs.currentDirectory.GetChildren() {
		if target == key {
			result = append(result, key)
		}
	}

	return result
}
