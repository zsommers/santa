package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/zsommers/santa/pkg/santa"
	"github.com/zsommers/santa/pkg/twilio"
)

func main() {
	reallySend := flag.Bool("send", false, "Really send texts")
	nameFile := flag.String("namefile", "names.csv", "Path to names.csv")
	flag.Parse()

	var santas santa.People
	var err error
	for i := 1; i <= 10; i++ {
		log.Printf("Attempting to assign santas round %d", i)
		santas, err = santa.AssignSantas(*nameFile)
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Panic(err)
	}

	var t *twilio.Texter
	if t, err = twilio.NewTexter(*reallySend); err != nil {
		log.Panic(err)
	}

	for _, p := range santas {
		m := fmt.Sprintf(
			"Merry Christmas, %s! You'll be giving a gift to %s 🎅🎄🎁",
			p.Name,
			*p.Recipient,
		)
		if err := t.SendText(p.Mobile, m); err != nil {
			log.Panic(err)
		}
	}
}
