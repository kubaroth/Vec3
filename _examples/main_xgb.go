// Get dependencies for this example:
// go get -v -u ...
// go build ./main_xgb.go
// ./main_xgb


// Example create-window shows how to create a window, map it, resize it,
// and listen to structure and key events (i.e., when the window is resized
// by the window manager, or when key presses/releases are made when the
// window has focus). The events are printed to stdout.


package main

import (
	. "github.com/kubaroth/Vec3"
	"errors"
	"fmt"
	"image/png"
	"os"
	"sync"

	
	"time"
	"math/rand"
	"log"
	
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
_	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
_	"github.com/BurntSushi/xgbutil/xgraphics" // painting
_	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgbutil/xwindow"

)

var counter int32

// Just iniaitlize the RNG seed for generating random background colors.
func init() {
	rand.Seed(time.Now().UnixNano())
}


func renderSetup(){

	world := HittableList{}
	world.Add(Sphere{NewVec3(0,0,-1), 0.5})
	world.Add(Sphere{NewVec3(0,-100.5,-1), 100.0})
	_ = world

	// Enable this to see BVH culling in action. 5sec vs 28sec for []Hittablelist
	// for i:=0; i<500; i++ {
	// 	world.Add(Sphere{NewVec3(0,float32(i)/10., float32(i + 1)), 0.5})
	// }
	
	bvh := NewBVHSplit(world.Objects,0,len(world.Objects))
	_ = bvh


	cam := NewCamera(NewVec3(0,0,0), NewVec3(0,0,-1), 200)
	samples := 1
	
	start := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)
	
	go func() {

		defer wg.Done()
		img := Render(cam, samples, &world, nil)  // pass bvh instead of nil to use BVH_node container

		// saving png

		path := os.Getenv("HOME") + "/storage/downloads/img.png" // termux preview
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			path = "img.png"
		}
	
		fmt.Println("saving into:", path)

		f, err := os.Create(path)
		if err != nil {
			panic(err)
		}
	
		defer f.Close()
		if err = png.Encode(f, img); err != nil {
			fmt.Printf("failed to encode: %v", err)
		}

	}()

	fmt.Println("\nWaiting...")
	wg.Wait()

	fmt.Println("time", time.Since(start))
}





// newWindow creates a new window with a random background color. It sets the
// WM_PROTOCOLS property to contain the WM_DELETE_WINDOW atom. It also sets
// up a ClientMessage event handler so that we know when to destroy the window.
// We also set up a mouse binding so that clicking inside a window will
// create another one.
func newWindow(X *xgbutil.XUtil) *xwindow.Window {
	counter++

	win, err := xwindow.Generate(X)
	if err != nil {
		log.Fatal(err)
	}

	// Get a random background color, create the window (ask to receive button
	// release events while we're at it) and map the window.
	bgColor := rand.Intn(0xffffff + 1)
	win.Create(X.RootWin(), 0, 0, 200, 200,
		xproto.CwBackPixel|xproto.CwEventMask,
		uint32(bgColor), xproto.EventMaskButtonRelease)

	// WMGracefulClose does all of the work for us. It sets the appropriate
	// values for WM_PROTOCOLS, and listens for ClientMessages that implement
	// the WM_DELETE_WINDOW protocol. When one is found, the provided callback
	// is executed.
	win.WMGracefulClose(
		func(w *xwindow.Window) {
			// Detach all event handlers.
			// This should always be done when a window can no longer
			// receive events.
			xevent.Detach(w.X, w.Id)
			mousebind.Detach(w.X, w.Id)
			w.Destroy()

			// Exit if there are no more windows left.
			counter--
			if counter == 0 {
				os.Exit(0)
			}
		})

	// It's important that the map comes after setting WMGracefulClose, since
	// the WM isn't obliged to watch updates to the WM_PROTOCOLS property.
	win.Map()

	// A mouse binding so that a left click will spawn a new window.
	// Note that we don't issue a grab here. Typically, window managers will
	// grab a button press on the client window (which usually activates the
	// window), so that we'd end up competing with the window manager if we
	// tried to grab it.
	// Instead, we set a ButtonRelease mask when creating the window and attach
	// a mouse binding *without* a grab.
	err = mousebind.ButtonReleaseFun(
		func(X *xgbutil.XUtil, ev xevent.ButtonReleaseEvent) {
			newWindow(X)
		}).Connect(X, win.Id, "1", false, false)
	if err != nil {
		log.Fatal(err)
	}

	return win
}

func main(){

	// XCB - version - determin bounds with User's input
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}
	mousebind.Initialize(X)
	keybind.Initialize(X)

	
	win := newWindow(X)
	

	err = mousebind.ButtonPressFun(
		func(X *xgbutil.XUtil, e xevent.ButtonPressEvent) {
			log.Println("Painting")
		}).Connect(X, win.Id, "1", false, true)

	// var prev_rx, prev_ry int
	// mousebind.Drag(X, win.Id, win.Id, "1", false,
	// 	func(X *xgbutil.XUtil, rx, ry, ex, ey int) (bool, xproto.Cursor) {
	// 		log.Println("starting", rx, ry)
	// 		bounds.Min.X = rx
	// 		bounds.Min.Y = ry
	// 		prev_rx = rx
	// 		prev_ry = ry
	// 		return true, 0
	// 	},
	// 	func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
	// 		// log.Println("painting", rx, ry)
	// 		paint(canvas, win, rx, ry, prev_rx, prev_ry)
	// 		prev_rx = rx
	// 		prev_ry = ry

	// 	},
	// 	func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
	// 		log.Println("release", rx, ry)
	// 		bounds.Max.X = rx
	// 		bounds.Max.Y = ry
	// 		// push on the undo stack
	// 		undo_step := make([]byte, len(canvas.Pix))
	// 		copy(undo_step, canvas.Pix)
	// 		SETTINGS.undos = append(SETTINGS.undos, undo_step)

	// 	})

	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			log.Println("quiting...")
			// graceful exit
			xevent.Detach(win.X, win.Id)
			mousebind.Detach(win.X, win.Id)
			keybind.Detach(win.X, win.Id)
			win.Destroy()
			xevent.Quit(X)

		}).Connect(X, win.Id, "Escape", true)

	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println("Rendering...")
			renderSetup()
		}).Connect(X, win.Id, "bracketright", true)

	if err != nil {
		log.Fatal(err)
	}

	win.Listen(xproto.EventMaskPropertyChange, xproto.EventMaskStructureNotify)

	// Main Event loop
	xevent.Main(X)

}
