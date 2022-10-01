package main

import (
	"fmt"
	"os"

	"github.com/go-ole/go-ole"
	wa "github.com/hnakamur/w32uiautomation"
)

func main() {

	//https://docs.microsoft.com/en-us/windows/win32/winauto/uiauto-creatingcuiautomation
	//https://en.it1352.com/article/3b1b3f1c45034514b9747512a77ace6d.html

	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	auto, err := wa.NewUIAutomation()
	if err != nil {
		//return err
		fmt.Println("ERR1::", err.Error())
		os.Exit(1)
	}

	root, err := auto.GetRootElement()
	if err != nil {
		//return err
		fmt.Println("ERR2::", err.Error())
	}
	defer root.Release()
	className, err := root.Get_CurrentCurrentClassName()
	fmt.Println(className, err)

	name, err := root.Get_CurrentName()
	fmt.Println(name, err)

	// condVal := wa.NewVariantString("Pane")
	// condition, err := auto.CreatePropertyCondition(wa.UIA_NamePropertyId, condVal)
	// fmt.Println(condition, err)

	// if err != nil {
	// 	fmt.Println("ERR3::", err.Error())
	// 	os.Exit(1)
	// }

	// elm1, err := root.FindFirst(wa.TreeScope_Children, condition)

	// eName, err := elm1.Get_CurrentName()
	// fmt.Println(eName, err)

}
