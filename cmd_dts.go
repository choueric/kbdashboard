package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/choueric/kbdashboard/tree"
)

type DtsContext struct {
	wait       sync.WaitGroup
	baseDir    string
	includeDir string
}

type FileNode struct {
	filePath string // absolute path
}

var dtsHelp = &helpNode{
	cmd:      "dts",
	synopsis: "List relevant DTS files.",
	usage: func(w io.Writer) {
		cmdTitle(w, true, "dts list [-t|-v]")
		cmdUsage(w, "List all relevant DTS files. By default, print as a list.\n"+
			"-t: Print as a tree-like graph.\n"+
			"-v: Print absolute file path.\n")

		cmdTitle(w, false, "dts link [-o directory]")
		cmdUsage(w, "Create all relevant DTSs' soft-link into a dirctory.\n"+
			"-o: Specify the output directory.\n"+
			"    Without -o, the default is '[profile_name]_dts' in current path.\n")
	},
}

// return the absolute path.
func parseIncludeItem(line string, context *DtsContext) string {
	if strings.Contains(line, "\"") {
		s := strings.Split(line, "\"")[1]
		return path.Join(context.baseDir, s)
	} else if strings.Contains(line, "<") {
		s := strings.Split(line, "<")[1]
		s = s[:len(s)-1]
		return path.Join(context.includeDir, s)
	} else {
		return ""
	}
}

func searchInclude(node *tree.Node, context *DtsContext) error {
	defer context.wait.Done()
	val := node.Value.(*FileNode)

	f, err := os.Open(val.filePath)
	if err != nil {
		logger.Println(err)
		return err
	}
	defer f.Close()

	doContainMark := func(line string) bool {
		if strings.Contains(line, "#include") ||
			strings.Contains(line, "/include/") {
			return true
		}
		return false
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if doContainMark(line) {
			n := parseIncludeItem(line, context)
			sub := node.AddSubNode(&FileNode{n})

			context.wait.Add(1)
			go searchInclude(sub, context)
		}
	}

	return nil
}

// parseDTSFiles scans all included fils from a DTS file and builds a tree to
// contain all dts files. At last, return root of this tree.
func parseDTSFiles(dtsFile string) (*tree.Node, error) {
	var context DtsContext
	context.baseDir = path.Dir(dtsFile)
	context.includeDir = path.Join(context.baseDir, "../../../../include")

	filePath, err := filepath.Abs(dtsFile)
	if err != nil {
		return nil, err
	}

	root := tree.New(&FileNode{filePath})
	context.wait.Add(1)

	go searchInclude(root, &context)
	context.wait.Wait()

	return root, nil
}

// make a list of all nodes, remove the duplication.
// f is used to make the node string.
func makeFileList(n *tree.Node, f func(*tree.Node) string) map[string]int {
	maps := map[string]int{}
	n.Walk(func(node *tree.Node) {
		maps[f(node)] = 0
	})

	return maps
}

////////////////////////////////////////////////////////////////////////////////

func getDtsAbsPath(p *Profile) string {
	return path.Join(p.SrcDir, "arch", p.Arch, "boot/dts", p.DTB[0:len(p.DTB)-1]+"s")
}

func dtsListHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	if p.DTB == "" {
		return 0, errors.New("no specified DTS")
	}

	var printTree, verbose bool

	flagSet := flag.NewFlagSet("dtsList", flag.ExitOnError)
	flagSet.BoolVar(&printTree, "t", false, "print Tree.")
	flagSet.BoolVar(&verbose, "v", false, "print complete path.")
	flagSet.Parse(args)

	root, err := parseDTSFiles(getDtsAbsPath(p))
	if err != nil {
		return 0, err
	}

	f := func(n *tree.Node) string {
		val := n.Value.(*FileNode)
		if verbose {
			return val.filePath
		} else {
			return filepath.Base(val.filePath)
		}
	}

	if printTree {
		root.PrintTree(os.Stdout, f)
	} else {
		maps := makeFileList(root, f)
		for k, _ := range maps {
			fmt.Fprintln(os.Stdout, k)
		}
	}

	return 0, nil
}

////////////////////////////////////////////////////////////////////////////////

func dtsLinkHandler(args []string, data interface{}) (int, error) {
	p, _, err := getCurrentProfile(gConfig)
	if err != nil {
		return 0, err
	}

	var outputDir string
	flagSet := flag.NewFlagSet("dtsLink", flag.ExitOnError)
	flagSet.StringVar(&outputDir, "o", p.Name+"_dts", "Output directory name.")
	flagSet.Parse(args)

	fmt.Println(cWrap(cGREEN, "dts link"), "to:", cWrap(cGREEN, outputDir))
	outputDir, _ = filepath.Abs(outputDir)

	err = checkDirExist(outputDir)
	if err != nil {
		return 0, err
	}

	root, err := parseDTSFiles(getDtsAbsPath(p))
	if err != nil {
		return 0, err
	}

	f := func(n *tree.Node) string {
		val := n.Value.(*FileNode)
		return val.filePath
	}
	maps := makeFileList(root, f)

	for k, _ := range maps {
		err := os.Symlink(k, path.Join(outputDir, filepath.Base(k)))
		if err != nil {
			return 0, err
		}
	}

	return 0, nil
}
