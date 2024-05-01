package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func InputString(option string) string {
	fmt.Print(option)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}
func InputInt(option string) int {
	fmt.Print(option)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	i, _ := strconv.Atoi(text)
	return i
}

func ClearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//func Login() user.User {
//	fmt.Println("Login")
//
//}

func main() {
	//err := postgres.Migration()
	//if err != nil {
	//	panic(err)
	//}
	//
	//usr, err := user.ReadUser("Diyorbek", "+_+diyor2005+_+")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(usr)

}
