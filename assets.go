package streamdeck

import (
	"embed"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

//go:embed assets
var assetDirectory embed.FS

var MonoRegular *truetype.Font
var MonoMedium *truetype.Font

func init() {
	err := initHelp()
	if err != nil {
		panic(err)
	}
}

func initHelp() error {
	f, err := loadFromAssets("assets/mplus-1m-regular.ttf")
	if err != nil {
		return err
	}
	MonoRegular = f

	f, err = loadFromAssets("assets/mplus-1m-medium.ttf")
	if err != nil {
		return err
	}
	MonoMedium = f

	return nil
}

func loadFromAssets(fn string) (*truetype.Font, error) {
	f, err := assetDirectory.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return freetype.ParseFont(data)
}
