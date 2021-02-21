package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	author  = "starx"
	current = "current"
	target  = "target"
	gbk     = "gbk"
	gb2312  = "gb2312"
	big5    = "big5"
	empty 	= ""
)


func check(err error)  {
	if err != nil{
		panic(err)
	}
}

func isDir(path string)bool{
	if info,err := os.Stat(path); err == nil{
		return info.IsDir()
	}
	return true
}

func detect(path string,formats []string)  bool{
	if len(formats) == 0{
		return true
	}
	for _,f :=range formats{
		ext := filepath.Ext(path)
		if strings.Contains(ext,f){
			return true
		}
	}
	return false
}

func process(src string,eg string,dt string,df []string) {
	var sourcePath string

	if src == current{
		sourcePath,_ = os.Getwd()
	}else {
		sourcePath = src
	}


	if !isDir(dt){
		panic(errors.New("dest path must be directory"))
	}
		var files []string
		if isDir(sourcePath) {
			_ = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
				check(err)
				if !info.IsDir() {
					files = append(files, path)
				}
				return nil
			})
		} else {
			for _,v := range strings.Split(src,","){
				files = append(files, v)
			}
		}

	for i,v := range files{
		//buf := new(bytes.Buffer)
		//sReader := strings.NewReader(src)
		//_,err:= buf.ReadFrom(sReader)
		//check(err)
		if !detect(v,df){
			continue
		}

		data, err := ioutil.ReadFile(v)
		check(err)
		sReader := bytes.NewReader(data)
		var dReader *transform.Reader
		//buf.Bytes()
		switch eg {
		case gbk:
			dReader = transform.NewReader(sReader, simplifiedchinese.GBK.NewDecoder())
		case gb2312:
			dReader = transform.NewReader(sReader, simplifiedchinese.HZGB2312.NewDecoder())
		case big5:
			dReader = transform.NewReader(sReader, traditionalchinese.Big5.NewDecoder())
		}
		out,_ := ioutil.ReadAll(dReader)
		var dest string
		switch dt {
		case current:
			wd,_ := os.Getwd()
			dest = filepath.Join(wd,filepath.Base(v))
		case target:
			dest = v
		default:
			if _, err := os.Stat(dt); err != nil {
				if os.IsNotExist(err) {
					_ = os.MkdirAll(dt,os.ModePerm)
				} else {
					panic(err)
				}
			}
			dest = filepath.Join(dt,filepath.Base(v))
		}
		err = ioutil.WriteFile(dest,out,os.ModePerm)
		check(err)
		fmt.Printf("index: %d source: %s encoding: %s dest: %s\n",i,v,eg,dest)
	}
}

func appRun(c *cli.Context) error{
	argSP := c.String("source")
	argEG := c.String("encoding")
	argDT := c.String("destination")
	argDF := strings.Split(c.String("format"),",")

	if argSP == current {
		argSP,_ = os.Getwd()
	}

	if argDT == current {
		argDT,_ = os.Getwd()
	}
	process(argSP,argEG,argDT,argDF)
	return nil
}

func main() {
	app := &cli.App{
		Name:        "Simple Text Converter",
		HelpName:    "stc",
		ArgsUsage:   "./stc [src_path] [src_encoding] [dest_path]",
		Version:     "v1.0",
		Description: "a simple tool that converts your files text to uft8 format.",
		Flags:       []cli.Flag{&cli.StringFlag{
			Name:     "source",
			Aliases:  []string{"s"},
			Usage:    "source file or dir path",
			Required: false,
			Value:    current,
		},&cli.StringFlag{
			Name:     "format",
			Aliases:  []string{"f"},
			Usage:    "source file[s] format",
			Required: false,
			Value:    empty,
		}, &cli.StringFlag{
			Name:     "encoding",
			Aliases:  []string{"e"},
			Usage:    "source file[s] encoding",
			Required: false,
			Value:    gbk,
		},&cli.StringFlag{
			Name:     "destination",
			Aliases:  []string{"d"},
			Usage:    "file[s] destination path",
			Required: false,
			Value:    current,
		}},
		Action:      appRun,
		Copyright:   author,
	}
	err := app.Run(os.Args)
	check(err)
}
