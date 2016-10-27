// main
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/kurehajime/tmpl"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/transform"
)

func main() {
	var csvStr string
	var templateStr string
	var err error
	var encodeTemplate string
	var encodeCsv string
	var pathTemlate string
	var pathCsv string
	var nameCol int
	var tsv bool

	var defaultEncoding string
	var output string

	if runtime.GOOS == "windows" {
		defaultEncoding = "sjis"
	} else {
		defaultEncoding = "utf-8"
	}

	flag.StringVar(&pathTemlate, "t", "", "[must]template path")
	flag.StringVar(&pathCsv, "c", "", "[must]csv path")
	flag.StringVar(&encodeTemplate, "te", defaultEncoding, "template encoding")
	flag.StringVar(&encodeCsv, "ce", defaultEncoding, "csv encoding")
	flag.StringVar(&output, "o", "./", "output path or file")
	flag.IntVar(&nameCol, "n", -1, "Name column no")
	flag.BoolVar(&tsv, "tsv", false, "TSV:Tab-Separated Values")

	flag.Usage = func() {
		fmt.Println("tmpl makes files that replaced template text with csv by matched column name")
		flag.PrintDefaults()
	}
	flag.Parse()
	if pathTemlate == "" || pathCsv == "" {
		flag.Usage()
		os.Exit(0)
	}

	csvStr, err = readFileByArg(pathCsv)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	csvStr, err = transEnc(csvStr, encodeCsv)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	templateStr, err = readFileByArg(pathTemlate)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	templateStr, err = transEnc(templateStr, encodeTemplate)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	ch := make(chan tmpl.Result)
	go tmpl.Generate(templateStr, csvStr, nameCol, ch, tsv)
	writeFile(output, pathTemlate, ch)
}

func readFileByArg(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func transEnc(text string, encode string) (string, error) {
	body := []byte(text)
	var f []byte

	encodings := []string{"sjis", "utf-8"}
	if encode != "" {
		encodings = append([]string{encode}, encodings...)
	}
	for _, enc := range encodings {
		if enc != "" {
			ee, _ := charset.Lookup(enc)
			if ee == nil {
				continue
			}
			var buf bytes.Buffer
			ic := transform.NewWriter(&buf, ee.NewDecoder())
			_, err := ic.Write(body)
			if err != nil {
				continue
			}
			err = ic.Close()
			if err != nil {
				continue
			}
			f = buf.Bytes()
			break
		}
	}
	return string(f), nil
}

func writeFile(pathOrFile string, pathTemlate string, ch chan tmpl.Result) error {
	var name string
	var ext string
	var dir string
	stat, err := os.Stat(pathOrFile)

	if err == nil && stat.IsDir() == true {
		dir = pathOrFile
		_, name = filepath.Split(pathTemlate)
		name = strings.Split(name, ".")[0]
		ext = filepath.Ext(pathTemlate)
	} else {
		dir, name = filepath.Split(pathOrFile)
		name = strings.Split(name, ".")[0]
		ext = filepath.Ext(pathOrFile)
	}
	for res := range ch {
		if res.Err != nil {
			return res.Err
		}
		idx := fmt.Sprintf("%0"+strconv.Itoa(int(math.Floor(math.Log10(float64(res.Total)))))+"d", res.No)
		if res.Name != "" {
			name = res.Name + ext
		} else {
			name = name + "_" + idx + ext
		}
		ioutil.WriteFile(filepath.Join(dir, name), []byte(res.Str), os.ModePerm)
	}
	return nil
}
