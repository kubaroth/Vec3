// An example which sets up a GUI windows with using xgb (X11).
// A new window listen for keyboard input.
// ] - starts slow render
// [ - interrupts render and saves the image
// Esc - exits
// w - move camera forward
// s - move camera backward


// Get dependencies for this example:
// go get -v -u ...
// go build ./main_xgb.go
// ./main_xgb

package main

import (
	. "github.com/kubaroth/Vec3"
	"errors"
	"fmt"
_	"image/png"
	"os"
_	"sync"

	
	"time"
	"math/rand"
	"log"
	
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
_	"github.com/BurntSushi/xgbutil/xwindow"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xgraphics" // painting
_	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgbutil/xwindow"

)

var counter int32
var world HittableList

// Just iniaitlize the RNG seed for generating random background colors.
func init() {
	rand.Seed(time.Now().UnixNano())
	world = HittableList{}
}

type RenderSetup struct {
	cam Camera
	world HittableList
	done chan int
	X *xgbutil.XUtil
	win *xwindow.Window
}

// func renderSetup(cam Camera, world HittableList, done chan int, X *xgbutil.XUtil, win *xwindow.Window){
func renderSetup(parms RenderSetup){

	// Enable this to see BVH culling in action. 5sec vs 28sec for []Hittablelist
	// for i:=0; i<500; i++ {
	// 	world.Add(Sphere{NewVec3(0,float32(i)/10., float32(i + 1)), 0.5})
	// }
	
	bvh := NewBVHSplit(parms.world.Objects,0,len(parms.world.Objects))
	_ = bvh


	samples := 16 // increase samples to see the problem of updating image during interrupted render
	
	start := time.Now(); _ = start

	path := os.Getenv("HOME") + "/storage/downloads/img.png" // termux preview
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		path = "img.png"
	}
	
	img := Render(parms.cam, samples, &parms.world, nil, parms.done)  // pass bvh instead of nil to use BVH_node container

	// if len(img.Pix) == 0 { // interupted - don't paint
	// 	return
	// }
	
	// Write image to pixmap and update content of the window
	ximg := xgraphics.NewConvert(parms.X, img)
	ximg.XSurfaceSet(parms.win.Id)
	ximg.XDraw()
	ximg.XPaint(parms.win.Id)
	parms.win.Resize(parms.cam.Width, parms.cam.Height)
	
	// saving png
	
	// fmt.Println("saving into:", path)

	// f, err := os.Create(path)
	// if err != nil {
	// 	panic(err)
	// }
	
	// defer f.Close()
	// if err = png.Encode(f, img); err != nil {
	// 	fmt.Printf("failed to encode: %v", err)
	// }

	// fmt.Println("time", time.Since(start))
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
	// 		prev_rx = rx
	// 		prev_ry = ry

	// 	},
	// 	func(X *xgbutil.XUtil, rx, ry, ex, ey int) {
	// 		log.Println("release", rx, ry)
	// 		bounds.Max.X = rx
	// 		bounds.Max.Y = ry
	// 	})

	done := make(chan int)
	cam := NewCamera(NewVec3(0,0,0), NewVec3(0,0,-1), 400)
	// TODO: Include into camera samples
	//       Rethink how to rework camera struct to make update its state simpler


	render_parms := RenderSetup{cam, world, done, X, win}
	_ = render_parms
	

	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			log.Println("quiting...")
			// graceful exit
			xevent.Detach(win.X, win.Id)
			mousebind.Detach(win.X, win.Id)
			keybind.Detach(win.X, win.Id)
			win.Destroy()
			xevent.Quit(X)

			// close the done chanell on exit
			close(done)

		}).Connect(X, win.Id, "Escape", true)

	first_run := func(){
		render_parms.world.Add(Sphere{NewVec3(0,0,-1), 0.5})
		render_parms.world.Add(Sphere{NewVec3(0,-100.5,-1), 100.0})
		renderSetup(render_parms)
	}
	first_run()
	
	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println("Rendering...")
			first_run()

		}).Connect(X, win.Id, "bracketright", true)

	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println("cancelling...")
			go func(){
				done <- 1
			}()
		}).Connect(X, win.Id, "bracketleft", true)

	// A test with different number of objects in the scene
	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println("key P was pressed...")
			go func(){
				render_parms.world.Objects = nil // clear the slice
				render_parms.world.Add(Sphere{NewVec3(0,0,-1), 0.5}) // add only a single sphere
				renderSetup(render_parms)
			}()
		}).Connect(X, win.Id, "p", true)

	
	// Move forward
	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			fmt.Println("key W was pressed...")
		}).Connect(X, win.Id, "w", true)

	// We should consider this method where we refrain from re-rendering on multiple
	// key presses only trigger the render once the key is release
	// NOTE: The 'release' is also triggered on press!!!!
	keybind.KeyReleaseFun( 
		func(X *xgbutil.XUtil, e xevent.KeyReleaseEvent) {
			fmt.Println("key W was released...")

			go func(){  // stop previous run
				done <- 1 
			}()
			go func(){
				render_parms.cam = NewCamera( render_parms.cam.Origin.Add(NewVec3(0,0,-0.01)), NewVec3(0,0,-1), render_parms.cam.Width)
				renderSetup(render_parms)
			}()
			
		}).Connect(X, win.Id, "w", true)

	
	// Move bacwkward
	keybind.KeyPressFun(
		func(X *xgbutil.XUtil, e xevent.KeyPressEvent) {
			go func(){ // stop previous run
				done <- 1 
			}()
			go func(){
				render_parms.cam = NewCamera( render_parms.cam.Origin.Add(NewVec3(0,0,0.01)), NewVec3(0,0,-1), render_parms.cam.Width)
				renderSetup(render_parms)
			}()
		}).Connect(X, win.Id, "s", true)

	
	if err != nil {
		log.Fatal(err)
	}

	win.Listen(xproto.EventMaskPropertyChange, xproto.EventMaskStructureNotify)

	// Main Event loop
	xevent.Main(X)

}
