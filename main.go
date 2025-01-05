/*
Copyright (C) 2025 adzsx
*/

package main

import (
	"github.com/adzsx/gwire/cmd"
	"github.com/adzsx/gwire/internal/scan"
)

func main() {
	print(scan.Ping("192.168.178.80"))
	cmd.Execute()
}
