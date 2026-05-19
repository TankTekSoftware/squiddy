package main

import (
	"fmt"
	"os"

	"tankteksoftware.com/squiddy/internal/ask"
	"tankteksoftware.com/squiddy/internal/config"
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
		if err := ask.Run(question); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "api_key":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: squiddy api_key <your-api-key>")
			os.Exit(1)
		}
		if err := config.SetAPIKey(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		path, _ := config.Path()
		fmt.Println("API key saved to", path)
	case "provider":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: squiddy provider <anthropic|openai>")
			os.Exit(1)
		}
		if err := config.SetProvider(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		path, _ := config.Path()
		fmt.Println("Provider saved to", path)
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
