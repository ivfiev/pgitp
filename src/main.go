package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

func main() {
	err := checkDeps()
	if err != nil {
		log.Fatalln(err)
	}
	if len(os.Args) < 4 {
		showHelp()
		return
	}
	ix := 1
	if runtime.GOOS == "android" {
		ix = 2
	}
	cmd := os.Args[ix]
	local := os.Args[ix+1]
	remote := os.Args[ix+2]
	switch cmd {
	case "push":
		err = putDir(local, remote)
	case "pull":
		err = getDir(local, remote)
	default:
		fmt.Printf("invalid command [%s]", cmd)
		showHelp()
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
}
