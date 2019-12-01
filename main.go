package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type person struct {
	name      string
	mobile    string
	excluded  []string
	recipient *string
	santa     *string
}

func (p *person) String() string {
	r := "?"
	if p.recipient != nil {
		r = *p.recipient
	}
	return fmt.Sprintf("%s\tRecipeient: %s\tExcluded: %s", p.name, r, p.excluded)
}

func (p *person) excludes(n string) bool {
	for _, e := range p.excluded {
		if n == e {
			return true
		}
	}

	return false
}

type people []*person

func (ps people) getByName(n string) (*person, error) {
	for _, p := range ps {
		if p.name == n {
			return p, nil
		}
	}

	return nil, fmt.Errorf("%s not found", n)
}

func (ps people) needSanta() people {
	np := make(people, 0)
	for _, p := range ps {
		if p.santa == nil {
			np = append(np, p)
		}
	}

	return np
}

func (ps people) needRecipient() people {
	np := make(people, 0)
	for _, p := range ps {
		if p.recipient == nil {
			np = append(np, p)
		}
	}

	return np
}

func loadPeople() (people, error) {
	file, err := os.Open("names.csv")
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)

	// Burn header
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	ps := make(people, 0, len(records))
	for _, r := range records {
		ps = append(ps, &person{
			name:     r[0],
			mobile:   r[1],
			excluded: strings.Split(r[2], " "),
		})
	}

	return ps, nil
}

func findRecipient(giver *person, needSanta people) (*person, error) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(needSanta), func(i, j int) {
		needSanta[i], needSanta[j] = needSanta[j], needSanta[i]
	})

	for _, ns := range needSanta {
		if giver.name == ns.name {
			continue
		}
		if giver.excludes(ns.name) {
			continue
		}
		return ns, nil
	}

	return nil, fmt.Errorf("no suitable recipient found")
}

func assignSantas() (people, error) {
	ps, err := loadPeople()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		r, err := findRecipient(p, ps.needSanta())
		if err != nil {
			return nil, err
		}

		p.recipient = &r.name
		r.santa = &p.name
	}

	return ps, nil
}

func main() {
	var santas people
	var err error
	for i := 1; i <= 10; i++ {
		log.Printf("Attempting to assign santas round %d", i)
		santas, err = assignSantas()
		if err == nil {
			break
		}
	}

	if err != nil {
		log.Panic(err)
	}

	var t *Texter
	if t, err = NewTexter(); err != nil {
		log.Panic(err)
	}

	for _, p := range santas {
		m := fmt.Sprintf("Merry Christmas, %s! You'll be giving a gift to %s 🎅🎄🎁", p.name, *p.recipient)
		if err := t.SendText(p.mobile, m); err != nil {
			log.Panic(err)
		}
	}
}
