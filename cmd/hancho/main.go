package main

import (
  "os"

  "github.com/daniel-bryant/hancho"
)

func main() {
  command := ""
  if len(os.Args) > 1 {
    command = os.Args[1]
  }

  hancho.HandleCommand(command)
}
