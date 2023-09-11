package main

import (
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/caris-events/gg"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

var (
	boldFont   *sfnt.Font
	mediumFont *sfnt.Font

	titleFace font.Face
	descFace  font.Face
	metaFace  font.Face
	urlFace   font.Face

	baseImage image.Image
)

// InitDraw
func InitDraw() (err error) {
	if boldFont != nil {
		return
	}

	boldFont, err = LoadFont("./assets/fonts/NotoSansCJKtc-Bold.otf")
	if err != nil {
		return err
	}
	mediumFont, err = LoadFont("./assets/fonts/NotoSansCJKtc-Medium.otf")
	if err != nil {
		return err
	}

	titleFace, err = LoadFontFace(boldFont, 10.5)
	if err != nil {
		return err
	}
	descFace, err = LoadFontFace(mediumFont, 7)
	if err != nil {
		return err
	}
	metaFace, err = LoadFontFace(mediumFont, 6.3)
	if err != nil {
		return err
	}
	urlFace, err = LoadFontFace(mediumFont, 6)
	if err != nil {
		return err
	}
	baseImage, err = gg.LoadPNG("./assets/images/cover.png")
	if err != nil {
		return err
	}
	return nil
}

// LoadFont
func LoadFont(v string) (*sfnt.Font, error) {
	file, err := os.ReadFile(v)
	if err != nil {
		return nil, err
	}
	font, err := opentype.Parse(file)
	if err != nil {
		return nil, err
	}
	return font, nil
}

// LoadFontFace
func LoadFontFace(font *sfnt.Font, size float64) (font.Face, error) {
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: size,
		DPI:  300,
	})
	if err != nil {
		return nil, err
	}
	return face, nil
}

// DrawObjectCover
func DrawObjectCover(v *Object) error {
	if err := InitDraw(); err != nil {
		return err
	}
	canvas := gg.NewContextForImage(baseImage)

	// Title
	canvas.SetFontFace(titleFace)
	canvas.SetRGB(0.2, 0.2, 0.2)
	canvas.DrawStringAnchored(v.Name, 90, 80, 0, 0.5)

	// Meta
	canvas.SetFontFace(metaFace)
	canvas.SetRGB(0.5, 0.5, 0.5)

	meta := []string{}
	if v.SecondaryName != "" {
		meta = append(meta, v.SecondaryName)
	}
	meta = append(meta, v.OwnerStr)
	canvas.DrawStringAnchored(strings.Join(meta, "．"), 90, 165, 0, 0)

	// Logo
	logo, err := gg.LoadImage(fmt.Sprintf("./../../docs/%s/%s.jpg", v.Code, v.Code))
	if err != nil {
		return err
	}

	// Scaling Logo
	origW := logo.Bounds().Dx()
	origH := logo.Bounds().Dy()
	aspect := float64(origW) / float64(origH)
	w := 180
	h := int(float64(w) / aspect)

	logoScaled := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.BiLinear.Scale(logoScaled, logoScaled.Rect, logo, logo.Bounds(), draw.Over, nil)
	canvas.DrawImageAnchored(logoScaled, 1000, 125, 0.5, 0.5)

	// URL
	canvas.SetFontFace(urlFace)
	canvas.SetRGB(0.4, 0.4, 0.4)
	canvas.DrawStringAnchored(fmt.Sprintf("https://baka-invade.org/%s/", v.Code), 90, 510, 0, 0)

	// Description
	canvas.SetFontFace(descFace)
	canvas.SetRGB(0.2, 0.2, 0.2)
	canvas.DrawStringWrapped(v.Summary, 90, 220, 0, 0, 700, 1.4, gg.AlignLeft)

	return canvas.SavePNG(fmt.Sprintf("./../../docs/%s/%s_cover.png", v.Code, v.Code))
}

// DrawDictCover
func DrawDictCover(v *Dict) error {
	if err := InitDraw(); err != nil {
		return err
	}
	canvas := gg.NewContextForImage(baseImage)

	// Title
	canvas.SetFontFace(titleFace)
	canvas.SetRGB(0.2, 0.2, 0.2)
	canvas.DrawStringAnchored(v.Word, 90, 80, 0, 0.5)

	// Meta
	canvas.SetFontFace(metaFace)
	canvas.SetRGB(0.5, 0.5, 0.5)
	canvas.DrawStringAnchored(fmt.Sprintf("通常是指：%s", v.ExampleStr), 90, 165, 0, 0)

	// URL
	canvas.SetFontFace(urlFace)
	canvas.SetRGB(0.4, 0.4, 0.4)
	canvas.DrawStringAnchored(fmt.Sprintf("https://baka-invade.org/dict/%s/", v.Code), 90, 510, 0, 0)

	// Description
	canvas.SetFontFace(descFace)
	canvas.SetRGB(0.2, 0.2, 0.2)

	canvas.DrawStringWrapped(fmt.Sprintf("誤用《%s》詞彙可能會傳達不正確的訊息。長期使用經過中國言論審查、中共思想而產生的侵略性詞彙，可能會在無形之中影響自己的思考方式。", v.Word), 90, 220, 0, 0, 700, 1.4, gg.AlignLeft)

	return canvas.SavePNG(fmt.Sprintf("./../../docs/dict/%s/%s_cover.png", v.Code, v.Code))
}
