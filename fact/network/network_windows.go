package network

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/sapcc/arc/fact/helper"
)

func (h Source) Facts() (map[string]interface{}, error) {
	facts := newFacts()

	adapters, _ := getAdapterList()
	for ; adapters != nil; adapters = adapters.Next {
		//name := helper.BytePtrToString(&adapters.Description[0])
		gw := helper.BytePtrToString(&adapters.GatewayList.IpAddress.String[0])
		if gw != "0.0.0.0" {
			facts["default_gateway"] = gw
			facts["ipaddress"] = helper.BytePtrToString(&adapters.IpAddressList.IpAddress.String[0])
			facts["default_interface"] = fmt.Sprintf("%d", adapters.Index)
		}

	}

	return facts, nil
}

// borrowed from src/pkg/net/interface_windows.go
func getAdapterList() (*syscall.IpAdapterInfo, error) {
	b := make([]byte, 1000)
	l := uint32(len(b))
	a := (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0])) // #nosec
	err := syscall.GetAdaptersInfo(a, &l)
	if err == syscall.ERROR_BUFFER_OVERFLOW {
		b = make([]byte, l)
		a = (*syscall.IpAdapterInfo)(unsafe.Pointer(&b[0])) // #nosec
		err = syscall.GetAdaptersInfo(a, &l)
	}
	if err != nil {
		return nil, os.NewSyscallError("GetAdaptersInfo", err)
	}
	return a, nil
}
