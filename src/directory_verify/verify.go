package directory_verify

import (
	"fmt"
	"log"
	"os"
)

func Verify_directory() {
	if _, err := os.Stat("dados/"); err == nil {
		fmt.Printf("File exists\n")
	} else {
		if err := os.MkdirAll("dados", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("certificados/"); err == nil {
		fmt.Printf("File exists\n")
	} else {
		if err := os.MkdirAll("certificados", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("logs/"); err == nil {
		fmt.Printf("File exists\n")
	} else {
		if err := os.MkdirAll("logs", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	if _, err := os.Stat("stage/"); err == nil {
		fmt.Printf("File exists\n")
	} else {
		if err := os.MkdirAll("stage", os.ModePerm); err != nil {
			log.Fatal(err)
		}
		if err := os.MkdirAll("stage/done", os.ModePerm); err != nil {
			log.Fatal(err)
		}
		if err := os.MkdirAll("stage/queue", os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

}
