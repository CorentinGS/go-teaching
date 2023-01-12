package goroutines_simple_vs_complex

import "time"

// Pineapple is a struct that represents a database object with sensitive data that should be hidden
type Pineapple struct {
	Paro       string `faker:"name"`
	Turkey     string `faker:"name"`
	Banana     string `faker:"name"`
	Age        int    `faker:"number"`
	Size       int    `faker:"number"`
	IsAlive    bool
	ID         uint
	SecretCode []byte
	Created    time.Time
	Updated    time.Time
}

func (p *Pineapple) ToSafePineApple() SafePineApple {
	return SafePineApple{
		Paro:    p.Paro,
		Turkey:  p.Turkey,
		Banana:  p.Banana,
		IsAlive: p.IsAlive,
		Age:     p.Age,
		ID:      p.ID,
	}
}

// SafePineApple is a struct that represents a Pineapple object without sensitive data
type SafePineApple struct {
	Paro    string
	Turkey  string
	Banana  string
	IsAlive bool
	Age     int
	ID      uint
}
