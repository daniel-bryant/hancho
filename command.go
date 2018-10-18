package main

import (
  "fmt"
)

func handleCommand(command string) {
  const unknown = "hancho: '%s' is not a hancho command. See 'hancho help'.\n"

  const usage = `usage: hancho <command>

With <command> being one of the following:

  start, s   Start proxy server (default)
  help, h    Show help`

  switch command {
  case "start", "s", "up", "":
    handleStartCommand()
  case "stop", "down":
    handleStopCommand()
  case "help", "h":
    fmt.Println(usage)
  default:
    fmt.Printf(unknown, command)
  }
}
