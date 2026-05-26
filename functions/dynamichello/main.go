package main

import (
"fmt"
"io"
"os"
)

func main(){
_, _ = io.ReadAll(os.Stdin)
fmt.Println("Hello from dynamically deployed function")
}