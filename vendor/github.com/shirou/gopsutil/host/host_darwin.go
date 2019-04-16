// +build darwin

package host

import (
	"bytes"
	"encoding/binary"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"unsafe"

	common "github.com/shirou/gopsutil/common"
)

const (
	UTXUserSize = 256 /* include/NetBSD/utmpx.h */
	UTXIDSize   = 4
	UTXLineSize = 32
	UTXHostSize = 256
)

type utmpx32 struct {
	UtUser [UTXUserSize]byte /* login name */
	UtID   [UTXIDSize]byte   /* id */
	UtLine [UTXLineSize]byte /* tty name */
	//TODO	UtPid  pid_t              /* process id creating the entry */
	UtType [4]byte /* type of this entry */
	//TODO	UtTv   timeval32          /* time entry was created */
	UtHost [UTXHostSize]byte /* host name */
	UtPad  [16]byte          /* reserved for future use */
}

func HostInfo() (*HostInfoStat, error) {
	ret := &HostInfoStat{
		OS:             runtime.GOOS,
		PlatformFamily: "darwin",
	}

	hostname, err := os.Hostname()
	if err != nil {
		return ret, err
	}
	ret.Hostname = hostname

	platform, family, version, err := GetPlatformInformation()
	if err == nil {
		ret.Platform = platform
		ret.PlatformFamily = family
		ret.PlatformVersion = version
	}
	system, role, err := GetVirtualization()
	if err == nil {
		ret.VirtualizationSystem = system
		ret.VirtualizationRole = role
	}

	values, err := common.DoSysctrl("kern.boottime")
	if err == nil {
		// ex: { sec = 1392261637, usec = 627534 } Thu Feb 13 12:20:37 2014
		v := strings.Replace(values[2], ",", "", 1)
		t, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return ret, err
		}
		ret.Uptime = t
	}

	return ret, nil
}

func BootTime() (int64, error) {
	values, err := common.DoSysctrl("kern.boottime")
	if err != nil {
		return 0, err
	}
	// ex: { sec = 1392261637, usec = 627534 } Thu Feb 13 12:20:37 2014
	v := strings.Replace(values[2], ",", "", 1)

	boottime, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, err
	}

	return boottime, nil
}

func Users() ([]UserStat, error) {
	utmpfile := "/var/run/utmpx"
	var ret []UserStat

	file, err := os.Open(utmpfile)
	if err != nil {
		return ret, err
	}

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return ret, err
	}

	u := utmpx32{}
	entrySize := int(unsafe.Sizeof(u))
	count := len(buf) / entrySize

	for i := 0; i < count; i++ {
		b := buf[i*entrySize : i*entrySize+entrySize]

		var u utmpx32
		br := bytes.NewReader(b)
		err := binary.Read(br, binary.LittleEndian, &u)
		if err != nil {
			continue
		}
		user := UserStat{
			User: common.ByteToString(u.UtUser[:]),
			//			Terminal: ByteToString(u.UtLine[:]),
			Host: common.ByteToString(u.UtHost[:]),
			//			Started:  int(u.UtTime),
		}
		ret = append(ret, user)
	}

	return ret, nil

}

func GetPlatformInformation() (string, string, string, error) {
	platform := ""
	family := ""
	version := ""

	out, err := exec.Command("uname", "-s").Output()
	if err == nil {
		platform = strings.ToLower(strings.TrimSpace(string(out)))
	}

	out, err = exec.Command("uname", "-r").Output()
	if err == nil {
		version = strings.ToLower(strings.TrimSpace(string(out)))
	}

	return platform, family, version, nil
}

func GetVirtualization() (string, string, error) {
	system := ""
	role := ""

	return system, role, nil
}
