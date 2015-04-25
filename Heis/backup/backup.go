package backup

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	NUM_BACKUPS = 3
)

// Creates child processes of the main program, and creates new ones if the current dies.
// This is the bottom Fail-Safe state, which will restart the program if necessary.
func CreateBackupsAndListen(sigc chan os.Signal) {
	go func() {
		iter := 0
		for iter < NUM_BACKUPS {
			<-sigc
			iter++
			fmt.Println()
		}
		os.Exit(0)
	}()

	fmt.Println("\x1b[34;1m::: Start Backup Service :::\x1b[0m")

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	for {
		process, err := os.StartProcess(dir+os.Args[0][1:], []string{os.Args[0], "-c"}, &procAttr)
		if err != nil {
			fmt.Printf("\t\x1b[31;1mError\x1b[0m |CreateBackupsAndListen| [%v], continue\n", err)
		}
		_, err = process.Wait()
		if err != nil {
			fmt.Printf("\t\x1b[31;1mError\x1b[0m |CreateBackupsAndListen| [%v], continue\n", err)
		}

		fmt.Println("\t\x1b[31;1mError\x1b[0m |CreateBackupsAndListen| [Child process has exited], create new one")

		time.Sleep(500 * time.Millisecond)
	}
}

// Read any internal orders stored locally on the computer.
func ReadInternalBackup() []int {
	file, err := os.Open("internal.gob")
	if err != nil {
		fmt.Printf("\t\x1b[31;1mError\x1b[0m |1ReadInternalBackup| [%v], returning default\n", err)
		return []int{-1, -1, -1, -1}
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var rinternal []int
	err = decoder.Decode(&rinternal)
	if err != nil {
		fmt.Printf("\x1b\t[31;1mError\x1b[0m |2ReadInternalBackup| [%v], returning default\n", err)
		return []int{-1, -1, -1, -1}
	}
	return rinternal
}

// Writes internal orders locally on the computer.
func WriteInternalBackup(winternal []int) error {
	file, err := os.Create("internal.gob")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(winternal)
	if err != nil {
		return err
	}
	return nil
}
