# wmic-proccess-tree
Windows only - use WMIC to get the process tree


## Quick Start

Install with `go get -u github.com/go-zero-boilerplate/wmic-proccess-tree`

```
package main

import (
    "fmt"
    "log"
    "github.com/golang-devops/exec-logger/process_tree"
)

func main() {
    procId := 123 //The process id you want to get the tree for
    procTree, err := process_tree.LoadProcessTree(procId)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(fmt.Sprintf("%+v", procTree))
}
```