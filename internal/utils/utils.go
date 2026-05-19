package utils

import "fmt"

func JoinArgs(args []string) string {
	result := ""
	for i, a := range args {
		if i > 0 {
			result += " "
		}
		result += a
	}

	return result
}

func PrintHelp() {
	fmt.Print(`
███████╗     ██████╗     ██╗   ██╗    ██╗    ██████╗     ██████╗     ██╗   ██╗
██╔════╝    ██╔═══██╗    ██║   ██║    ██║    ██╔══██╗    ██╔══██╗    ╚██╗ ██╔╝
███████╗    ██║   ██║    ██║   ██║    ██║    ██║  ██║    ██║  ██║     ╚████╔╝ 
╚════██║    ██║▄▄ ██║    ██║   ██║    ██║    ██║  ██║    ██║  ██║      ╚██╔╝  
███████║    ╚██████╔╝    ╚██████╔╝    ██║    ██████╔╝    ██████╔╝       ██║   
╚══════╝     ╚══▀▀═╝      ╚═════╝     ╚═╝    ╚═════╝     ╚═════╝        ╚═╝   

Your shell command assistant 🦑 

USAGE:
  squiddy ask <question>
  squiddy provider <provider_name>
  squiddy api_key <api_key>
  squiddy version
  squiddy help

EXAMPLES:
  squiddy ask how do I kill a port process?
  squiddy ask how do I start a docker process?
  squiddy ask how do I center a div?
  squiddy ask what option do I use to view my git changes?
	`)
}
