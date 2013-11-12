package main

import (
	"os"
	"testing"
)

func assert(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("expected: %+v, actual: %+v", expected, actual)
	}
}

func (monit *Monit) getService(name string) *Service {
	for _, service := range monit.Service {
		if name == service.Name {
			return &service
		}
	}
	return nil
}

func TestMonitFromXml(t *testing.T) {
	xmlFile, err := os.Open("monit_test.xml")
	if err != nil {
		t.Fatal("Error opening file:", err)
	}
	defer xmlFile.Close()

	monit := MonitFromXml(xmlFile)

	// Monit
	assert(t, "782906abe0690f94dfce0752195640f3", monit.Id)
	assert(t, 1382272025, monit.Incarnation)
	assert(t, "5.5", monit.Version)

	// Server
	server := monit.Server
	assert(t, 576, server.Uptime)
	assert(t, 60, server.Poll)
	assert(t, 0, server.StartDelay)
	assert(t, "myhost.mydomain.tld", server.Localhostname)
	assert(t, "localhost", server.HttpAddress)
	assert(t, 2812, server.HttpPort)
	assert(t, "admin", server.Username)
	assert(t, "monit", server.Password)

	// Platform
	platform := monit.Platform
	assert(t, "Linux", platform.Name)
	assert(t, "3.11.4-1-ARCH", platform.Release)
	assert(t, "#1 SMP PREEMPT Sat Oct 5 21:22:51 CEST 2013", platform.Version)
	assert(t, "x86_64", platform.Machine)
	assert(t, 8, platform.Cpu)
	assert(t, 8081312, platform.Memory)
	assert(t, 7999484, platform.Swap)

	// Services
	assert(t, 5, len(monit.Service))

	system := *monit.getService("myhost.mydomain.tld")
	expectedSystem := Service{"myhost.mydomain.tld", SYSTEM,
		1382272565, 385252, 0, 1, 0, 0, ""}
	assert(t, expectedSystem, system)

	file := *monit.getService("bash_history")
	expectedFile := Service{"bash_history", FILE,
		1382272565, 385270, 0, 1, 0, 0, ""}
	assert(t, expectedFile, file)

	process := *monit.getService("dhcpcd")
	expectedProcess := Service{"dhcpcd", PROCESS,
		1382272565, 385323, 0, 1, 0, 0, ""}
	assert(t, expectedProcess, process)

	filesystem := *monit.getService("home")
	expectedFilesystem := Service{"home", FILESYSTEM,
		1382272565, 385470, 0, 1, 0, 0, ""}
	assert(t, expectedFilesystem, filesystem)

	directory := *monit.getService("bin")
	expectedDirectory := Service{"bin", DIRECTORY,
		1382272565, 385491, 0, 1, 0, 0, ""}
	assert(t, expectedDirectory, directory)
}
