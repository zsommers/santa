package santa

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Person is a collection of data for one secret santa participant
type Person struct {
	Name      string
	Mobile    string
	Excluded  []string
	Recipient *string
	Santa     *string
}

func (p *Person) String() string {
	r := "?"
	if p.Recipient != nil {
		r = *p.Recipient
	}
	return fmt.Sprintf("%s\tRecipeient: %s\tExcluded: %s", p.Name, r, p.Excluded)
}

func (p *Person) excludes(n string) bool {
	for _, e := range p.Excluded {
		if n == e {
			return true
		}
	}

	return false
}

// People is a collection of secret santa participants
type People []*Person

func (ps People) getByName(n string) (*Person, error) {
	for _, p := range ps {
		if p.Name == n {
			return p, nil
		}
	}

	return nil, fmt.Errorf("%s not found", n)
}

func (ps People) needSanta() People {
	np := make(People, 0)
	for _, p := range ps {
		if p.Santa == nil {
			np = append(np, p)
		}
	}

	return np
}

func (ps People) needRecipient() People {
	np := make(People, 0)
	for _, p := range ps {
		if p.Recipient == nil {
			np = append(np, p)
		}
	}

	return np
}

func loadPeople(nameFile string) (People, error) {
	file, err := os.Open(nameFile)
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

	ps := make(People, 0, len(records))
	for _, r := range records {
		ps = append(ps, &Person{
			Name:     r[0],
			Mobile:   r[1],
			Excluded: strings.Split(r[2], " "),
		})
	}

	return ps, nil
}

func findRecipient(giver *Person, needSanta People) (*Person, error) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(needSanta), func(i, j int) {
		needSanta[i], needSanta[j] = needSanta[j], needSanta[i]
	})

	for _, ns := range needSanta {
		if giver.Name == ns.Name {
			continue
		}
		if giver.excludes(ns.Name) {
			continue
		}
		return ns, nil
	}

	return nil, fmt.Errorf("no suitable Recipient found")
}

// AssignSantas will perform one round of assigning secret santas
func AssignSantas(nameFile string) (People, error) {
	ps, err := loadPeople(nameFile)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		r, err := findRecipient(p, ps.needSanta())
		if err != nil {
			return nil, err
		}

		p.Recipient = &r.Name
		r.Santa = &p.Name
	}

	return ps, nil
}
