package coldpixels

import (
	"bytes"
	"strconv"
	"time"
)

type Sysinfo struct {
	display *Display
}

func NewSysinfo(d *Display) *Sysinfo {
	return &Sysinfo{
		display: d,
	}
}

func (s *Sysinfo) DisplayCPUInfo(cpuUtil, cpuTemp uint16, utilColor, tempColor TextColor) error {
	b := bytes.Buffer{}
	b.WriteString(strconv.Itoa(int(utilColor)))
	b.WriteString(strconv.Itoa(int(tempColor)))

	_, err := s.display.device.Control(0x40, 21, cpuUtil, cpuTemp, b.Bytes())

	if err != nil {
		return err
	}

	time.Sleep(s.display.displaySysinfoWait / 1000)

	return nil
}

func (s *Sysinfo) DisplayGPUInfo(ram, gpuTemp uint16, ramColor, tempColor TextColor) error {
	b := bytes.Buffer{}
	b.WriteString(strconv.Itoa(int(ramColor)))
	b.WriteString(strconv.Itoa(int(tempColor)))

	_, err := s.display.device.Control(0x40, 22, ram, gpuTemp, b.Bytes())

	if err != nil {
		return err
	}

	time.Sleep(s.display.displaySysinfoWait / 1000)

	return nil
}

func (s *Sysinfo) DisplayNetworkInfo(recv, sent uint16, recvColor, sentColor TextColor, recvMb, sentMb bool) error {
	boolMap := map[bool]int{true: 1, false: 0}

	b := bytes.Buffer{}
	b.WriteString(strconv.Itoa(int(boolMap[recvMb])))
	b.WriteString(strconv.Itoa(int(boolMap[sentMb])))
	b.WriteString(strconv.Itoa(int(recvColor)))
	b.WriteString(strconv.Itoa(int(sentColor)))

	_, err := s.display.device.Control(0x40, 20, recv, sent, b.Bytes())

	if err != nil {
		return err
	}

	time.Sleep(s.display.displaySysinfoWait / 1000)

	return nil
}

func (s *Sysinfo) DisplayFanInfo(cpu, chassis uint16, cpuColor, chassisColor TextColor) error {
	b := bytes.Buffer{}
	b.WriteString(strconv.Itoa(int(cpuColor)))
	b.WriteString(strconv.Itoa(int(chassisColor)))

	_, err := s.display.device.Control(0x40, 23, cpu, chassis, b.Bytes())

	if err != nil {
		return err
	}

	time.Sleep(s.display.displaySysinfoWait / 1000)

	return nil
}
