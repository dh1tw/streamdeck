package streamdeck

type Config struct {
	ProductID        uint16 // ProductID is the USB ProductID
	NumButtonColumns int
	NumButtonRows    int
	Spacer           int // Spacer is the spacing distance (in pixel) of two buttons on the Stream Deck.
	ButtonSize       int
	ImageFormat      string
	ImageRotate      bool
}

func (c Config) NumButtons() int {
	return c.NumButtonRows * c.NumButtonColumns
}

// PanelWidth is the total screen width of the Stream Deck (including spacers).
func (c *Config) PanelWidth() int {
	return c.NumButtonColumns*c.ButtonSize + c.Spacer*(c.NumButtonColumns-1)
}

// PanelHeight is the total screen height of the stream deck (including spacers).
func (c *Config) PanelHeight() int {
	return c.NumButtonRows*c.ButtonSize + c.Spacer*(c.NumButtonRows-1)
}

var Original = Config{
	ProductID:        0x60,
	NumButtonColumns: 5,
	NumButtonRows:    3,
	Spacer:           19,
	ButtonSize:       72,
	ImageFormat:      "bmp",
}

var Original2 = Config{
	ProductID:        0x80,
	NumButtonColumns: 5,
	NumButtonRows:    3,
	Spacer:           19,
	ButtonSize:       72,
	ImageFormat:      "jpg",
	ImageRotate:      true,
}

var Plus = Config{
	ProductID:        0x0084,
	NumButtonColumns: 4,
	NumButtonRows:    2,
	Spacer:           19,
	ButtonSize:       120,
	ImageFormat:      "jpg",
}
