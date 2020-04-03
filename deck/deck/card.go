package deck

type CardValue int
type CardType int

const (
	Joker CardValue = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Knight
	Queen
	King
)

const (
	NoType CardType = iota
	Spades
	Diamonds
	Clubs
	Hearts
)

type Card struct {
	value     CardValue
	cType     CardType
	isVisible bool
}

func (v *CardValue) String() string {
	values := []string{
		"Joker",
		"Ace",
		"Two",
		"Three",
		"Four",
		"Five",
		"Six",
		"Seven",
		"Eight",
		"Nine",
		"Ten",
		"Knight",
		"Queen",
		"King",
	}

	return values[*v]
}

func (t *CardType) String() string {
	types := []string{
		"",
		"Spades",
		"Diamonds",
		"Clubs",
		"Hearts",
	}

	return types[*t]
}

func (c *Card) SetVisible(isVisible bool) {
	c.isVisible = isVisible
}

func (c *Card) IsVisible() bool {
	return c.isVisible
}

func (c *Card) GetValue() CardValue {
	return c.value
}

func getAllCardValues() []CardValue {
	values := []CardValue{
		Ace,
		Two,
		Three,
		Four,
		Five,
		Six,
		Seven,
		Eight,
		Nine,
		Ten,
		Knight,
		Queen,
		King,
	}
	return values
}

