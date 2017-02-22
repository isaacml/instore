package main

import (
	"bufio"
	"fmt"
	"github.com/isaacml/instore/libs"
	"strings"
)

func main() {
	var line string
	lectura, err := libs.DownloadPage("http://www.streamrus.com/en/index.php", 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	lector := bufio.NewReader(lectura)
	for {
		line, err = lector.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimRight(line, "\n")

		fmt.Println(line)
	}
}
