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
    regBranchId    = regexp.MustCompile("(.*[bB]ranch_?[iI]d\\s*=\\s*)(.\\d*)(.*)")
    regBranchName  = regexp.MustCompile("(.*[bB]ranch_?[nN]ame\\s*=\\s*[\"'])(.*)([\"'].*)")
    regBuildId     = regexp.MustCompile("(.*[bB]uild_?[iI]d\\s*=\\s*)(\\d*)(.*)")
    regBuildConfig = regexp.MustCompile("(.*[bB]uild_?[cC]onfig\\s*=\\s*[\"'])(.*)([\"'].*)")
)

// config vars
var (
    path        string
    branch      string
    branchId    int
    buildConfig string
    buildId     int
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

        replaced, newLine := replaceInfo(
            regBranchId,
            line,
            fmt.Sprintf("${1}%d${3}", branchId),
        )
        if replaced {
            buffer.WriteString(newLine)
            continue
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

func replaceInfo(regex *regexp.Regexp, line, repl string) (bool, string) {
    if len(regex.FindStringSubmatch(line)) < 3 {
        return false, ""
    }

    return true, regex.ReplaceAllString(line, repl)
}

func parseArgs() {
    flag.StringVar(&path, "path", "", "(required) Path to the source file containing build information")
    flag.StringVar(&branch, "branchName", "", "(required) The branch to label the source file with")
    flag.IntVar(&branchId, "branchId", -1, "(required) The numeric branch id to label the source file with")
    flag.StringVar(&buildConfig, "buildConfig", "Debug", "The build configuration to label the source file with")
    flag.IntVar(&buildId, "buildId", 0, "The BuildID to label the source file with")
    flag.Parse()

    if len(path) < 1 || len(branch) < 1 || branchId == -1 {
        flag.PrintDefaults()
        os.Exit(1)
    }
}
