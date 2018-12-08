package main

import "github.com/envzo/zorm/cmd"

func main() {
	if err := cmd.Exec(); err != nil {
		panic(err)
	}
}
