package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gin-gonic/gin"
	"github.com/golang/freetype/truetype"
	"github.com/xyproto/palgen"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type yunshi struct {
	Text    string `json:"text"`
	Emotion int    `json:"emotion"`
}

type placement struct {
	Rotate   float64     `json:"rotate"`
	Rotatex  float64     `json:"rotatex"`
	Rotatey  float64     `json:"rotatey"`
	Pos      [][]float64 `json:"pos"`
	FontSize float64     `json:"fontsize"`
}

type ysCfg struct {
	Yunshi         []yunshi    `json:"yunshi"`
	Font           string      `json:"font"`
	Color          []float64   `json:"color"`
	Emotion        []string    `json:"emotion"`
	Placement      []placement `json:"placement"`
	EmotionPalette []color.Palette
	FontFace       *truetype.Font
}

func DrawRP(rp string) image.Image {
	bg, _ := os.ReadFile("./bg.png")
	img, _, _ := image.Decode(bytes.NewReader(bg))
	dc := gg.NewContext(500, 500)
	dc.DrawImage(img, 0, 0)
	dc.SetFontFace(truetype.NewFace(ysConfig.FontFace, &truetype.Options{
		Size: 100,
	}))

	dc.SetRGB(0, 0, 0)

	if len(rp) == 1 {
		dc.DrawString(rp, 347, 300)
	} else if len(rp) == 3 {
		dc.DrawString(rp, 292, 300)
	} else {
		dc.DrawString(rp, 320, 300)
	}
	return dc.Image()
}

func DrawYS(ys string, emo int) image.Image {
	bg, _ := os.ReadFile(ysConfig.Emotion[emo])
	img, _, _ := image.Decode(bytes.NewReader(bg))
	dc := gg.NewContext(img.Bounds().Dx(), img.Bounds().Dy())
	dc.DrawImage(img, 0, 0)
	dc.SetRGB(ysConfig.Color[0], ysConfig.Color[1], ysConfig.Color[2])

	ysLetter := strings.Split(ys, "")

	if len(ysLetter) > len(ysConfig.Placement) {
		log.Fatalln("placement length not enough!")
		return nil
	}

	l := len(ysLetter) - 1

	dc.SetFontFace(truetype.NewFace(ysConfig.FontFace, &truetype.Options{
		Size: ysConfig.Placement[l].FontSize,
	}))
	dc.RotateAbout(gg.Radians(ysConfig.Placement[l].Rotate), ysConfig.Placement[l].Rotatex, ysConfig.Placement[l].Rotatey)
	for i, e := range ysLetter {
		dc.DrawString(e, ysConfig.Placement[l].Pos[i][0], ysConfig.Placement[l].Pos[i][1])
		fmt.Println(e, ysConfig.Placement[l].Pos[i][0], ysConfig.Placement[l].Pos[i][1])
	}
	return dc.Image()
}

var ysConfig ysCfg

func main() {
	jsonFile, _ := os.ReadFile("./ys.json")
	json.Unmarshal(jsonFile, &ysConfig)

	bg, _ := os.ReadFile("bg.png")
	img, _, _ := image.Decode(bytes.NewReader(bg))
	palRP, _ := palgen.Generate(img, 255)
	palRP = append(palRP, image.Transparent)

	for _, e := range ysConfig.Emotion {
		bg, _ = os.ReadFile(e)
		img, _, _ = image.Decode(bytes.NewReader(bg))
		pal, _ := palgen.Generate(img, 255)
		pal = append(pal, image.Transparent)
		ysConfig.EmotionPalette = append(ysConfig.EmotionPalette, pal)
	}

	fontBytes, _ := os.ReadFile(ysConfig.Font)
	ysConfig.FontFace, _ = truetype.Parse(fontBytes)
	fmt.Println(ysConfig)

	rand.Seed(time.Now().UnixNano())

	ginServer := gin.Default()
	ginServer.GET("/jrrp", func(c *gin.Context) {
		rp := c.Query("rp")
		fmt.Println(rp)
		if rp == "" {
			rp = fmt.Sprintf("%d", rand.Intn(101))
		}
		log.Println("rp: " + rp)
		imgCtx := image.NewPaletted(image.Rect(0, 0, 500, 500), palRP)
		buf := new(bytes.Buffer)
		img := DrawRP(rp)
		c.Header("Content-Type", "image/gif")
		draw.Draw(imgCtx, img.Bounds(), img, image.ZP, draw.Src)
		_ = gif.EncodeAll(buf, &gif.GIF{
			Image:     []*image.Paletted{imgCtx},
			Delay:     []int{0},
			LoopCount: 0,
		})
		_, _ = c.Writer.Write(buf.Bytes())
	})

	ginServer.GET("/jrys", func(c *gin.Context) {
		ys := c.Query("ys")
		emo := c.Query("emo")
		if ys == "" {
			ysCol := ysConfig.Yunshi[rand.Intn(len(ysConfig.Yunshi))]
			ys = ysCol.Text
			emo = strconv.Itoa(ysCol.Emotion)
		}
		log.Println("ys: " + ys)
		log.Println("emo: " + emo)
		emoNum, _ := strconv.ParseInt(emo, 10, 64)
		imgCtx := image.NewPaletted(image.Rect(0, 0, 500, 500), ysConfig.EmotionPalette[emoNum])
		buf := new(bytes.Buffer)
		img := DrawYS(ys, int(emoNum))
		c.Header("Content-Type", "image/gif")
		draw.Draw(imgCtx, img.Bounds(), img, image.ZP, draw.Src)
		_ = gif.EncodeAll(buf, &gif.GIF{
			Image:     []*image.Paletted{imgCtx},
			Delay:     []int{0},
			LoopCount: 0,
		})
		_, _ = c.Writer.Write(buf.Bytes())
	})
	_ = ginServer.Run(":11451")
}
