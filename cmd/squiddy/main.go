package main

import (
	"fmt"
	"os"

	"tankteksoftware.com/squiddy/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		utils.PrintHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "ask":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: squiddy ask <your question>")
			os.Exit(1)
		}
		question := utils.JoinArgs(os.Args[2:])
		fmt.Println("Q: ", question)
		// TODO: Run ask
	case "version":
		fmt.Println("squiddy v1.0.0")
	case "help", "--help", "-h":
		utils.PrintHelp()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		utils.PrintHelp()
		os.Exit(1)
	}
}
