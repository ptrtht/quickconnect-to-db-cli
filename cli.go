package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	cmdrc := getCmdrcPath()

	command := createCommandString(reader)

	fmt.Print("This will be appended to your ~zshrc\n", command)
	fmt.Print("\n Are you sure you want to continue? (y/n) ")

	sure, err := reader.ReadString('\n')
	checkErr(err)
	if sure[0:len(sure)-1] != "y" {
		panic("Aborted")
	}
	file, err := os.OpenFile(cmdrc, os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()

	length, err := io.WriteString(file, command)
	checkErr(err)

	fmt.Printf("%d chars written\n", length)
	fmt.Println("Done, please refresh your shell")
}

func createCommandString(reader *bufio.Reader) string {
	// TODO: add validation for values.
	ip := promptForValue(reader, "IP: ", "NaN")
	dbname := promptForValue(reader, "DB Name: ", "postgres")
	user := promptForValue(reader, "User: ", "postgres")
	port := promptForValue(reader, "Port: ", "5432")
	password := promptForValue(reader, "Password: ", "NaN")

	commandName := promptForValue(reader, "Command Name: ", "NaN")

	return fmt.Sprintf(
		`
# Next function is added by cli
function %s()
{
	echo "%s" 
	psql -h %s -p %s -U %s -d %s -W
}

`, commandName, password, ip, port, user, dbname)
}

func promptForValue(reader *bufio.Reader, flag string, def string) string {
	fmt.Print(fmt.Sprintf("Please enter a value for: %s (default: %s): ", flag, def))
	input, err := reader.ReadString('\n')
	checkErr(err)
	input = input[:len(input)-1]

	if input == "" {
		if def == "NaN" || def == "" {
			panic(fmt.Sprintf("No default for %s, please enter a value", flag))
		}
		return def
	}
	return input
}

func getHome() string {
	home, err := os.UserHomeDir()
	checkErr(err)
	return home
}

func getCmdrcPath() string {
	home := getHome()
	return fmt.Sprintf("%s/.zshrc", home)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
