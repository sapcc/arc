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

	facts["memory_total"] = fmt.Sprintf("%d", m.Total)
	facts["memory_used"] = fmt.Sprintf("%d", m.Used)
	facts["memory_used_percent"] = fmt.Sprintf("%d", int(m.UsedPercent+.5))
	facts["memory_available"] = fmt.Sprintf("%d", m.Available)

	return facts, nil
}
