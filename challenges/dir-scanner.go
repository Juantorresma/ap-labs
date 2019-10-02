//Primer parcial A01227885

package main

import (
	"fmt"
	"os"
	"path/filepath"
)

//Aqui hago la magía
func scanDir(dir string) error {

    //inicializamos nuestras variables en 0
	directories := 0
	symbolicLinks := 0
	devices := 0
	sockets := 0
	otherFiles := 0

    //checamos cada una y en caso que sean del tipo, agregamos una
	err := filepath.Walk(dir, func(dir string, fileI os.FileInfo, err error) error {
		fileI, err = os.Lstat(dir)
		if fileI.Mode()&os.ModeDir != 0 {
			directories = directories + 1
		} else if fileI.Mode()&os.ModeSymlink != 0 {
			symbolicLinks = symbolicLinks + 1
		} else if fileI.Mode()&os.ModeDevice != 0 {
			devices = devices + 1
		} else if fileI.Mode()&os.ModeSocket != 0 {
			sockets = sockets + 1
		} else {
			otherFiles = otherFiles + 1
		}
		return err
	})

    //im`primimos los resultados
	fmt.Println("Directory Scanner Tool")
	fmt.Println("+-------------------------+------+")
	fmt.Println("| Path                    |", dir, "|")
	fmt.Println("+-------------------------+------+")
	fmt.Println("| Directories             | ", directories, " |")
	fmt.Println("| Symbolic Links          | ", symLinks, " |")
	fmt.Println("| Devices                 | ", devices, " |")
	fmt.Println("| Sockets                 | ", sockets, " |")
	fmt.Println("| Other files             | ", other, " |")
	fmt.Println("+-------------------------+------+")
	return err
}

//Aqui solo verificaremos que nos den el argumento de directorio, sino mandaremos warning
func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: ./dir-scan <directory>")
		os.Exit(1)
	}
	scanDir(os.Args[1])
}
