package memory

import "github.com/shirou/gopsutil/mem"

type Source struct{}

func New() Source {
	return Source{}
}

func (h Source) Name() string {
	return "memory"
}

func (h Source) Facts() (map[string]interface{}, error) {
	facts := make(map[string]interface{})
	m, err := mem.VirtualMemory()
	if err != nil {
		return facts, nil
	}

	facts["memory_total"] = m.Total
	facts["memory_used"] = m.Used
	facts["memory_used_percent"] = int(m.UsedPercent + .5)
	facts["memory_available"] = m.Available

	return facts, nil
}
