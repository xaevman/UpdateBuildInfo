package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
)

// regexp
var (
	regBranchIdNumeric = []*regexp.Regexp{
		// cs, go, py
		regexp.MustCompile("(.*[bB][rR][aA][nN][cC][hH]_?[iI][dD]\\s*=\\s*)(.\\d*)(.*)"),
		// json
		regexp.MustCompile("(.*\"[bB][rR][aA][nN][cC][hH]_?[iI][dD]\"\\s*:\\s*)(.\\d*)(.*)"),
		// c
		regexp.MustCompile("(#define [bB][rR][aA][nN][cC][hH]_?[iI][dD]\\s*)(.\\d*)(.*)"),
	}
	regBranchIdStr = []*regexp.Regexp{
		// cs, go, py
		regexp.MustCompile("(.*[bB][rR][aA][nN][cC][hH]_?[iI][dD]\\s*=\\s*[\"'])(.*)([\"'].*)"),
		// json
		regexp.MustCompile("(.*\"[bB][rR][aA][nN][cC][hH]_?[iI][dD]\"\\s*:\\s*\")(.*)(\".*)"),
		// c
		regexp.MustCompile("(#define [bB][rR][aA][nN][cC][hH]_?[iI][dD]\\s*\")(.*)(\".*)"),
	}
	regBranchName = []*regexp.Regexp{
		// cs, go, py
		regexp.MustCompile("(.*[bB][rR][aA][nN][cC][hH]_?[nN][aA][mM][eE]\\s*=\\s*[\"'])(.*)([\"'].*)"),
		// json
		regexp.MustCompile("(.*\"[bB][rR][aA][nN][cC][hH]_?[nN][aA][mM][eE]\"\\s*:\\s*\")(.*)(\".*)"),
	}
	regBuildId = []*regexp.Regexp{
		// cs, go, py
		regexp.MustCompile("(.*[bB][uU][iI][lL][dD]_?[iI][dD]\\s*=\\s*)(\\d*)(.*)"),
		// json
		regexp.MustCompile("(.*\"[bB][uU][iI][lL][dD]_?[iI][dD]\"\\s*:\\s*)(\\d*)(.*)"),
		// c
		regexp.MustCompile("(#define [bB][uU][iI][lL][dD]_?[iI][dD]\\s*)(\\d*)(.*)"),
	}
	regBuildConfig = []*regexp.Regexp{
		// cs, go, py
		regexp.MustCompile("(.*[bB][uU][iI][lL][dD]_?[cC][oO][nN][fF][iI][gG]\\s*=\\s*[\"'])(.*)([\"'].*)"),
		// json
		regexp.MustCompile("(.*\"[bB][uU][iI][lL][dD]_?[cC][oO][nN][fF][iI][gG]\"\\s*:\\s*\")(.*)(\".*)"),
		// c
		regexp.MustCompile("(#define [bB][uU][iI][lL][dD]_?[cC][oO][nN][fF][iI][gG]\\s*\")(.*)(\".*)"),
	}
)

// config vars
var (
	path              string
	branch            string
	branchId          string
	buildConfig       string
	buildId           int
	useStringBranchId bool
)

func main() {
	parseArgs()
	labelBuild()
}

func labelBuild() {
	var buffer bytes.Buffer

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}

		replaced := false
		newLine := line

		if useStringBranchId {
			replaced, newLine = replaceInfo(
				regBranchIdStr,
				line,
				fmt.Sprintf("${1}%s${3}", branchId),
			)
			if replaced {
				buffer.WriteString(newLine)
				continue
			}
		} else {
			replaced, newLine = replaceInfo(
				regBranchIdNumeric,
				line,
				fmt.Sprintf("${1}%s${3}", branchId),
			)
			if replaced {
				buffer.WriteString(newLine)
				continue
			}
		}

		replaced, newLine = replaceInfo(
			regBranchName,
			line,
			fmt.Sprintf("${1}%s${3}", branch),
		)
		if replaced {
			buffer.WriteString(newLine)
			continue
		}

		replaced, newLine = replaceInfo(
			regBuildConfig,
			line,
			fmt.Sprintf("${1}%s${3}", buildConfig),
		)
		if replaced {
			buffer.WriteString(newLine)
			continue
		}

		replaced, newLine = replaceInfo(
			regBuildId,
			line,
			fmt.Sprintf("${1}%d${3}", buildId),
		)
		if replaced {
			buffer.WriteString(newLine)
			continue
		}

		// write line as-is
		buffer.WriteString(line)
	}

	f.Close()

	ioutil.WriteFile(path, buffer.Bytes(), 0660)
}

func replaceInfo(regexList []*regexp.Regexp, line, repl string) (bool, string) {
	for _, regex := range regexList {
		if len(regex.FindStringSubmatch(line)) == 4 {
			return true, regex.ReplaceAllString(line, repl)
		}
	}

	return false, ""
}

func parseArgs() {
	flag.StringVar(&path, "path", "", "(required) Path to the source file containing build information")
	flag.StringVar(&branch, "branchName", "", "(required) The branch to label the source file with")
	flag.StringVar(&branchId, "branchId", "-1", "(required) The numeric branch id to label the source file with")
	flag.StringVar(&buildConfig, "buildConfig", "Debug", "The build configuration to label the source file with")
	flag.IntVar(&buildId, "buildId", 0, "The BuildID to label the source file with")
	flag.BoolVar(&useStringBranchId, "useStringBranchId", false, "Whether to interpret branch id as a string or number")
	flag.Parse()

	if useStringBranchId {
		fmt.Println("Using string branch IDs")
	}

	if len(path) < 1 || len(branch) < 1 || branchId == "-1" {
		flag.PrintDefaults()
		os.Exit(1)
	}
}
