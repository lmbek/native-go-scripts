package main

import (
	"fmt"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image"
	"image/color"
	"log"
)

var (
	blue0    = color.RGBA{0x00, 0x00, 0x1f, 0xff}
	blue1    = color.RGBA{0x00, 0x00, 0x3f, 0xff}
	darkGray = color.RGBA{0x3f, 0x3f, 0x3f, 0xff}
	green    = color.RGBA{0x00, 0x7f, 0x00, 0x7f}
	red      = color.RGBA{0x7f, 0x00, 0x00, 0x7f}
	yellow   = color.RGBA{0x3f, 0x3f, 0x00, 0x3f}
)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title: "Basic Shiny Example",
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		size0 := image.Point{256, 256}
		b, err := s.NewBuffer(size0)
		if err != nil {
			log.Fatal(err)
		}
		defer b.Release()

		t0, err := s.NewTexture(size0)
		if err != nil {
			log.Fatal(err)
		}
		defer t0.Release()
		t0.Upload(image.Point{}, b, b.Bounds())

		size1 := image.Point{32, 20}
		t1, err := s.NewTexture(size1)
		if err != nil {
			log.Fatal(err)
		}
		defer t1.Release()
		t1.Fill(t1.Bounds(), green, screen.Src)
		t1.Fill(t1.Bounds().Inset(2), red, screen.Over)
		t1.Fill(t1.Bounds().Inset(4), red, screen.Src)

		//var sz size.Event
		for {
			e := w.NextEvent()

			// This print message is to help programmers learn what events this
			// example program generates. A real program shouldn't print such
			// messages; they're not important to end users.
			format := "got %#v\n"
			if _, ok := e.(fmt.Stringer); ok {
				format = "got %v\n"
			}
			fmt.Printf(format, e)

			switch e := e.(type) {
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				if e.Code == key.CodeEscape {
					return
				}

			case paint.Event:

				w.Publish()

			case size.Event:
				//sz = e

			case error:
				log.Print(e)
			}
		}
	})
}
