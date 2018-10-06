package main

import (
  "os"
)

func main() {
  command := ""
  if len(os.Args) > 1 {
    command = os.Args[1]
  }

  handleCommand(command)
}
