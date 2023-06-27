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
	"log"
)

func main() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Hello World",
			Width:  1000,
			Height: 600,
		})
		if err != nil {
			log.Fatal(err)
		}
		defer w.Release()

		var t screen.Texture

		buffer, err := s.NewBuffer(image.Point{X: 1000, Y: 600})
		if err != nil {
			log.Fatal(err)
		}
		defer buffer.Release()
		size2 := buffer.Bounds().Size()
		//width, height := size.X, size.Y
		fmt.Println(size2)

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
				// Re-draw the texture with white text on black background
				if t != nil {
					t.Release()
				}

				fmt.Println("here is: --> ", w, " <---- ")
				/*
								size := w.Bounds().Size() // TODO: commented

								b, err := s.NewBuffer(size)
								if err != nil {
									log.Fatal(err)
								}
								defer b.Release()

							d := &font.Drawer{
								Src:  image.White,
								Face: basicfont.Face7x13,
								Dst:  b.RGBA(),
								Dot:  fixed.Point26_6{X: fixed.I(10), Y: fixed.I(100)},
							}
							d.DrawString("Hello World")

						t, err = s.NewTexture(size)
						if err != nil {
							log.Fatal(err)
						}
						defer t.Release()


					t.Upload(image.Point{}, b, b.Bounds())
				*/
				w.Publish()

			case size.Event:
				// Re-size the texture and re-draw it
				w.Send(paint.Event{})
			}
		}
	})
}
