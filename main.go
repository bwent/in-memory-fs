package main

import (
	"bufio"
	"fmt"
	"in-memory-fs/src"
	"os"
	"strconv"
	"strings"
)

// Maps a valid method to its acceptable number of inputs
var ValidInputMap = map[string][]int{
	"pwd":    {0},
	"mkdir":  {1},
	"cd":     {1},
	"ls":     {0, 1},
	"rm":     {1, 2},
	"mkfile": {1},
	// -1 indicates we have no bounds on the input size
	"writefile": {-1},
	"readfile":  {1},
	"mvfile":    {2},
	"find":      {2},
	"symlink":   {2},
	"link":      {2},
}

const HelpText string = `Commands:
pwd              	Prints the current working directory.
mkdir <path>        	Creates a new directory within the current working directory.
cd <path>           	Changes the current working directory to the specified path.
ls [path]           	Lists the contents (files and subdirectories) of the specified path.
rm <path> <useRecursion>    	Removes a file (not a directory). Set useRecursion to true to remove directories recursively.
mkfile <name>       	Creates a new empty file in the current directory.
writeFile <name>    	Writes contents to the specified file in the current directory.
readFile <name>     	Reads the contents of the specified file in the current directory.
mvfile <name> <target>  	Moves the specified file to the given target directory.
find <name> <useRecursion>     	Finds files or directories with the specified name. Set useRecursion to true to search subdirectories.
link <target> <name>   Creates a hard link to the specified target with the given name. Only supports hard links for files
symlink <target> <name>  Creates a symbolic link (symlink) to the specified target (file or directory) with the given name.  
-------------------------
help                	Displays this help menu.
exit                	Exits the program.`

func main() {
	fs := src.NewFileSystem()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command (or 'exit' to quit): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error parsing input: ", err)
			return
		}

		keyword := strings.TrimSpace(input)

		switch keyword {
		case "exit":
			fmt.Println("Exiting")
			return
		case "help":
			fmt.Println(HelpText)
			return
		default:
			err := parseUserInputs(fs, strings.Split(input, " "))
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}

func validateInputs(method string, inputs []string) error {
	validInputSizes := ValidInputMap[method]
	if validInputSizes == nil {
		return fmt.Errorf("Invalid method %s- run 'help' for guidance", method)
	}
	if !contains(validInputSizes, len(inputs)) && validInputSizes[0] != -1 {
		validSizesStr := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(validInputSizes)), ", "), "[]")
		return fmt.Errorf("Invalid input length (saw=%d, expected=%s). Run 'help' for guidance", len(inputs), validSizesStr)
	}
	return nil
}

func parseUserInputs(fs *src.Filesystem, inputs []string) error {
	method := inputs[0]
	method = strings.ToLower(method)
	method = strings.TrimSpace(method)

	params := inputs[1:]

	err := validateInputs(method, params)
	if err != nil {
		return err
	}
	for i := range params {
		params[i] = strings.TrimSpace(params[i])
	}

	switch method {
	case "pwd":
		fmt.Println(fs.Pwd())
	case "mkdir":
		printResults(fs.MkDir(params[0]))
	case "cd":
		printResults(fs.Cd(params[0]))
	case "ls":
		if len(params) == 0 {
			printResults(fs.Ls())
		} else {
			printResults(fs.Ls(params[0]))
		}
	case "rm":
		useRecursion := false
		var err error
		if len(params) == 2 {
			useRecursion, err = strconv.ParseBool(params[1])
			if err != nil {
				fmt.Println("Invalid second parameter: must be among {true, false, T, F, 0, 1}")
			}
		}
		printResults(fs.Rm(params[0], useRecursion))
	case "mkfile":
		printResults(fs.MkFile(params[0]))
	case "writefile":
		printResults(fs.WriteFile(params[0], params[1:]...))
	case "readfile":
		printResults(fs.ReadFile(params[0]))
	case "mvfile":
		printResults(fs.MvFile(params[0], params[1]))
	case "find":
		bVal, err := strconv.ParseBool(params[1])
		if err != nil {
			fmt.Println("Invalid second parameter: must be among {true, false, T, F, 0, 1}")
		}
		res := fs.FindFileOrDir(params[0], bVal)
		fmt.Println(strings.Join(res, ","))
	case "link":
		printResults(fs.CreateHardlink(params[0], params[1]))
	case "symlink":
		printResults(fs.CreateSymlink(params[0], params[1]))
	default:
		return fmt.Errorf("Invalid method call %s - please run 'help' for more details", method)
	}
	return nil
}

func printResults(res string, err error) {
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}

func contains(slice []int, val int) bool {
	elementsMap := make(map[int]bool)
	for _, num := range slice {
		elementsMap[num] = true
	}

	return elementsMap[val]
}
