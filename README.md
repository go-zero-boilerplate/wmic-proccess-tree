# wmic-proccess-tree
Windows only - use WMIC to get the process tree


## Quick Start

Install with `go get -u github.com/go-zero-boilerplate/wmic-proccess-tree`

```
import (
    "fmt"
    "log"

    "github.com/go-zero-boilerplate/wmic-proccess-tree/process"
)

func main() {
    procId := 123 //The process id you want to get the tree for
    procTree, err := process.LoadProcessTree(procId)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(fmt.Sprintf("%+v", procTree))
}
```