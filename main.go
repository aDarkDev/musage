package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getChilds(pid int) []int {
	var checked []int
	var allPids []int
	allPids = append(allPids, pid)
	checked = append(checked, pid)

	for i := 0; i < len(allPids); {
		listChild, _ := exec.Command("pgrep", "-P", strconv.Itoa(allPids[i])).Output()
		childPidsStr := strings.Split(string(listChild), "\n")

		for _, childPidStr := range childPidsStr {
			childPid, err := strconv.Atoi(childPidStr)
			if err == nil && !contains(checked, childPid) {
				checked = append(checked, childPid)
				allPids = append(allPids, childPid)
			}
		}
		i++
	}
	return allPids
}

func getPidMemory(pid int) float64 {
	regex := regexp.MustCompile(`VmRSS:\s+(\d+)\s+`)
	dataProc, err := os.ReadFile(fmt.Sprintf("/proc/%d/status", pid))
	if err != nil {
		return 0
	}
	rssMemoryRaw := regex.FindStringSubmatch(string(dataProc))
	if len(rssMemoryRaw) == 0 {
		return 0
	}
	rssMemory := rssMemoryRaw[1]
	memory, err := strconv.Atoi(rssMemory)
	if err != nil {
		return 0
	}
	return float64(memory) / 1024
}

func getFullPidMemory(pid int) float64 {
	pids := getChilds(pid)
	total := 0.
	for _, i := range pids {
		total += getPidMemory(i)
	}
	return total
}

func memory_writer(pid int, command string) {
	for {
		usage_mb := getFullPidMemory(pid)
		timestamp := time.Now().Unix()
		text := fmt.Sprintf(
			"{\"process_id\": %d,\"memory_usage\": %.6f, \"timestamp\": %d, \"command\": \"%s\"}",
			pid,
			usage_mb,
			timestamp,
			strings.Replace(command, "\"", "\\\"", -1),
		)
		os.WriteFile("musage.log", []byte(text), 0644)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Printf("usage: %s <command>\r\nexample: %s sleep 10; echo hi", os.Args[0], os.Args[0])
		return
	}

	torun := strings.Join(os.Args[1:], " ")
	cmd := exec.Command("bash", "-c", torun)
	go cmd.Run()
	time.Sleep(200*time.Millisecond)
	go memory_writer(cmd.Process.Pid, torun)
	cmd.Wait()
}
