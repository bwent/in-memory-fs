package util

import (
	"fmt"
	"strings"
)

// Limit the number of bytes that can be written to any file to 2M bytes, or ~2MB
const MaxFileSize int = 2000000

// Limit the size of the string that can be returned when reading a file to 2000 chars
const MaxFileReadSize int = 2000

// Stores information about a File or Directory object
type File struct {
	name        string
	contents    []byte
	isDirectory bool
	children    map[string]*File
	parent      *File
}

// NewFile creates a new File instance with the given name, isDir flag, and parent file.
func NewFile(name string, isDir bool, parent *File) *File {
	return &File{
		name:        name,
		isDirectory: isDir,
		contents:    []byte{},
		children:    make(map[string]*File),
		parent:      parent,
	}
}

// Simple Getters
func (f *File) GetName() string {
	return f.name
}

func (f *File) IsDirectory() bool {
	return f.isDirectory
}

func (f *File) GetChildren() map[string]*File {
	return f.children
}

func (f *File) GetChildrenNames() []string {
	var childrenNames []string
	for _, c := range f.children {
		if c != nil {
			childrenNames = append(childrenNames, c.name)
		}
	}
	return childrenNames
}

func (f *File) GetChildByName(name string) *File {
	return f.children[name]
}

func (f *File) GetParent() *File {
	return f.parent
}

// Reads the contents of a file into a string, cutting off after `MaxFileReadSize` chars
func (f *File) ReadFileContents() string {
	str := string(f.contents)
	if len(str) > MaxFileReadSize {
		strSpl := strings.SplitAfterN(str, ",", MaxFileReadSize)
		str = fmt.Sprintf("%s ...[trunated contents after %d chars]", strSpl[0], MaxFileReadSize)
	}
	return str
}

// Returns the full path name of a given file (e.g.'/Users/bwent/test1')
func (f *File) GetFullPathName(root *File) string {
	return getFullPathNameHelper(f, root)
}

// Write methods
func (f *File) UpsertChild(name string, file *File) {
	f.children[name] = file
}

func (f *File) RemoveChild(name string) {
	delete(f.children, name)
}

func (f *File) SetParent(parent *File) {
	f.parent = parent
}

// Writes the specified data (represented as a byte slice) to a file
// Returns an error if the newData + exisitng contents exceeds `MaxFileSize`
func (f *File) WriteFileData(data []byte) error {
	totalSize := len(f.contents) + len(data)
	if totalSize > MaxFileSize {
		return fmt.Errorf("Exceeded max file size: size=%d, max=%d", totalSize, MaxFileSize)
	}
	f.contents = append(f.contents, data...)
	return nil
}

// Helper function to get the full path name of a file by recursively traversing up the tree
func getFullPathNameHelper(curr *File, start *File) string {
	if curr == start || curr.parent == nil {
		// Root directory or nil - base case
		return ""
	}

	parent := curr.parent
	return getFullPathNameHelper(parent, start) + "/" + curr.name
}
