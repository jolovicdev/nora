package main

import (
	"fmt"
	"os"

	"github.com/jolovicdev/nora/internal/app"
)


func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: nora <command> [arguments]")
        fmt.Println("Commands:")
        fmt.Println("  init                  - Initialize a new story")
        fmt.Println("  prepare <files...>    - Prepare files for snapshot")
        fmt.Println("  capture <message>     - Create a new snapshot")
        fmt.Println("  recall <snapshot-id>  - View snapshot details")
        fmt.Println("  diff <file>          - Show changes in prepared file")
        os.Exit(1)
    }

    app := app.New(".nora")
    var err error

    switch os.Args[1] {
    case "init":
        err = app.Initialize()
    case "prepare":
        if len(os.Args) < 3 {
            fmt.Println("Usage: nora prepare <files...>")
            os.Exit(1)
        }
        err = app.PrepareFiles(os.Args[2:])
    case "forget":
        if len(os.Args) < 3 {
            fmt.Println("Usage: nora forget <files...>")
            os.Exit(1)
        }
        err = app.Forget(os.Args[2:])
    case "capture":
        if len(os.Args) < 3 {
            fmt.Println("Usage: nora capture <message>")
            os.Exit(1)
        }
        err = app.CreateSnapshot(os.Args[2])
    case "recall":
        if len(os.Args) < 3 {
            fmt.Println("Usage: nora recall <snapshot-id>")
            os.Exit(1)
        }
        err = app.RecallSnapshot(os.Args[2])
    case "diff":
        if len(os.Args) < 3 {
            fmt.Println("Usage: nora diff <file>")
            os.Exit(1)
        }
        err = app.ShowDiff(os.Args[2])
    case "status":
        err = app.GetStatus()
        if err != nil {
            fmt.Printf("Error getting status: %v\n", err)
            os.Exit(1)
        }
    default:
        fmt.Printf("Unknown command: %s\n", os.Args[1])
        os.Exit(1)
    }

    if err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}