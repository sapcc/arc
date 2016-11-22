package network

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"regexp"
)

var (
	ipIntRegex *regexp.Regexp = regexp.MustCompile(`^(\d+): ([0-9a-zA-Z@:\.\-_]*?)(@[0-9a-zA-Z]+|):\s`)
	linkRegex  *regexp.Regexp = regexp.MustCompile(`link\/(\w+) ([\da-f\:]+) `)
	inetRegex  *regexp.Regexp = regexp.MustCompile(`inet (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(\/(\d{1,2}))?`)

	defaultGwRegex *regexp.Regexp = regexp.MustCompile(`default via ([^\s]+) dev ([^\s]+)`)
)

type interf struct {
	Name string
	IPs  []string
	Type string
	Mac  string
}

func ipBinary() string {
	//possible_locations := []string{"/sbin/ip", "/usr/bin/ip", "/bin/ip"}
	for _, path := range []string{"/sbin/ip", "/usr/bin/ip", "/bin/ip"} {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ""
}

func (h Source) Facts() (map[string]interface{}, error) {

	facts := newFacts()
	cmd := exec.Command(ipBinary(), "addr") // #nosec
	interfaces := make(map[string]*interf)
	var currentInterface string
	if out, err := cmd.Output(); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		//fmt.Println(string(out))
		for scanner.Scan() {
			line := scanner.Text()
			if match := ipIntRegex.FindStringSubmatch(line); match != nil {
				//fmt.Println("ipIntRegex match", match[0])
				currentInterface = match[2]
				interfaces[currentInterface] = &interf{Name: currentInterface, IPs: make([]string, 0)}

			} else if match := linkRegex.FindStringSubmatch(line); match != nil {
				i := interfaces[currentInterface]
				i.Type = match[1]
				i.Mac = match[2]

			} else if match := inetRegex.FindStringSubmatch(line); match != nil {
				//fmt.Println("inet Regex match", match[0])
				i := interfaces[currentInterface]
				i.IPs = append(i.IPs, match[1])
			}
		}

	}
	//for _, i := range interfaces {
	//  fmt.Println("interface", i)
	//}

	cmd = exec.Command(ipBinary(), "-o", "-f", "inet", "route", "show") // #nosec
	if out, err := cmd.Output(); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(out))
		for scanner.Scan() {
			line := scanner.Text()
			if match := defaultGwRegex.FindStringSubmatch(line); match != nil {
				facts["default_interface"] = match[2]
				facts["ipaddress"] = interfaces[match[2]].IPs[0]
				facts["macaddress"] = interfaces[match[2]].Mac
				facts["default_gateway"] = match[1]
			}
		}
	}

	//If no default gatway can be found fall back to the first ethernet interface
	if ip := facts["ipaddress"]; ip == nil {
		for name, i := range interfaces {
			if i.Type == "ether" {
				facts["ipaddress"] = interfaces[name].IPs[0]
				facts["macaddress"] = interfaces[name].Mac
				break
			}
		}
	}

	return facts, nil

}
