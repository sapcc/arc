// +build windows

package process

import (
	"errors"
	"syscall"
	"unsafe"

	"github.com/shirou/w32"

	common "github.com/shirou/gopsutil/common"
	cpu "github.com/shirou/gopsutil/cpu"
	net "github.com/shirou/gopsutil/net"
)

const (
	NoMoreFiles   = 0x12
	MaxPathLength = 260
)

type SystemProcessInformation struct {
	NextEntryOffset   uint64
	NumberOfThreads   uint64
	Reserved1         [48]byte
	Reserved2         [3]byte
	UniqueProcessID   uintptr
	Reserved3         uintptr
	HandleCount       uint64
	Reserved4         [4]byte
	Reserved5         [11]byte
	PeakPagefileUsage uint64
	PrivatePageCount  uint64
	Reserved6         [6]uint64
}

// Memory_info_ex is different between OSes
type MemoryInfoExStat struct {
}

type MemoryMapsStat struct {
}

func Pids() ([]int32, error) {

	var ret []int32

	procs, err := processes()
	if err != nil {
		return ret, nil
	}
	for _, proc := range procs {
		ret = append(ret, proc.Pid)
	}

	return ret, nil
}

func (p *Process) Ppid() (int32, error) {
	ret, _, _, err := p.getFromSnapProcess(p.Pid)
	if err != nil {
		return 0, err
	}
	return ret, nil
}
func (p *Process) Name() (string, error) {
	name := ""
	return name, common.NotImplementedError
}
func (p *Process) Exe() (string, error) {
	_, _, ret, err := p.getFromSnapProcess(p.Pid)
	if err != nil {
		return "", err
	}
	return ret, nil
}
func (p *Process) Cmdline() (string, error) {
	return "", common.NotImplementedError
}
func (p *Process) Cwd() (string, error) {
	return "", common.NotImplementedError
}
func (p *Process) Parent() (*Process, error) {
	return p, common.NotImplementedError
}
func (p *Process) Status() (string, error) {
	return "", common.NotImplementedError
}
func (p *Process) Username() (string, error) {
	return "", common.NotImplementedError
}
func (p *Process) Uids() ([]int32, error) {
	var uids []int32

	return uids, common.NotImplementedError
}
func (p *Process) Gids() ([]int32, error) {
	var gids []int32
	return gids, common.NotImplementedError
}
func (p *Process) Terminal() (string, error) {
	return "", common.NotImplementedError
}
func (p *Process) Nice() (int32, error) {
	return 0, common.NotImplementedError
}
func (p *Process) IOnice() (int32, error) {
	return 0, common.NotImplementedError
}
func (p *Process) Rlimit() ([]RlimitStat, error) {
	var rlimit []RlimitStat

	return rlimit, common.NotImplementedError
}
func (p *Process) IOCounters() (*IOCountersStat, error) {
	return nil, common.NotImplementedError
}
func (p *Process) NumCtxSwitches() (*NumCtxSwitchesStat, error) {
	return nil, common.NotImplementedError
}
func (p *Process) NumFDs() (int32, error) {
	return 0, common.NotImplementedError
}
func (p *Process) NumThreads() (int32, error) {
	_, ret, _, err := p.getFromSnapProcess(p.Pid)
	if err != nil {
		return 0, err
	}
	return ret, nil
}
func (p *Process) Threads() (map[string]string, error) {
	ret := make(map[string]string, 0)
	return ret, common.NotImplementedError
}
func (p *Process) CPUTimes() (*cpu.CPUTimesStat, error) {
	return nil, common.NotImplementedError
}
func (p *Process) CPUPercent() (int32, error) {
	return 0, common.NotImplementedError
}
func (p *Process) CPUAffinity() ([]int32, error) {
	return nil, common.NotImplementedError
}
func (p *Process) MemoryInfo() (*MemoryInfoStat, error) {
	return nil, common.NotImplementedError
}
func (p *Process) MemoryInfoEx() (*MemoryInfoExStat, error) {
	return nil, common.NotImplementedError
}
func (p *Process) MemoryPercent() (float32, error) {
	return 0, common.NotImplementedError
}

func (p *Process) Children() ([]*Process, error) {
	return nil, common.NotImplementedError
}

func (p *Process) OpenFiles() ([]OpenFilesStat, error) {
	return nil, common.NotImplementedError
}

func (p *Process) Connections() ([]net.NetConnectionStat, error) {
	return nil, common.NotImplementedError
}

func (p *Process) IsRunning() (bool, error) {
	return true, common.NotImplementedError
}

func (p *Process) MemoryMaps(grouped bool) (*[]MemoryMapsStat, error) {
	ret := make([]MemoryMapsStat, 0)
	return &ret, common.NotImplementedError
}

func NewProcess(pid int32) (*Process, error) {
	p := &Process{Pid: pid}

	return p, nil
}

func (p *Process) SendSignal(sig syscall.Signal) error {
	return common.NotImplementedError
}

func (p *Process) Suspend() error {
	return common.NotImplementedError
}
func (p *Process) Resume() error {
	return common.NotImplementedError
}
func (p *Process) Terminate() error {
	return common.NotImplementedError
}
func (p *Process) Kill() error {
	return common.NotImplementedError
}

func (p *Process) getFromSnapProcess(pid int32) (int32, int32, string, error) {
	snap := w32.CreateToolhelp32Snapshot(w32.TH32CS_SNAPPROCESS, uint32(pid))
	if snap == 0 {
		return 0, 0, "", syscall.GetLastError()
	}
	defer w32.CloseHandle(snap)
	var pe32 w32.PROCESSENTRY32
	pe32.DwSize = uint32(unsafe.Sizeof(pe32))
	if w32.Process32First(snap, &pe32) == false {
		return 0, 0, "", syscall.GetLastError()
	}

	if pe32.Th32ProcessID == uint32(pid) {
		szexe := syscall.UTF16ToString(pe32.SzExeFile[:])
		return int32(pe32.Th32ParentProcessID), int32(pe32.CntThreads), szexe, nil
	}

	for w32.Process32Next(snap, &pe32) {
		if pe32.Th32ProcessID == uint32(pid) {
			szexe := syscall.UTF16ToString(pe32.SzExeFile[:])
			return int32(pe32.Th32ParentProcessID), int32(pe32.CntThreads), szexe, nil
		}
	}
	return 0, 0, "", errors.New("Couldn't find pid:" + string(pid))
}

// Get processes
func processes() ([]*Process, error) {
	ps := make([]uint32, 255)
	var read uint32
	if w32.EnumProcesses(ps, uint32(len(ps)), &read) == false {
		return nil, syscall.GetLastError()
	}

	var results []*Process
	dwardSize := uint32(4)
	for _, pid := range ps[:read/dwardSize] {
		if pid == 0 {
			continue
		}
		p, err := NewProcess(int32(pid))
		if err != nil {
			break
		}
		results = append(results, p)
	}

	return results, nil
}

func getProcInfo(pid int32) (*SystemProcessInformation, error) {
	initialBufferSize := uint64(0x4000)
	bufferSize := initialBufferSize
	buffer := make([]byte, bufferSize)

	var sysProcInfo SystemProcessInformation
	ret, _, _ := common.ProcNtQuerySystemInformation.Call(
		uintptr(unsafe.Pointer(&sysProcInfo)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bufferSize)),
		uintptr(unsafe.Pointer(&bufferSize)))
	if ret != 0 {
		return nil, syscall.GetLastError()
	}

	return &sysProcInfo, nil
}
