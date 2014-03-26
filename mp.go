// This is a re-implementation of "maint-modify-patch" which is written in
// python.
//

package main

import (
    "fmt"
    "strings"
    "os"
    "bufio"
    "github.com/jessevdk/go-flags"
    "io/ioutil"
    "os/user"
    "path"
    "encoding/json"
    "regexp"
    )


var args struct {
    Acks        []string `long:"ack"          description:"An irc nickname which will be expanded to full a email address and added as an 'Acked-by:' line."`
    Sobs        []string `long:"sob"          description:"An irc nickname which will be expanded to full a email address and added as an 'Signed-off-by:' line."`
    Bugs        []string `long:"bug"          description:"A Launchpad bug-id which will be used to generate Buglink URL lines."`
    CPs         []string `long:"cp"           description:"The SHA1 of a commit that will make up 'cherry-pick' lines in the signed-off-by block of the commit."`
    BPs         []string `long:"bp"           description:"The SHA1 of a commit that will make up 'back-port' lines in the signed-off-by block of the commit."`
    CVE           string `long:"cve"          description:"A CVE number which will be inserted into the commit body."`
    ListAliases   bool   `long:"list-aliases" description:"Print out all the aliases found in the aliases file."`
}


var aliases     map[string]interface{}

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

// sobBlock
//
func sobBlock(existingAcks []string, existingSobs []string, existingCps []string, existingBps []string) (string) {
    var retval string

    for _, ack := range args.Acks {
        if !hasString(existingAcks, ack) {
            retval += fmt.Sprintf("Acked-by: %s\n", aliases[ack])
        }
    }

    for _, bp := range args.BPs {
        if !hasString(existingBps, bp) {
            retval += fmt.Sprintf("(backported from commit %s upstream)\n", bp)
        }
    }

    for _, cp := range args.CPs {
        if !hasString(existingCps, cp) {
            retval += fmt.Sprintf("(cherry picked from commit %s)\n", cp)
        }
    }

    for _, sob := range args.Sobs {
        if !hasString(existingSobs, sob) {
            retval += fmt.Sprintf("Signed-off-by: %s\n", aliases[sob])
        }
    }

    return retval
}

// modify
//    Open a single file; modify it's contents as they are copied to a temporary file;
//    mv the temporary file to the origian file.
//
func modify(fid string) (error) {
    var (
        justCopy                       bool = false
        looking4SOBInsertionPoint      bool = false
        looking4SubjectLine            bool = true
        looking4BuglinkInsertionPoint  bool = false
        subjectLine                    bool = false
        cveInsertionPoint              bool = false
        buglinkInsertionPoint          bool = false
        sobInsertionPoint              bool = false
        existingCVEs                   []string
        existingBugIds                 []string
        existingAcks                   []string
        existingSobs                   []string
        existingCps                    []string
        existingBps                    []string
        buglinkBaseUrl                 string = "http://bugs.launchpad.net/bugs/"
        cpRC                                  = regexp.MustCompile("cherry picked from commit ([0-9a-zA-Z]+)")
        bpRC                                  = regexp.MustCompile("backported from commit ([0-9a-zA-Z]+) upstream")
    )

    // Open the input file, return the error if there is one.
    //
    inputFile, err := os.Open(fid)
    if err != nil {
        return err
    }
    defer inputFile.Close()

    // Open the temp. file, return the error if there is one.
    //
    dst, err := ioutil.TempFile("./", "mp__")
    if err != nil {
        return err
    }
    defer dst.Close()

    scanner := bufio.NewScanner(inputFile)
    for scanner.Scan() {
        line := scanner.Text()

        if justCopy {
            dst.WriteString(line)
            dst.WriteString("\n")
            continue
        }

        // If we are looking for the Sob insertion point then we've handled
        // all the other cases and we just need to find the sob section.
        //
        if looking4SOBInsertionPoint {
            looking4SOBInsertionPoint = true
            if line == "---" {
                sobInsertionPoint = true
                looking4SOBInsertionPoint = false
            } else {
                if strings.Contains(line, "Acked-by:") {
                    id := strings.Replace(line, "Acked-by:", "", -1)
                    if !hasString(existingAcks, id) {
                        existingAcks = append(existingAcks, id)
                    }
                }

                if strings.Contains(line, "Signed-off-by:") {
                    id := strings.Replace(line, "Signed-off-by:", "", -1)
                    if !hasString(existingSobs, id) {
                        existingSobs = append(existingSobs, id)
                    }
                }

                if strings.Contains(line, "cherry picked") {
                    result := cpRC.FindStringSubmatch(line)
                    existingCps = append(existingCps, result[1])
                }

                if strings.Contains(line, "backported from") {
                    result := bpRC.FindStringSubmatch(line)
                    existingBps = append(existingBps, result[1])
                }
            }
        }

        if sobInsertionPoint {
            dst.WriteString(sobBlock(existingAcks, existingSobs, existingCps, existingBps))
            sobInsertionPoint = false
            justCopy = true
        }

        // After the first blank line after the subject line is where we
        // want to insert our CVE lines if we need to insert any.
        //
        if cveInsertionPoint {
            cveInsertionPoint = true
            if strings.Contains(line, "CVE-") {
                cve := strings.Replace(line, "CVE-", "", -1)
                existingCVEs = append(existingCVEs, cve)
            } else {
                // Add the CVE id here.
                //
                if args.CVE != "" {
                    if !hasString(existingCVEs, args.CVE) {
                        dst.WriteString("CVE-")
                        dst.WriteString(args.CVE)
                        dst.WriteString("\n")
                        dst.WriteString("\n") // One blank line after the CVE line (this assumes there is only one CVE)
                    }
                }
                cveInsertionPoint = false
                looking4BuglinkInsertionPoint = true

                // We don't know at this point if we are going to insert a Buglink
                // so we can't write out the current line of text.
            }
        }

        // After the first blank line after the CVE lines is where the Buglinks are to be
        // inserted.
        //
        if looking4BuglinkInsertionPoint {
            if line != "" {
                looking4BuglinkInsertionPoint = false
                buglinkInsertionPoint = true
            }
        }

        if buglinkInsertionPoint {
            buglinkInsertionPoint = true
            // Just like the CVEs we skip past any existing BugLink lines and build a list of existing
            // buglinks so we don't duplicate any.
            //
            if strings.Contains(line, "BugLink:") {
                s := strings.Split(line, "/")
                id := s[len(s)-1]
                existingBugIds = append(existingBugIds, id)
            } else {
                if len(args.Bugs) > 0 {
                    for _, id := range args.Bugs {
                        if !hasString(existingBugIds, id) {
                            dst.WriteString(fmt.Sprintf("BugLink: %s%s\n", buglinkBaseUrl, id))
                        }
                    }
                    dst.WriteString("\n") // One blank line after the BugLink line
                }
                buglinkInsertionPoint = false
                looking4SOBInsertionPoint = true
            }
        }

        // Once we've found the subject line, we look for the first blank line after it.
        //
        if subjectLine {
            if line == "" {
                cveInsertionPoint = true
                subjectLine = false
            }
        }

        // All modificatins that we make are made after the subject line, therefore that's
        // the first thing we look for.
        //
        if looking4SubjectLine {
            if strings.Contains(line, "Subject:") {
                subjectLine = true
                looking4SubjectLine = false
            }
        }

        dst.WriteString(line)
        dst.WriteString("\n")
    }

    // If the scanner encountered an error, return it.
    //
    if err := scanner.Err(); err != nil {
        return err
    }

    os.Rename(dst.Name(), inputFile.Name())
    return nil
}

// main
//    This is where the bugs start.
//
func main() {
    // Load the aliases, JSON, file.
    //
    usr, err := user.Current()
    if err != nil {
        fmt.Printf("  Unable to determine the current user.\n")
        os.Exit(1)
    }
    aliasesPath := path.Join(usr.HomeDir, ".sob-aliases")

    file, e := ioutil.ReadFile(aliasesPath)
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }

    if err := json.Unmarshal(file, &aliases); err != nil {
        fmt.Printf("    File error: %v\n", err)
        os.Exit(1)
    }

    // Basic command line arg parsing
    //
    fileNames, err := flags.Parse(&args)
    if err != nil {
        os.Exit(1)
    }

    if args.ListAliases {
        fmt.Printf("list aliases\n")
        for k, v := range aliases {
            fmt.Printf("%20s : %s\n", k, v)
        }
        os.Exit(1)
    }

    if len(fileNames) == 0 {
        fmt.Printf("  Error: No files were specified on the command line.\n")
        os.Exit(1)
    }

    for _, fid := range fileNames {
        err := modify(fid)
        if err != nil {
            fmt.Printf("%s", err)
        }
    }
}

