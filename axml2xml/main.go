package main

import (
    "os"
    "fmt"
    "io"
    "encoding/xml"
    "flag"

    "binxml"
    "strings"
)

func main() {
    isApk := flag.Bool("a", false, "The input file is an apk")

    flag.Parse()

    if len(os.Args) != 2 {
        fmt.Printf("%s INPUT\n", os.Args[0])
        os.Exit(1)
    }

    var r io.Reader
    input := os.Args[1]

    if strings.HasSuffix(input, ".apk") {
        *isApk = true
    }

    if input == "-" {
        r = os.Stdin
    } else if *isApk {
        zr, err := binxml.OpenFileInZip(input, "AndroidManifest.xml")
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        defer zr.Close()
        r = zr
    } else {
        f, err := os.Open(input)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        defer f.Close()
        r = f
    }

    enc := xml.NewEncoder(os.Stdout)
    enc.Indent("", "    ")

    err := binxml.Parse(r, enc)
    fmt.Println()
    if err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
