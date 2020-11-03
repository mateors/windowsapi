package main

import (
	"encoding/hex"
	"fmt"
	"strings"
	"syscall"
	"unsafe"
)

// Windows API functions
var (
	modKernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32      = syscall.NewLazyDLL("user32.dll")
	uiauto      = syscall.NewLazyDLL("UIAutomationCore.dll")

	getCurrentProcessId          = modKernel32.NewProc("GetCurrentProcessId")
	openProcess                  = modKernel32.NewProc("OpenProcess") //???????????
	procCloseHandle              = modKernel32.NewProc("CloseHandle")
	procCreateToolhelp32Snapshot = modKernel32.NewProc("CreateToolhelp32Snapshot")
	procProcess32First           = modKernel32.NewProc("Process32FirstW")
	procProcess32Next            = modKernel32.NewProc("Process32NextW")

	enumWindows              = user32.NewProc("EnumWindows")
	getWindowTextLength      = user32.NewProc("GetWindowTextLengthW")
	getWindowText            = user32.NewProc("GetWindowTextW")
	isWindow                 = user32.NewProc("IsWindow")
	isDialogMessage          = user32.NewProc("IsDialogMessageW") // IsDialogMessage function is intended for modeless dialog boxes, you can use it with any window that contains controls, enabling the windows to provide the same keyboard selection as is used in a dialog box.
	isWindowEnabled          = user32.NewProc("IsWindowEnabled")
	isWindowVisible          = user32.NewProc("IsWindowVisible")
	getWindowRect            = user32.NewProc("GetWindowRect")
	getClassName             = user32.NewProc("GetClassNameW")
	getWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	isIconic                 = user32.NewProc("IsIconic")
	getWindowInfo            = user32.NewProc("GetWindowInfo")
	enumPropsW               = user32.NewProc("EnumPropsW")
	getWindowLongPtrW        = user32.NewProc("GetWindowLongPtrW") //???
)

// Some constants from the Windows API
const (
	ERROR_NO_MORE_FILES = 0x12
	MAX_PATH            = 260
)

// PROCESSENTRY32 is the Windows API structure that contains a process's
// information. [https://docs.microsoft.com/en-us/windows/win32/api/tlhelp32/ns-tlhelp32-processentry32]
type PROCESSENTRY32 struct {
	Size              uint32
	CntUsage          uint32
	ProcessID         uint32
	DefaultHeapID     uintptr
	ModuleID          uint32
	CntThreads        uint32
	ParentProcessID   uint32
	PriorityClassBase int32
	Flags             uint32
	ExeFile           [MAX_PATH]uint16
}

// type (
// 	BOOL          uint32
// 	BOOLEAN       byte
// 	BYTE          byte
// 	DWORD         uint32
// 	DWORD64       uint64
// 	HANDLE        uintptr
// 	HLOCAL        uintptr
// 	LARGE_INTEGER int64
// 	LONG          int32
// 	LPVOID        uintptr
// 	SIZE_T        uintptr
// 	UINT          uint32
// 	ULONG_PTR     uintptr
// 	ULONGLONG     uint64
// 	WORD          uint16
//	ATOM          uint16
//    )
//LPCTSTR

//https://docs.microsoft.com/en-us/windows/win32/winprog/windows-data-types
type WINDOWINFO struct {
	cbSize          uint32
	rcWindow        RECT
	rcClient        RECT
	dwStyle         uint32
	dwExStyle       uint32
	dwWindowStatus  uint32
	cxWindowBorders uint32
	cyWindowBorders uint32
	atomWindowType  uint16
	wCreatorVersion uint16
}

// WindowsProcess is an implementation of Process for Windows.
type WindowsProcess struct {
	pid  int
	ppid int
	exe  string
}

//Pid ...
func (p *WindowsProcess) Pid() int {
	return p.pid
}

//PPid ...
func (p *WindowsProcess) PPid() int {
	return p.ppid
}

//Executable ...
func (p *WindowsProcess) Executable() string {
	return p.exe
}

func newWindowsProcess(e *PROCESSENTRY32) *WindowsProcess {
	// Find when the string ends for decoding
	end := 0
	for {
		if e.ExeFile[end] == 0 {
			break
		}
		end++
	}

	return &WindowsProcess{
		pid:  int(e.ProcessID),
		ppid: int(e.ParentProcessID),
		exe:  syscall.UTF16ToString(e.ExeFile[:end]),
	}
}

func findProcess(pid int) (Process, error) {
	ps, err := processes()
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.Pid() == pid {
			return p, nil
		}
	}

	return nil, nil
}

func processes() ([]Process, error) {

	handle, _, _ := procCreateToolhelp32Snapshot.Call(
		0x00000002,
		0)
	if handle < 0 {
		return nil, syscall.GetLastError()
	}
	defer procCloseHandle.Call(handle)

	var entry PROCESSENTRY32
	entry.Size = uint32(unsafe.Sizeof(entry))
	ret, _, _ := procProcess32First.Call(handle, uintptr(unsafe.Pointer(&entry)))
	if ret == 0 {
		return nil, fmt.Errorf("error retrieving process info!")
	}

	results := make([]Process, 0, 50)
	for {
		results = append(results, newWindowsProcess(&entry))

		ret, _, _ := procProcess32Next.Call(handle, uintptr(unsafe.Pointer(&entry)))
		if ret == 0 {
			break
		}
	}

	return results, nil
}

//

type Process interface {
	// Pid is the process ID for this process.
	Pid() int

	// PPid is the parent process ID for this process.
	PPid() int

	// Executable name running this process. This is not a path to the
	// executable.
	Executable() string
}

func Processes() ([]Process, error) {
	return processes()
}

func FindProcess(pid int) (Process, error) {
	return findProcess(pid)
}

type (
	HANDLE uintptr
	HWND   HANDLE
	DWORD  uint32
)

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

type POINT struct {
	X, Y int32
}

type RECT struct {
	Left, Top, Right, Bottom int32
}

func (r *RECT) Width() int32 {
	return r.Right - r.Left
}

func (r *RECT) Height() int32 {
	return r.Bottom - r.Top
}

// BOOL EnumWindows(
// 	WNDENUMPROC lpEnumFunc,
// 	LPARAM      lParam
//   );

//WNDENUMPROC::
// BOOL CALLBACK EnumWindowsProc(
// 	_In_ HWND   hwnd,
// 	_In_ LPARAM lParam
//   );

//https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-enumwindows
func EnumWindows(callback func(window HWND) bool) bool {

	f := syscall.NewCallback(func(w, _ uintptr) uintptr {
		if callback(HWND(w)) {
			return 1
		}
		return 0
	})
	ret, _, _ := enumWindows.Call(f, 0)
	return ret != 0
}

// int EnumPropsA(
// 	HWND          hWnd,
// 	PROPENUMPROCA lpEnumFunc
//   );

// BOOL Propenumproca(
// 	HWND Arg1,
// 	LPCSTR Arg2,
// 	HANDLE Arg3
//   )
func EnumPropsW(parent HWND, callback func(window HWND, lpcstr *string, handle uintptr) bool) int {

	//uintptr(hwnd)
	//var lpcstr *string

	// f := syscall.NewCallback(func(parent, lpcstr, _) uintptr {

	// 	if callback(HWND(w)) {
	// 		return 1
	// 	}
	// 	return 0
	// })

	f := syscall.NewCallback(func(w HWND, lpcstr *string, handle uintptr) uintptr {
		if callback(w, lpcstr, handle) {
			return 1
		}
		return 0
	})

	ret, _, _ := enumPropsW.Call(uintptr(parent), f)

	isTrue := ret != 0
	if isTrue == true {
		return int(ret)
	}

	return -1
}

func GetWindowTextLength(hwnd HWND) int {
	ret, _, _ := getWindowTextLength.Call(uintptr(hwnd))
	return int(ret)
}

func GetWindowText(hwnd HWND) string {
	textLen := GetWindowTextLength(hwnd) + 1
	buf := make([]uint16, textLen)
	len, _, _ := getWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen),
	)
	return syscall.UTF16ToString(buf[:len])
}

func IsWindow(hwnd HWND) bool {
	ret, _, _ := isWindow.Call(uintptr(hwnd))
	return ret != 0
}

func IsIconic(hwnd HWND) bool {
	ret, _, _ := isIconic.Call(uintptr(hwnd))
	return ret != 0
}

func IsDialogMessage(hwnd HWND, msg *MSG) bool {
	ret, _, _ := isDialogMessage.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(msg)),
	)
	return ret != 0
}

func IsWindowEnabled(hwnd HWND) bool {
	ret, _, _ := isWindowEnabled.Call(uintptr(hwnd))
	return ret != 0
}

func IsWindowVisible(hwnd HWND) bool {
	ret, _, _ := isWindowVisible.Call(uintptr(hwnd))
	return ret != 0
}

func GetClassName(window HWND) (string, bool) {
	var output [256]uint16
	ret, _, _ := getClassName.Call(
		uintptr(window),
		uintptr(unsafe.Pointer(&output[0])),
		uintptr(len(output)),
	)
	return syscall.UTF16ToString(output[:]), ret != 0
}

func GetCurrentProcessId() DWORD {
	id, _, _ := getCurrentProcessId.Call()
	return DWORD(id)
}

func GetWindowRect(hwnd HWND) *RECT {
	var rect RECT
	getWindowRect.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)),
	)
	return &rect
}

func GetWindowThreadProcessId(hwnd HWND) (HANDLE, DWORD) {

	var processId DWORD
	ret, _, _ := getWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&processId)),
	)
	return HANDLE(ret), processId
}

//GetWindowInfo ,pwi *WINDOWINFO
func GetWindowInfo(hwnd HWND) (bool, WINDOWINFO) {

	var pwi WINDOWINFO
	ret, _, _ := getWindowInfo.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&pwi)),
	)
	return ret != 0, pwi
}

//GetWindowLongPtrW ..
func GetWindowLongPtrW(hwnd HWND, nIndex int) int32 {

	ret, _, _ := getWindowLongPtrW.Call(
		uintptr(hwnd),
		uintptr(nIndex),
	)

	return int32(ret)

}

type Process2 uintptr

//https://docs.microsoft.com/en-us/windows/win32/procthread/process-security-and-access-rights
const PROCESS_ALL_ACCESS = 0x1F0FFF
const PROCESS_QUERY_INFORMATION = 0x0400
const PROCESS_SET_INFORMATION = 0x0200

func OpenProcessHandle(processId int) Process2 {
	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	proc := kernel32.MustFindProc("OpenProcess")
	handle, _, _ := proc.Call(ptr(PROCESS_ALL_ACCESS), ptr(true), ptr(processId))
	return Process2(handle)
}

func ptr(val interface{}) uintptr {
	switch val.(type) {
	case string:
		return uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(val.(string))))
	case int:
		return uintptr(val.(int))
	default:
		return uintptr(0)
	}
}

// void Edit_GetText(
// 	hwndCtl,
// 	lpch,
// 	cchMax
//  );

func Edit_GetText(hwndCtl HWND, lpch *string, cchMax int) {

}

//enumPropsW
func callBackFunc(window HWND, lpcstr *string, handle uintptr) bool {

	strval := fmt.Sprintf("0x%X", lpcstr)
	//lpcstr==strval==GetWindowText()
	if len(strval) > 0 {
		_, processID := GetWindowThreadProcessId(window)

		bs, _ := hex.DecodeString(strval)
		fmt.Println("@", window, handle, lpcstr, &lpcstr, strval, processID, bs)
		return true
	}

	return false
}

func main() {

	//https://contabo.com/?show=vps
	//https://github.com/mitchellh/go-ps

	// ui, err := syscall.LoadDLL(`E:\GOLANG\src\mateors\windowsapi\process\UIAutomationClient.dll`)
	// if err != nil {
	// 	fmt.Println("ERRR::", err.Error())
	// }

	// _, err = ui.FindProc("IUIAutomationElement::FindAll")
	// fmt.Println("", ui.Name, err)

	find := uiauto.NewProc("IUIAutomationElement::FindAll")
	err := find.Find()
	fmt.Println(">>", err)

	fmt.Println(PROCESS_SET_INFORMATION, 0x00800000)

	//svar := fmt.Sprintf("%#x", 512)
	svar := fmt.Sprintf("0x%x", 512)
	fmt.Println(svar)

	EnumWindows(func(w HWND) bool {

		text := GetWindowText(w)
		isWindow := IsWindow(w)
		isEnable := IsWindowEnabled(w)
		isVisible := IsWindowVisible(w)
		//rect := GetWindowRect(w)
		//dword := GetCurrentProcessId()
		clssName, _ := GetClassName(w)
		handle, processID := GetWindowThreadProcessId(w)

		isMinimized := IsIconic(w)

		isExist, wInfo := GetWindowInfo(w)

		var msg MSG

		if isWindow && isEnable && isVisible == true {

			if isExist == true && strings.Contains(text, "Google Chrome") == true {

				res := EnumPropsW(w, callBackFunc)
				//fmt.Println("windoInfo:", windoInfo.dwStyle)
				//uint16, _ := syscall.UTF16FromString()
				//str := fmt.Sprint(wInfo.dwExStyle)
				str := fmt.Sprintf("0x%X", wInfo.dwExStyle)
				//str2 := fmt.Sprintf("0x%X", wInfo.dwStyle)

				//gwp := GetWindowLongPtrA(w, -21)
				//sss := string(gwp)
				//sss := strconv.QuoteRuneToASCII(gwp)

				fmt.Println(res, clssName, handle, processID, isMinimized, text, wInfo.dwWindowStatus, str, wInfo.dwStyle)
				//fmt.Println(wInfo.dwWindowStatus, wInfo.cxWindowBorders, wInfo.cyWindowBorders)
				//wInfo.dwExStyle, wInfo.dwStyle
				//fmt.Println()
			}
		}

		isValid := IsDialogMessage(w, &msg)
		if isValid == true {
			fmt.Println("###", msg.Hwnd, msg.LParam, msg.Message, msg.Pt, msg.Time, msg.WParam)

		}

		return true
	})

	// open process
	//https://gist.github.com/castaneai/ed8cc2aaedf9d1eafd68
	// pid := 0
	// fmt.Print("Input PID: ")
	// fmt.Scanf("%d", &pid)
	// handle := OpenProcessHandle(pid)
	// fmt.Printf("handle: %d", handle)

	// ps, err := Processes()
	// if err != nil {
	// 	fmt.Println("ERR:", err.Error())
	// }

	// var c int = 0
	// for i, p := range ps {

	// 	exeName := p.Executable()
	// 	if strings.Contains(exeName, "chrome") == true {
	// 		c++
	// 		fmt.Println(i, p.Executable(), p.PPid(), p.Pid())

	// 	}

	// }
	// fmt.Println(c)
}
