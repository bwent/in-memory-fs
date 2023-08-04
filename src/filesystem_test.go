// filesystem_test.go
package src

import (
	"fmt"
	"in-memory-fs/src/util"
	"strings"
	"testing"
)

func TestNewFileSystem(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// Verify that the root directory is created correctly
	if fs.root.GetName() != "/" {
		t.Errorf("Expected root directory name to be / but got: %s", fs.root.GetName())
	}

	// Verify that the current directory is set to the root directory
	if fs.currentDirectory != fs.root {
		t.Errorf("Expected current directory to be the root directory")
	}
}

func TestPwd(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// Working directory should be "/"
	res := fs.Pwd()
	if res != "/" {
		t.Errorf("Expected the current working directory to be / but is %s", res)
	}

	// Create a new directory and navigate to it
	fs.MkDir("/home")
	fs.Cd("/home")

	res = fs.Pwd()
	if res != "/home" {
		t.Errorf("Expected the current working directory to be /home but is %s", res)
	}

	// Now add another directory and navigate to it
	fs.MkDir("/test")
	fs.Cd("/test")

	res = fs.Pwd()
	if res != "/home/test" {
		t.Errorf("Expected the current working directory to be /home/test but is %s", res)
	}
}

func TestMkDir(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// Create a new directory
	res, err := fs.MkDir("/home")

	// Assert no errors
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
	}

	// Assert directory was created with the right name
	if res != "home" {
		t.Errorf("Expected the new directory to be /home but got %s", res)
	}

	res, err = fs.MkDir("/home/bwent")

	// Assert no errors
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
	}

	// Assert directory was created with the right name
	if res != "bwent" {
		t.Errorf("Expected the new directory to be /bwent but got %s", res)
	}

	// Create a new invalid directory should throw error
	res, err = fs.MkDir("/invalid/path")
	assertErrorAndEmptyResult(res, err, "Directory not found: invalid", t)

	// Create a new empty directory should throw error
	res, err = fs.MkDir(" ")
	assertErrorAndEmptyResult(res, err, "Must provide at least one directory name", t)
}

func TestCd(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// First attempt to CD into a non-existing directory
	res, err := fs.Cd("/dir1")
	assertErrorAndEmptyResult(res, err, "Directory not found: dir1", t)

	fs.MkDir("/dir1")
	fs.MkDir("/dir1/dir2")
	fs.MkDir("/dir1/dir2/dir3")

	// Test cases
	res, err = fs.Cd("/dir1")
	assertMatchesAndNoErrors(res, err, "dir1", t)

	// Should start from the root and traverse the path
	res, err = fs.Cd("~/dir1/dir2")
	assertMatchesAndNoErrors(res, err, "dir2", t)

	// Should navigate one level up to "home"
	res, err = fs.Cd("../")
	assertMatchesAndNoErrors(res, err, "dir1", t)

	// Invalid path - should return an error
	res, err = fs.Cd("/test1")
	assertErrorAndEmptyResult(res, err, "Directory not found: test1", t)
}

func TestLs(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// List the contents of the root directory. Should be empty
	res, err := fs.Ls()
	assertMatchesAndNoErrors(res, err, "", t)

	fs.MkDir("/home")
	fs.MkDir("/home/test")
	fs.MkDir("/home/test/foo")
	// List contents of directory by path should return the single directory under that path
	res, err = fs.Ls("/home/test")
	assertMatchesAndNoErrors(res, err, "foo", t)

	// Test with absolute path
	res, err = fs.Ls("~/home")
	assertMatchesAndNoErrors(res, err, "test", t)

	// Test with ".."
	fs.Cd("/home/test")
	res, err = fs.Ls("../")
	assertMatchesAndNoErrors(res, err, "test", t)
}

func TestRm(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	fs.MkDir("dir1")
	fs.MkDir("dir1/dir2")
	// Shouldn't be able to non-recursively remove a directory
	res, err := fs.Rm("/dir1", false)
	assertErrorAndEmptyResult(res, err, "Method does not support removing non-empty directories. Use the recursive option", t)

	// Shouldn't be able to remove a nonexistent directory
	res, err = fs.Rm("/test", false)
	assertErrorAndEmptyResult(res, err, "Directory not found: test", t)

	// Happy path 1
	res, err = fs.Rm("/dir1", true)
	assertMatchesAndNoErrors(res, err, "dir1", t)

	res, err = fs.Ls()
	if res != "" {
		t.Errorf("Expected the current directory to be empty after removing dir1 but instead was %s", res)
	}
}

func TestMkFile(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// Test creating a new file
	res, err := fs.MkFile("test.txt")
	assertMatchesAndNoErrors(res, err, "test.txt", t)

	// Test collisions
	res, err = fs.MkFile("test.txt")
	assertMatchesAndNoErrors(res, err, "test1.txt", t)
}

func TestWriteAndReadFile(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	// Write some data to a non-existent file; should thrown an error
	res, err := fs.WriteFile("test.txt", "hello world!")
	assertErrorAndEmptyResult(res, err, "File test.txt does not exist", t)

	// Create a new file
	res, err = fs.MkFile("test.txt")

	// Write some data to the file
	text := "hello world!"
	res, err = fs.WriteFile("test.txt", text)
	assertMatchesAndNoErrors(res, err, "test.txt", t)

	// Read the data back
	res, err = fs.ReadFile("test.txt")
	assertMatchesAndNoErrors(res, err, text, t)

	moreText := " I am a computer."
	fs.WriteFile("test.txt", moreText)
	res, err = fs.ReadFile("test.txt")
	assertMatchesAndNoErrors(res, err, text+moreText, t)

	fs.MkFile("large-text-test.txt")
	var builder strings.Builder
	// Write 1 more than the max allowed char limit
	for i := 0; i < util.MaxFileReadSize+1; i++ {
		builder.WriteRune('x')
	}
	fs.WriteFile("large-text-test.txt", builder.String())
	res, err = fs.ReadFile("large-text-test.txt")
	expected := fmt.Sprintf("%s ...[trunated contents after %d chars]", builder.String()[:util.MaxFileReadSize+1], util.MaxFileReadSize)
	assertMatchesAndNoErrors(res, err, expected, t)
}

func TestMoveFile(t *testing.T) {
	// Set up test subject
	fs := NewFileSystem()

	fs.MkDir("dir1")
	// Test moving nonexistent file
	res, err := fs.MvFile("file1", "dir1")
	assertErrorAndEmptyResult(res, err, "File file1 does not exist", t)

	fs.MkFile("file1")
	// Test moving file to invalid directory
	res, err = fs.MvFile("file1", "dir2")
	assertErrorAndEmptyResult(res, err, "Directory not found: dir2", t)

	fs.MkDir("dir2")
	// Test moving directory
	res, err = fs.MvFile("dir1", "dir2")
	assertErrorAndEmptyResult(res, err, "File dir1 is a directory; cannot move", t)

	fs.MkFile("file2")
	// Test moving directory
	res, err = fs.MvFile("file2", "file1")
	assertErrorAndEmptyResult(res, err, "Directory not found: file1", t)

	// Happy path
	res, err = fs.MvFile("file1", "dir1")
	assertMatchesAndNoErrors(res, err, "dir1", t)
	res, err = fs.MvFile("file2", "dir1/")
	assertMatchesAndNoErrors(res, err, "dir1", t)

	fs.MkDir("dir1/test1")
	fs.MkFile("file3")
	res, err = fs.MvFile("file3", "~/dir1/test1")
	assertMatchesAndNoErrors(res, err, "~/dir1/test1", t)
}

func TestFind(t *testing.T) {
	// Set up the test subject
	fs := NewFileSystem()

	fs.MkDir("dir1")
	fs.MkDir("dir2")
	fs.MkFile("file1.txt")
	fs.MkFile("file2.txt")

	// Test searching in the current directory
	res := fs.FindFileOrDir("file1.txt", false)
	expected := []string{"file1.txt"}
	if !stringSliceEqual(res, expected) {
		t.Errorf("Invalid results: got: %v, expected: %v", res, expected)
	}

	fs.MvFile("file1.txt", "dir1")
	// Test searching in subdirectories
	res = fs.FindFileOrDir("file1.txt", true)
	expected = []string{"/dir1/file1.txt"}
	if !stringSliceEqual(res, expected) {
		t.Errorf("Invalid results: got: %v, expected: %v", res, expected)
	}

	// Test finding a nonexistent file yields empty results
	res = fs.FindFileOrDir("file3.txt", true)
	expected = []string{}
	if !stringSliceEqual(res, expected) {
		t.Errorf("Invalid results: got: %v, expected: %v", res, expected)
	}
}

// HELPER METHODS

func assertMatchesAndNoErrors(res string, err error, expected string, t *testing.T) {
	if err != nil {
		t.Errorf("Expected no errors but got %s", err.Error())
	}

	if res != expected {
		t.Errorf("Expected path to be %s but was %s", expected, res)
	}
}

func assertErrorAndEmptyResult(res string, err error, errTxt string, t *testing.T) {
	// Assert error is thrown
	if err == nil || err.Error() != errTxt {
		t.Errorf("Expected error: %s but got %s", errTxt, err)
	}

	// Assert result is empty
	if res != "" {
		t.Errorf("Expected an empty result but got %s", res)
	}
}

func stringSliceEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}
