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

	var defaultEncoding string
	var output string

	if runtime.GOOS == "windows" {
		defaultEncoding = "sjis"
	} else {
		defaultEncoding = "utf-8"
	}
	flag.StringVar(&encodeTemplate, "te", defaultEncoding, "template encoding")
	flag.StringVar(&encodeCsv, "ce", defaultEncoding, "csv encoding")
	flag.StringVar(&pathTemlate, "t", "", "template encoding")
	flag.StringVar(&pathCsv, "c", "", "csv encoding")
	flag.StringVar(&output, "o", "./", "putput path or file")

	flag.Parse()

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

	result, err := tmpl.Generate(templateStr, csvStr)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	err = writeFiles(output, result, pathTemlate)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
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

func writeFiles(pathOrFile string, strs []string, pathTemlate string) error {
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
	keta := math.Floor(math.Log10(float64(len(strs)))) + 1
	for i := range strs {
		idx := fmt.Sprintf("%0"+fmt.Sprint(keta)+"d", i+1)
		ioutil.WriteFile(filepath.Join(dir, name+"_"+idx+ext), []byte(strs[i]), os.ModePerm)
	}
	return nil
}
