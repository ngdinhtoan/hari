package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/ngdinhtoan/flagstruct"
	"github.com/ngdinhtoan/hari/parser"
)

type appConfig struct {
	InputDir   string `flag:"input-dir" usage:"Path to directory which contains json file. Default is working directory."`
	OutputDir  string `flag:"output-dir" usage:"Path to directory to which generated go files will be exported. Same as input directory if omit."`
	ShowHelp   bool   `flag:"help" usage:"Print this message."`
	OnTheFly   bool   `flag:"on-the-fly" usage:"Read JSON string from stdin pipe, this action require --struct-name."`
	StructName string `flag:"struct-name" usage:"Name of root struct, this is required for --on-the-fly."`
}

var conf appConfig
var newline = []byte("\n\n")

func init() {
	conf = appConfig{}
	flagstruct.Parse(&conf)

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if conf.InputDir == "" {
		conf.InputDir = pwd
	}

	if !strings.HasSuffix(conf.InputDir, "/") {
		conf.InputDir += "/"
	}

	if conf.OutputDir == "" {
		conf.OutputDir = conf.InputDir
	} else if !strings.HasSuffix(conf.OutputDir, "/") {
		conf.OutputDir += "/"
	}

	if conf.OnTheFly && conf.StructName == "" {
		fmt.Println("You are using on-the-fly option, struct name is required.")
		conf.ShowHelp = true
	}
}

func main() {
	if conf.ShowHelp {
		flag.Usage()
		os.Exit(0)
	}

	if conf.OnTheFly {
		reader := bufio.NewReader(os.Stdin)
		contentBuf := bytes.Buffer{}
		reader.WriteTo(&contentBuf)

		f, err := os.Create(conf.OutputDir + conf.StructName + ".go")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		if err := parseJSONString(contentBuf.Bytes(), conf.StructName, w); err != nil {
			log.Fatal(err)
		}
		w.Flush()
	} else {
		files, _ := filepath.Glob(conf.InputDir + "/*.json")

		if len(files) == 0 {
			fmt.Printf("No json files in %q folder.\n", conf.InputDir)
			os.Exit(1)
		}

		for _, file := range files {
			parseFile(file)
		}
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

	filename := strings.TrimPrefix(strings.TrimSuffix(file, ".json"), conf.InputDir)

	f, err := os.Create(conf.OutputDir + filename + ".go")
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	if err := parseJSONString(content, filename, w); err != nil {
		return err
	}

	return w.Flush()
}

func parseJSONString(content []byte, structName string, w io.Writer) error {
	w.Write([]byte("package main"))

	rs := make(chan *parser.Struct)
	errs := make(chan error)
	done := make(chan bool)

	defer close(rs)
	defer close(errs)
	defer close(done)

	go parser.Parse(parser.ToCamelCase(structName), content, rs, errs, done)

	for {
		select {
		case err := <-errs:
			return err
		case <-done:
			return nil
		case s := <-rs:
			w.Write(newline)
			s.WriteTo(w)
		}
	}
}
