package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ngdinhtoan/hari/parser"
)

var (
	inputDir  string
	outputDir string
	showHelp  bool
)

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&inputDir, "input-dir", pwd, "Path to directory which contains json file. Default is working directory.")
	flag.StringVar(&outputDir, "output-dir", "", "Path to directory to which generated go files will be exported. Same as input directory if omit.")
	flag.BoolVar(&showHelp, "help", false, "Print this message")

	flag.Parse()

	if !strings.HasSuffix(inputDir, "/") {
		inputDir += "/"
	}

	if outputDir == "" {
		outputDir = inputDir
	} else if !strings.HasSuffix(outputDir, "/") {
		outputDir += "/"
	}
}

func main() {
	if showHelp {
		flag.Usage()
		os.Exit(0)
	}

	files, _ := filepath.Glob(inputDir + "/*.json")

	if len(files) == 0 {
		fmt.Printf("No json files in %q folder.\n", inputDir)
		os.Exit(1)
	}

	for _, file := range files {
		parseFile(file)
	}

	log.Println("Finish!")
	os.Exit(0)
}

// parseFile will read and parse JSON string from input file,
// then write GO struct to output file.
func parseFile(file string) error {
	log.Printf("Parsing %q\n", file)

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	filename := strings.TrimPrefix(strings.TrimSuffix(file, ".json"), inputDir)

	rs := make(chan *parser.Struct)
	errs := make(chan error)
	done := make(chan bool)

	f, err := os.Create(outputDir + filename + ".go")
	if err != nil {
		return err
	}
	defer f.Close()

	newline := []byte("\n\n")
	w := bufio.NewWriter(f)

	w.WriteString("package main")
	w.Flush()

	go parser.Parse(parser.ToCamelCase(filename), content, rs, errs, done)

	for {
		select {
		case <-done:
			close(rs)
			close(errs)
			close(done)
			return nil
		case s := <-rs:
			w.Write(newline)
			s.WriteTo(w)
			w.Flush()
		}
	}
}
