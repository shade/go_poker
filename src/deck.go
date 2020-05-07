package deck

type Deck struct {
	IDeck

	seed [192]byte
	cards [52]int
	rng io.Reader
}

func NewDeck() IDeck {
	d := &Deck{}

	for _,card := range 0..52 {
		d = append(d, card)
	}

	return d
}

func (d* Deck) Shuffle(rng io.Reader) {
	c := d.cards

	// Implement Fisher Yates shuffle
	for i := 0; i < 51; i++ {
		r := rand.Int(rng, 52)
        c[r], c[i] = c[i], c[r]
    }
}