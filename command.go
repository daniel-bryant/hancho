package hancho

import (
  "fmt"
  "strings"
)

type Command struct {
  Names []string
  Handle func()
  Description string
}

var availableCommands = []Command{
  {[]string{"up", "u"}, handleStartCommand, "Create and start all docker containers"},
  {[]string{"down", "d"}, handleStopCommand, "Stop and remove all containers"},
  {[]string{"proxies", "p"}, handleProxiesCommand, "List the ports being proxied"},
}

func handleHelpCommand() {
  fmt.Println("Usage: hancho <command>\n")
  fmt.Println("With <command> being one of the following:\n")
  printCommands()
}

func printCommands() {
  max := 0
  for _, command := range availableCommands {
    length := len(strings.Join(command.Names,", "))
    if length > max {
      max = length
    }
  }

  for _, command := range availableCommands {
    names := strings.Join(command.Names,", ")
    padding := strings.Repeat(" ", max - len(names))
    fmt.Println(names, padding, command.Description)
  }
}

func HandleCommand(arg string) {
  helpCommand := Command{[]string{"help", "h"}, handleHelpCommand, "Show help"}
  availableCommands = append(availableCommands, helpCommand)

  for _, command := range availableCommands {
    for _, name := range command.Names {
      if arg == name {
        command.Handle()
        return
      }
    }
  }

  if len(arg) == 0 {
    handleHelpCommand()
  } else {
    fmt.Printf("'%s' is not a hancho command. See 'hancho help'.\n\n", arg)
    printCommands()
  }
}
