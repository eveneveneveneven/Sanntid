package backup

import (
	"encoding/gob"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	NUM_BACKUPS = 5
)

func CreateBackupsAndListen(sigc chan os.Signal) {
	go func() {
		iter := 0
		for iter < NUM_BACKUPS {
			s := <-sigc
			iter++
			fmt.Println("\n", s)
			fmt.Println("Number of backups:", iter)
		}
		os.Exit(0)
	}()
	fmt.Println("Start backup process!")
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	var procAttr os.ProcAttr
	procAttr.Files = []*os.File{nil, os.Stdout, os.Stderr}
	for {
		process, err := os.StartProcess(dir+os.Args[0][1:], []string{os.Args[0], "-c"}, &procAttr)
		if err != nil {
			fmt.Printf("\x1b[31;1mError\x1b[0m |CreateBackupsAndListen| [%v], continue\n", err)
		}
		_, err = process.Wait()
		if err != nil {
			fmt.Printf("\x1b[31;1mError\x1b[0m |CreateBackupsAndListen| [%v], continue\n", err)
		}
		fmt.Println("Child process has exited!")
		time.Sleep(500 * time.Millisecond)
	}
}

func ReadInternalBackup() []int {
	file, err := os.Open("internal.gob")
	if err != nil {
		fmt.Printf("\x1b[31;1mError\x1b[0m |1ReadInternalBackup| [%v], returning default\n", err)
		return []int{-1, -1, -1, -1}
	}
	defer file.Close()
	decoder := gob.NewDecoder(file)
	var rinternal []int
	err = decoder.Decode(&rinternal)
	if err != nil {
		fmt.Printf("\x1b[31;1mError\x1b[0m |2ReadInternalBackup| [%v], returning default\n", err)
		return []int{-1, -1, -1, -1}
	}
	return rinternal
}

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
