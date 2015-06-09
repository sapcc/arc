package memory

import (
	"fmt"

	"github.com/shirou/gopsutil/mem"
)

type Source struct{}

func New() Source {
	return Source{}
}

func (h Source) Name() string {
	return "memory"
}

func (h Source) Facts() (map[string]string, error) {
	facts := make(map[string]string)
	m, _ := mem.VirtualMemory()

	fmt.Println("mem", m)
	return facts, nil
}
