package tools

import (
	"log"
)

func HandleErr(mes string, err error) {
	if err != nil {
		log.Fatal(mes, err)
	}
}
