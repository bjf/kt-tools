package main

import (
    "flag"
    "fmt"
    "log"
    "github.com/bjf/lpad"
    "strconv"
    "github.com/bjf/ktCore"
    )

func main() {
    // Basic command line arg parsing
    //
    flag.Parse()
    args := flag.Args()
    target := "0"
    if len(args) > 0 {
        target = flag.Args()[0]
    } else {
        fmt.Printf("  Error: Failed to specify a bug number.\n")
        return
    }

    // Fetch a launchpad bug
    //
    {
        oath := &lpad.OAuth{Anonymous: true, Consumer: "bradf-go"}
        root, err := lpad.Login(lpad.Production, oath)
        if err != nil {
            log.Fatal(err)
        }

        //ubuntu_distro, err := root.Distro("ubuntu")
        //if err != nil {
        //    log.Fatal(err)
        //}


        target_id, _ := strconv.Atoi(target)
        bug, err := root.Bug(target_id)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Printf("\n    %d: %s\n\n", bug.Id(), bug.Title())

        owner, _ := bug.Owner()
        fmt.Printf("                 Owner: %s\n", owner.DisplayName())

        fmt.Printf("               Created: %s\n", bug.DateCreated())
        fmt.Printf("          Last Message: %s\n", bug.DateLastMessage())
        fmt.Printf("          Last Updated: %s\n", bug.DateLastUpdated())
        fmt.Printf("            Is Private: %t\n", bug.Private())
        fmt.Printf("   Is Security Related: %t\n", bug.SecurityRelated())
        //fmt.Printf("             Duplicate: %s\n", bug.DuplicateOf())
        fmt.Printf("                  Heat: %d\n", bug.Heat())
        fmt.Printf("          Latest Patch: %s\n", bug.DateLatestPatchUploaded())
        fmt.Printf("          Is Expirable: %t\n", bug.IsExpirable())
        // fmt.Printf("           Series Name: %s\n", (bug.series[0]))
        // fmt.Printf("        Series Version: %s\n", (bug.series[1]))
        // fmt.Printf("          Problem Type: %s\n", (bug.problem_type))

        fmt.Print("\n")
        fmt.Print("        Description:\n")
        fmt.Print("        -----------------------------------------------------------------------------------\n")
        fmt.Print(bug.Description())

        fmt.Print("\n")
        fmt.Print("        Tags:\n")
        fmt.Print("        -----------------------------------------------------------------------------------\n")
        tags := bug.Tags()
        if len(tags) > 0 {
            fmt.Printf("            ")
            for i,element := range tags {
                fmt.Printf("%s", element)
                if i < len(tags) - 1 {
                    fmt.Printf(",")
                }
            }
            fmt.Print("\n")
        } else {
            fmt.Printf("            <empty>\n")
        }

        tasks, err := bug.Tasks()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Print("\n")
        fmt.Print("        Tasks:\n")
        fmt.Print("        -----------------------------------------------------------------------------------\n")

        tasks.For(func(task *lpad.BugTask) error {
            assignee := "none"
            a, err := task.Assignee()
            if err == nil {
                assignee = a.DisplayName()
            }
            owner := "none"
            o, err := task.Owner()
            if err == nil {
                owner = o.DisplayName()
            }

            fmt.Printf("            %s (%s)\n", task.BugTargetName(), task.BugTargetDisplayName())
            fmt.Printf("                       Status: %-20s  Importance: %-20s  Assignee: %-s\n", task.Status(), task.Importance(), assignee)
            fmt.Printf("                      Created: %s\n", task.DateCreated())
            fmt.Printf("                    Confirmed: %s\n", task.DateConfirmed())
            fmt.Printf("                     Assigned: %s\n", task.DateAssigned())
            fmt.Printf("                       Closed: %s\n", task.DateClosed())
            fmt.Printf("                Fix Committed: %s\n", task.DateFixCommitted())
            fmt.Printf("                 Fix Released: %s\n", task.DateFixReleased())
            fmt.Printf("                  In Progress: %s\n", task.DateInProgress())
            fmt.Printf("                   Incomplete: %s\n", task.DateIncomplete())
            fmt.Printf("                  Left Closed: %s\n", task.DateLeftClosed())
            fmt.Printf("                     Left New: %s\n", task.DateLeftNew())
            fmt.Printf("                      Triaged: %s\n", task.DateTriaged())
            fmt.Printf("                  Is Complete: %t\n", task.IsComplete())
            fmt.Printf("                        Owner: %s\n", owner)
            fmt.Printf("                        Title: %s\n", task.Title())

            fmt.Printf("\n")
            return nil
        })

        for k, _ := range core.DB {
            fmt.Printf("%s\n", k)
        }
    }
}
