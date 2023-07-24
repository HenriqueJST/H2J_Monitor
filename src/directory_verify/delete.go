package directory_verify

import (
	"log"
	"os"
)

func Delete_directory() {
	// Remove todos os arquivos
	e := os.RemoveAll("./dados")
	if e != nil {
		log.Fatal(e)
	}
	f := os.RemoveAll("./logs")
	if f != nil {
		log.Fatal(f)
	}
	g := os.RemoveAll("./stage")
	if g != nil {
		log.Fatal(g)
	}
	a := os.Remove(".env")
	if a != nil {
		log.Fatal(a)
	}
	b := os.RemoveAll("./certificados")
	if b != nil {
		log.Fatal(b)
	}

}
