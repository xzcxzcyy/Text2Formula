package main

import (
	"fmt"
	"os/exec"
)
var (
	path = "/home/ubuntu/tex2png/"
)
func main() {
	cmd := exec.Command("python3", path + "main.py", path, "x")
	cmd.Run()
	fmt.Println(cmd.Args)
	cmd2 := exec.Command("cairosvg", path+ "s2.svg","-o s2.png")
	cmd2.Run()
	fmt.Println(cmd2.Args)
}
