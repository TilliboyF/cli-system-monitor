package main

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

var cpuNow float64
var networkIn float64
var networkOut float64

func main() {
	go getCpu()
	go getNetwork()

	for {
		fmt.Print("\033[2J")
		fmt.Print("\033[H")

		fmt.Println("Current Time:", time.Now().Format("15:04:05"))
		memory, _ := mem.VirtualMemory()

		fmt.Printf("Total: %.2f Gb, Used: %.2f Gb, UsedPercent: %.2f%%\n", float64(memory.Total)/(1000*1000*1000), float64(memory.Used)/(1000*1000*1000), memory.UsedPercent)

		fmt.Printf("CPU Usage: %.2f%%\n", cpuNow)

		fmt.Printf("Network - IN: %.2f KB, OUT: %.2f KB\n", networkOut, networkIn)

		time.Sleep(1 * time.Second)
	}
}

func getCpu() {
	for {
		cpu, _ := cpu.Percent(time.Second, false)
		cpuNow = cpu[0]
	}
}

func getNetwork() {
	previousStats, _ := net.IOCounters(true)
	for {
		time.Sleep(1 * time.Second)
		currentStats, err := net.IOCounters(true)
		if err != nil {
			fmt.Println("Error fetching network stats:", err)
			return
		}

		var totalBytesSent uint64
		var totalBytesRecv uint64

		for i, stat := range currentStats {
			totalBytesSent += stat.BytesSent - previousStats[i].BytesSent
			totalBytesRecv += stat.BytesRecv - previousStats[i].BytesRecv
		}

		networkIn = float64(totalBytesRecv) / 8
		networkOut = float64(totalBytesSent) / 8

		previousStats = currentStats // update the previous stats for the next iteration
	}
}
