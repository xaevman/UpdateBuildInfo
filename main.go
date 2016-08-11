package main

import (
    "bufio"
    "bytes"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "regexp"
    "strings"

    "github.com/xaevman/buildinfo"
)

// regexp
var (
    regBranchId    = regexp.MustCompile("BranchId\\s*=\\s*(.*)")
    regBuildId     = regexp.MustCompile("BuildId\\s*=\\s*(.*)")
    regBuildConfig = regexp.MustCompile("BuildConfig\\s*=\\s*(.*)")
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

        matches := regBranchId.FindStringSubmatch(line)
        if len(matches) > 0 {
            // rewrite branch id
            newLine := strings.Replace(matches[0], matches[1], fmt.Sprintf("%d\n", branchId), 1)
            buffer.WriteString(newLine)
            continue
        }
        matches = regBuildConfig.FindStringSubmatch(line)
        if len(matches) > 0 {
            // rewrite build config
            newLine := strings.Replace(matches[0], matches[1], fmt.Sprintf("\"%s\"\n", buildConfig), 1)
            buffer.WriteString(newLine)
            continue
        }
        matches = regBuildId.FindStringSubmatch(line)
        if len(matches) > 0 {
            // rewrite build id
            newLine := strings.Replace(matches[0], matches[1], fmt.Sprintf("%d\n", buildId), 1)
            buffer.WriteString(newLine)
            continue
        }

        // write line as-is
        buffer.WriteString(strings.Replace(line, "\r", "", 1))
    }

    f.Close()

    ioutil.WriteFile(path, buffer.Bytes(), 0660)

    cmd := exec.Command("go", "fmt", path)
    cmd.Run()
}

func parseArgs() {
    flag.StringVar(&path, "path", "buildid.go", "Path to the source file containing build information")
    flag.StringVar(&branch, "branch", "Local", "The branch to label the source file with")
    flag.StringVar(&buildConfig, "buildConfig", "Debug", "The build configuration to label the source file with")
    flag.IntVar(&buildId, "buildId", 0, "The BuildID to label the source file with")
    flag.Parse()

    branchId = buildinfo.GetBranch(branch)
}
