# in-memory-fs
An in-memory filesystem implemented in Go. Supports basic filesystem operations, including creating new directories/files (`mkdir`), navigating to a directory (`cd`), listing the contents of a directory (`ls`), printing the working directory path (`pwd`) and several other operations. 

## Setup
### Requirements
You'll need [Go](https://go.dev) installed locally. Download the latest version [here](https://go.dev/dl/).

### Clone the repository
```
# via HTTPS
$ git clone https://github.com/bwent/in-memory-fs.git

# via Git cli
gh repo clone bwent/in-memory-fs
```

### Navigate to directory
```
$ cd in-memory-fs
```

### Run the application
```
$ go run main.go
```
You'll then be prompted for input. See the [Usage](#usage) section below for more details on how to use the filesystem.

### Run tetsts
```
# From in-memory-fs directory
$ cd src/
# Runs all the unit tests in the current directory
$ go test
```
Voila!

### Code Structure
* `src` contains all the source code for running theh application, apart from `main.go`
    * `filesystem.go` is where the main filesystem methods are implemented
* the `util` package contains auxiliary files and helpers
* `filesystem_test.go` contains unit tests for `filesystem` methods. 

## Usage

* `help` to view options
* `exit` to exit the program
* `mkdir <name>` - Creates a new directory with the specified name within the current directory. 
* `pwd`  - Prints the current working directory.
* `cd <path>` - Changes the current working directory to the specified path.
* `ls [path]` Lists the contents (files and subdirectories) of the specified path. If none provided, uses the current directory
* `rm <path> <useRecursion>` - Removes a file (not a directory). Set `useRecursion` to true to remove directories and all subdirectories.
* `mkfile <name>` - Creates a new empty file in the current directory.
* `writeFile <name>`  - Writes contents to the specified file in the current directory.
* `readFile <name>`    - Reads the contents of the specified file in the current directory (truncated after 2000 chars)
* `mvfile <name> <target>`  - Moves the specified file to the given target directory.
* `find <name> <useRecursion> `  - Finds files or directories with the specified name. Set `useRecursion` to true to search subdirectories.
* `link <target> <name> `  - Creates a hard link to the specified target with the given name. Only supports hard links for files
* `symlink <target> <name>` - Creates a symbolic link (symlink) to the specified target (file or directory) with the given name.  

### Testing
```
# Run from in-memory-fs directory
$ cd src/
# This will run all unit tests in the module
$ go test
```

## Notes
### TODOs
* Add unit tests for all `util` class files; add additional unit tests to check for more edge cases
* Add concurrency controls by building a Mutex into the File class and locking around write (and possibly also read) operations
* Add symlink and hard link support

