package main

import (
"encoding/json"
"fmt"
"io"
"os"
)

type Request struct {
A int `json:"a"`
B int `json:"b"`
}

func main(){
input, _ := io.ReadAll(os.Stdin)

var req Request

json.Unmarshal(input, &req)

sum := req.A + req.B

fmt.Println(sum)
}