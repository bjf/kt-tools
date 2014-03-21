package main

import (
    "fmt"
    "github.com/bjf/ktCore"
    )

func printInfo(record map[string]interface{}) {
    var ok bool

    if _, ok = record["development"]; ok {
        fmt.Printf("           development : %t\n", record["development"])
    } else {
        fmt.Printf("           development : false\n")
    }
    fmt.Printf("        series version : %s\n", record["series_version"])
    fmt.Printf("           series name : %s\n", record["name"])
    fmt.Printf("        kernel version : %s\n", record["kernel"])
    fmt.Printf("          is supported : %t\n", record["supported"])

    // packages
    //
    fmt.Printf("              packages : ")
    pkgs := record["packages"].([]string)
    if len(pkgs) > 0 {
        for _, p := range pkgs {
            fmt.Printf("%s  ", p)
        }
    } else {
        fmt.Printf("none")
    }
    fmt.Printf("\n")

    // dependent-packages
    //
    fmt.Printf("    dependent packages : ")
    if _, ok = record["dependent-packages"]; ok {
        depPkgs := record["dependent-packages"].(map[string]map[string]string)
        for k1, v1 := range depPkgs {
            fmt.Printf("[%s] ", k1)
            depPkg := v1
            for _, v2 := range depPkg {
                fmt.Printf("%s  ", v2)
            }
        }
        fmt.Printf("\n")
    } else {
        fmt.Printf(" none\n")
    }

    // derivative-packages
    //
    fmt.Printf("   derivative packages : ")
    if _, ok = record["derivative-packages"]; ok {
        derPkgs := record["derivative-packages"].(map[string][]string)
        if len(derPkgs) > 0 {
            for _, p := range derPkgs {
                for _, p3 := range p {
                    fmt.Printf("%s  ", p3)
                }
            }
        } else {
            fmt.Printf("none")
        }
        fmt.Printf("\n")
    } else {
        fmt.Printf(" none\n")
    }
}

func main() {
    var kp = core.NewUbuntuKernelsDB()

    printInfo(kp.IndexByKernelVersion("3.13.0"))
    fmt.Printf("\n")
    printInfo(kp.IndexByKernelVersion("3.11.0"))
    fmt.Printf("\n")
    fmt.Printf("\n")
    printInfo(kp.IndexBySeriesName("precise"))

    // for seriesVersion, _ := range core.DB {
    //     fmt.Printf("%s\n", seriesVersion)
    //     printInfo(core.DB[seriesVersion])
    //     fmt.Printf("\n")
    // }
}
