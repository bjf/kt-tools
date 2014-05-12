// This is a re-implementation of "maint-modify-patch" which is written in
// python.
//

package main

import (
    "fmt"
    "os"
    "os/exec"
    "bytes"
    "regexp"
    "strings"
    "github.com/bjf/ktCore"
    "path/filepath"
    "github.com/jessevdk/go-flags"
    "github.com/daviddengcn/go-colortext"
    )


var args struct {
    Color         bool   `short:"c" long:"color"      description:"Make the output pretty."`
}

// hasString
//     Does the given slice contain the specified string?
//
func hasString(slice []string, target string) bool {
    var retval bool = false
    for _, str := range slice {
        if str == target {
            retval = true
            break
        }
    }

    return retval
}

// processFilename
//
func processFileName(fileName string) (string, string) {
    var (
        name        string
        version     string

        // Very simple view of what a changes filename looks like.
        //
        changesRC            = regexp.MustCompile("(.*)_(.*)\\.dsc")
    )

    result := changesRC.FindStringSubmatch(fileName)
    name    = result[1]
    version = result[2]

    return name, version
}

// kernelVersion
//
func kernelVersion(version string) string {
    var (
        versionRC = regexp.MustCompile("^([0-9]+\\.[0-9]+\\.[0-9]+)[-\\.]([0-9]+)\\.([0-9]+)([~\\S]*)$")
        kdb = core.NewUbuntuKernelsDB()
    )

    result := versionRC.FindStringSubmatch(version)
    return kdb.IndexByKernelVersion(result[1])["name"].(string)
}

// main
//    This is where the bugs start.
//
func main() {
    var (
        packageName    string
        packageVersion string
        buf = new(bytes.Buffer)
        eol = []byte{'\n'}
    )

    // Basic command line arg parsing
    //
    fileNames, err := flags.Parse(&args)
    if err != nil {
        os.Exit(1)
    }

    if len(fileNames) == 0 {
        fmt.Printf("  Error: No files were specified on the command line.\n")
        os.Exit(1)
    }

    for _, fid := range fileNames {
        packageName, packageVersion = processFileName(fid)
        series := kernelVersion(packageVersion)
        cmd := exec.Command("pull-lp-source", "-d", packageName, series)
        out, err := cmd.CombinedOutput()
        if err != nil {
            println(err.Error())
            print(string(out))
            return
        }

        files, _ := filepath.Glob(packageName + "*.dsc")

        debdiff := exec.Command("debdiff", files[0], files[1])
        filterdiff := exec.Command("filterdiff", "-X", "~/.filterdiff.filters")

        filterdiff.Stdin, _ = debdiff.StdoutPipe()
        filterdiff.Stdout = buf
        _ = filterdiff.Start()
        _ = debdiff.Run()
        _= filterdiff.Wait()

        if args.Color {
            for _, line := range bytes.Split(buf.Bytes(), eol) {
                str := string(line)
                ct.ResetColor()
                if strings.HasPrefix(str, "+++") || strings.HasPrefix(str, "---") {
                    ct.ChangeColor(ct.Yellow, false, ct.None, false)
                } else if strings.HasPrefix(str, "+") {
                    ct.ChangeColor(ct.Green, false, ct.None, false)
                } else if strings.HasPrefix(str, "-") {
                    ct.ChangeColor(ct.Red, false, ct.None, false)
                }
                print(str)
                print("\n")
            }
        } else {
            buf.WriteTo(os.Stdout)
        }
    }
}

