package main

import (
	"fmt"
	"log"
	"strings"
	"syscall"

	"github.com/gonutz/w32"
	"github.com/hnakamur/w32syscall"
)

func main() {

	//AppCalculatorWindowMove()
	//RunningApp()

	// process, isOk := w32.EnumAllProcesses()
	// fmt.Println(isOk)
	// for i, p := range process {

	// 	fmt.Println(i, p)
	// }

	hwnd := w32.GetActiveWindow()
	text := w32.GetWindowText(hwnd)
	fmt.Println("Active Window::", text, hwnd)

	w32.EnumWindows(func(w w32.HWND) bool {

		isVis := w32.IsWindowVisible(w)
		if isVis == true {

			fmt.Println(w, w32.GetWindowText(w))

			//131926
			wval := fmt.Sprintf("%v", w)
			if wval == "131926" {
				w32.EnumChildWindows(w, func(wh w32.HWND) bool {

					fmt.Println(">>", w32.GetWindowText(wh))
					return true
				})
			}

		}

		return true

	})

	// uintV, err := syscall.GetVersion()
	// computer := w32.GetComputerName()
	// fmt.Println(uintV, err, computer)

	//h1 := w32.HWND(hD)
	//EnumWindows func(callback func(window HWND)) bool
	//EnumWindows func(callback func(hwnd Handle, lparam uintptr) bool, lparam uintptr) (err error)
	//EnumWindows func(callback TYPE) bool

	//func Name() re{

	//}

	// w32.EnumWindows(func(hD w32.HWND) bool {

	// })

	/*
		hD := w32.GetDesktopWindow()
		// fmt.Println(hD)
		clss, isOk := w32.GetClassName(hD)
		fmt.Println(clss, isOk)

		rect := w32.GetWindowRect(hD)
		wText := w32.GetWindowText(hD)
		fmt.Println(wText, rect.Width(), rect.Height())
		hndle, dw := w32.GetWindowThreadProcessId(hD)
		fmt.Println(hndle, dw)
		w32.MessageBox(hD, "Hello ...", "Bismillah", 0)
	*/

	// w32.EnumWindows(func(hD w32.HWND) {

	// 	fmt.Println(hD)

	// })

	//w32.EnumWindows(func(window w32.HWND) bool))
	//func(callback func(window w32.HWND)) bool

	// err = w32syscall.EnumWindows(func(hwnd syscall.Handle, lparam uintptr) bool {

	// 	h := w32.HWND(hwnd)

	// 	text := w32.GetWindowText(h)
	// 	if strings.Contains(text, "Calculator") {
	// 		w32.MoveWindow(h, 0, 0, 200, 600, true)
	// 	}
	// 	return true
	// }, 0)

	// if err != nil {
	// 	log.Fatalln(err)
	// }

}

//RunningApp ..
func RunningApp() {

	err := w32syscall.EnumWindows(func(hwnd syscall.Handle, lparam uintptr) bool {

		h := w32.HWND(hwnd)
		text := w32.GetWindowText(h)
		isWindow := w32.IsWindow(h)
		isEnable := w32.IsWindowEnabled(h)
		isVisible := w32.IsWindowVisible(h)
		rect := w32.GetWindowRect(h)
		dword := w32.GetCurrentProcessId()
		clssName, _ := w32.GetClassName(h)
		//w32.EnumProcesses()
		//w32.GetActiveWindow()

		//clssName==CabinetWClass==FOLDER
		//Windows.UI.Core.CoreWindow
		//https://stackoverflow.com/questions/10246444/how-can-i-get-enumwindows-to-list-all-windows

		if isWindow && isEnable && isVisible == true {
			fmt.Println(isWindow, isEnable, isVisible, rect.Width(), dword, clssName, text)
		}

		// if strings.Contains(text, "Calculator") {
		// 	w32.MoveWindow(h, 0, 0, 200, 600, true)
		// }

		return true

	}, 0)

	if err != nil {
		log.Fatalln(err)
	}

}

//AppCalculatorWindowMove open the calculator program then run this func and see the magic
func AppCalculatorWindowMove() {

	//https://stackoverflow.com/questions/29447807/taking-control-of-another-window-with-go
	err := w32syscall.EnumWindows(func(hwnd syscall.Handle, lparam uintptr) bool {
		h := w32.HWND(hwnd)
		text := w32.GetWindowText(h)
		if strings.Contains(text, "Calculator") {
			w32.MoveWindow(h, 0, 0, 200, 600, true)
		}
		return true
	}, 0)
	if err != nil {
		log.Fatalln(err)
	}
}
