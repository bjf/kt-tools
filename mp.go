// This is a re-implementation of "maint-modify-patch" which is written in
// python.
//

package main

import (
    "fmt"
    "github.com/jessevdk/go-flags"
    "strings"
    "os"
    )


var args struct {
    Acks []string `long:"ack" description:"An irc nickname which will be expanded to full a email address and added as an 'Acked-by:' line."`
    Sobs []string `long:"sob" description:"An irc nickname which will be expanded to full a email address and added as an 'Signed-off-by:' line."`
    Bugs []string `long:"bug" description:"A Launchpad bug-id which will be used to generate Buglink URL lines."`
    CPs  []string `long:"cp"  description:"The SHA1 of a commit that will make up 'cherry-pick' lines in the signed-off-by block of the commit."`
    BPs  []string `long:"bp"  description:"The SHA1 of a commit that will make up 'back-port' lines in the signed-off-by block of the commit."`
    CVE    string `long:"cve" description:"A CVE number which will be inserted into the commit body."`
}

// main
//    This is where the bugs start.
//
func main() {
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

    fmt.Printf(" acks: %s\n", strings.Join(args.Acks, " "))
    fmt.Printf(" sobs: %s\n", strings.Join(args.Sobs, " "))
    fmt.Printf("  cps: %s\n", strings.Join(args.CPs, " "))
    fmt.Printf("  bps: %s\n", strings.Join(args.BPs, " "))
    fmt.Printf(" bugs: %s\n", strings.Join(args.Bugs, " "))
    fmt.Printf("  cve: %s\n", args.CVE)
    fmt.Printf("files: %s\n", strings.Join(fileNames, " "))
}

