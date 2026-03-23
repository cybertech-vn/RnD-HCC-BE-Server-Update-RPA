package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/denisbrodbeck/machineid"
	"github.com/jaypipes/ghw"
)

func GenerateFingerprint() string {
	// Lấy machine ID làm base
	machineID, err := machineid.ProtectedID("MayBocSo")
	if err != nil {
		machineID = "NO_MACHINE_ID"
	}

	// Lấy CPU info (combine vendor + model + cores)
	cpuInfo := getCPUInfo()

	// Lấy disk serial (serial của disk đầu tiên)
	diskSerial := getDiskSerial()

	// Lấy baseboard serial (mainboard serial)
	baseboardSerial := getBaseboardSerial()

	// Lấy product serial (system serial, gần như BIOS UUID)
	productSerial := getProductSerial()

	// Combo tất cả
	raw := fmt.Sprintf("%s|%s|%s|%s|%s", machineID, cpuInfo, diskSerial, baseboardSerial, productSerial)

	// Hash SHA256
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func getCPUInfo() string {
	cpu, err := ghw.CPU(ghw.WithDisableWarnings())
	if err != nil {
		return "NO_CPU"
	}
	if len(cpu.Processors) > 0 {
		proc := cpu.Processors[0]
		return fmt.Sprintf("%s-%s-%d", proc.Vendor, proc.Model, proc.NumCores)
	}
	return "NO_CPU"
}

func getDiskSerial() string {
	block, err := ghw.Block(ghw.WithDisableWarnings())
	if err != nil {
		return "NO_DISK"
	}
	if len(block.Disks) > 0 {
		disk := block.Disks[0] // Disk đầu tiên (thường là boot disk)
		if disk.SerialNumber != "" && disk.SerialNumber != "unknown" {
			return disk.SerialNumber
		}
	}
	return "NO_DISK"
}

func getBaseboardSerial() string {
	baseboard, err := ghw.Baseboard(ghw.WithDisableWarnings())
	if err != nil {
		return "NO_BASEBOARD"
	}
	if baseboard.SerialNumber != "" && baseboard.SerialNumber != "unknown" {
		return baseboard.SerialNumber
	}
	return "NO_BASEBOARD"
}

func getProductSerial() string {
	product, err := ghw.Product(ghw.WithDisableWarnings())
	if err != nil {
		return "NO_PRODUCT"
	}
	if product.SerialNumber != "" && product.SerialNumber != "unknown" {
		return product.SerialNumber
	}
	return "NO_PRODUCT"
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
