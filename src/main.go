package main

import (
	"fmt"
	"log"
	"os"
	"indite/setup"

	"github.com/tristanisham/colors"
)

func main() {
	args := os.Args[1:]
	if len(args) >= 1 {
		first_arg := args[0]
		switch first_arg {
		case "init":
			err := setup.Init(".")
			if err != nil {
				log.Fatal(err)
			}
		case "new":
			if len(args) > 1 && args[1] != "." {
				err := setup.Init(args[1])
				if err != nil {
					log.Fatal(err)
				}
			} else {
				cwd, _ := os.Getwd()
				fmt.Printf("%s destination %s already exists.\n\nUse `write init` to initialize the directory.\n", colors.As("WARNING:", colors.BgDarkRed, colors.White), cwd)
			}
		case "build", "b":
			blog, err := setup.Build()
			if err != nil {
				log.Fatal(err)
			}

			

			if err := blog.Dist(); err != nil {
				log.Fatal(err)
			}
		case "dev":
			blog, err := setup.Build()
			if err != nil {
				log.Fatal(err)
			}

			if err := blog.Dev(); err != nil {
				log.Fatal(err)
			}
		}
	}
}