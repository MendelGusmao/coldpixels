package coldpixels

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/kylelemons/gousb/usb"
)

const (
	defaultUSBTimeout         = 5000 * time.Millisecond
	defaultClearLineWait      = 1000 * time.Millisecond
	defaultEraseSectorWait    = 220 * time.Millisecond
	defaultMinDisplayIconWait = 10 * time.Millisecond
	defaultMaxDisplayIconWait = 700 * time.Millisecond
	defaultWritePageWait      = 15 * time.Millisecond
	defaultDisplaySysinfoWait = 50 * time.Millisecond
	defaultMaxDisplayTextWait = 85 * time.Millisecond
	defaultCharsPerIcon       = 2.75
	col2Left                  = "|||||||||___"

	ctrlGetDeviceInfo          = 12
	ctrlSetBrightness          = 13
	ctrlSaveBrightness         = 14
	ctrlSendCommandToFlash     = 15
	ctrlRawImageToFlashSend    = 16
	ctrlDimWhenIdle            = 17
	ctrlDisplayNetworkInfo     = 20
	ctrlDisplayCPUInfo         = 21
	ctrlDisplayRAMGPUInfo      = 22
	ctrlDisplayFanInfo         = 23
	ctrlDisplayTextOnLine      = 24
	ctrlDisplayTextAnywhere    = 25
	ctrlClearLines             = 26
	ctrlDisplayIcon            = 27
	ctrlDisplayIconAnywhere    = 29
	ctrlSetTextBackgroundColor = 30
)

var (
	largeImageIndexes = map[uint16]bool{
		180: true,
		218: true,
		256: true,
		294: true,
		332: true,
		370: true,
		408: true,
		446: true,
	}

	fontLengthTable = []uint16{
		0x11, 0x06, 0x08, 0x15, 0x0E, 0x19, 0x15, 0x03, 0x08, 0x08, 0x0F, 0x0D,
		0x05, 0x08, 0x06, 0x0B, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11, 0x11,
		0x11, 0x11, 0x06, 0x06, 0x13, 0x10, 0x13, 0x0C, 0x1A, 0x14, 0x10, 0x12,
		0x13, 0x0F, 0x0D, 0x13, 0x11, 0x04, 0x07, 0x11, 0x0E, 0x14, 0x11, 0x15,
		0x0F, 0x15, 0x12, 0x10, 0x13, 0x11, 0x14, 0x1C, 0x13, 0x13, 0x12, 0x07,
		0x0B, 0x07, 0x0B, 0x02, 0x08, 0x0E, 0x0F, 0x0E, 0x0F, 0x10, 0x0B, 0x0F,
		0x0E, 0x04, 0x07, 0x0F, 0x04, 0x18, 0x0E, 0x10, 0x0F, 0x0F, 0x0A, 0x0D,
		0x0B, 0x0E, 0x10, 0x16, 0x10, 0x10, 0x0E, 0x01, 0x11, 0x02,
	}
)

type Display struct {
	usbTimeout         time.Duration
	clearLineWait      time.Duration
	eraseSectorWait    time.Duration
	minDisplayIconWait time.Duration
	maxDisplayIconWait time.Duration
	writePageWait      time.Duration
	displaySysinfoWait time.Duration
	maxDisplayTextWait time.Duration
	charsPerIcon       float32

	device *usb.Device
}

type DeviceInfo struct {
	EEPROM           string
	Serial           string
	FlashId          int
	FlashData        []byte
	DeviceValid      bool
	PictureFrameMode bool
	EightMbFlash     bool
	Flashcap         string
	FirmwareVersion  float32
	FlashDataVersion float32
}

func NewDisplay(index int) *Display {
	return &Display{
		usbTimeout:         defaultUSBTimeout,
		clearLineWait:      defaultClearLineWait,
		eraseSectorWait:    defaultEraseSectorWait,
		minDisplayIconWait: defaultMinDisplayIconWait,
		maxDisplayIconWait: defaultMaxDisplayIconWait,
		writePageWait:      defaultWritePageWait,
		displaySysinfoWait: defaultDisplaySysinfoWait,
		maxDisplayTextWait: defaultMaxDisplayTextWait,
		charsPerIcon:       defaultCharsPerIcon,
	}
}

func (s *Display) OpenDevice() error {
	devices, err := usb.NewContext().ListDevices(func(desc *usb.Descriptor) bool {
		return desc.Vendor == usb.ID(vendorID) && desc.Device == usb.BCD(deviceID)
	})

	if err != nil {
		return err
	}

	if len(devices) == 0 {
		return fmt.Errorf("no lcd detected")
	}

	s.device = devices[0]

	return nil
}

func (s *Display) SetBrightness(value uint16) error {
	// ignoring usbTimeout???
	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlSetBrightness, value%255, value%255, []byte{})

	return err
}

func (s *Display) SaveBrightness(off, on uint16) error {
	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlSaveBrightness, off+on*256, 0, []byte{})

	return err
}

func (s *Display) DisplayIcon(position, iconNumber uint16) error {
	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlDisplayIconAnywhere, (position%47)*512+iconNumber, 25600, []byte{})

	s.sleepByIconNumber(iconNumber)

	return err
}

func (s *Display) DisplayIconAnywhere(x, y, iconNumber uint16) error {
	x = x % 320
	y = y % 240

	b := bytes.Buffer{}
	binary.Write(&b, binary.BigEndian, y>>8)
	binary.Write(&b, binary.BigEndian, y&0xFF)
	binary.Write(&b, binary.BigEndian, x>>8)
	binary.Write(&b, binary.BigEndian, x&0xFF)

	idx := (iconNumber << 8) + iconNumber

	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlDisplayIconAnywhere, idx, idx, b.Bytes())

	if err != nil {
		return err
	}

	s.sleepByIconNumber(iconNumber)

	return nil
}

func (s *Display) SetTextBackgroundColor(color uint16) error {
	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlSetTextBackgroundColor, color, 0, []byte{})

	return err
}

// TODO
func (s *Display) alignText(text string, alignment TextAlignment, screenPx, stringLengthPx uint16) string {
	spaces := math.Floor(float64((screenPx - stringLengthPx) / 17))
	pixels := math.Mod(float64(screenPx-stringLengthPx), 17)

	switch alignment {
	case TextAlignmentCentre:
		text = ""
		// mm = mm.center(len(mm) + spaces, " ").center(len(mm) + pixels, "{")
	case TextAlignmentLeft:
		text = ""
		// mm = mm + " " * spaces + "{" * pixels
	case TextAlignmentRight:
		text = ""
		// mm = " " * spaces + "{" * pixels + mm
	}

	fmt.Println(pixels, spaces)
	return text
}

func (s *Display) textConversion(text string, fieldSize uint16, alignment TextAlignment) string {
	screenPx := 40*fieldSize - 1
	text = strings.Replace(text, " ", "___", -1) // TODO strip
	text = strings.TrimSpace(text)
	stringLengthPx := uint16(0)

	for i := range text {
		ord := int((text)[i])
		charLengthPx := uint16(0)

		if ord >= 32 && ord <= 125 {
			charLengthPx = fontLengthTable[ord-32]
		}

		if stringLengthPx+charLengthPx > screenPx {
			text = (text)[0:i]
		}

		stringLengthPx += charLengthPx
	}

	return s.alignText(text, alignment, screenPx, stringLengthPx)
}

func (s *Display) DisplayTextOnLine(line uint64, text string, padForIcon bool, alignment IntSlice, color TextColor, fieldLength Uint16Slice) error {
	buffer := &bytes.Buffer{}
	div := uint16(8)

	if padForIcon {
		div = 7
	}

	fieldLength.Map(func(l uint16) uint16 {
		return l % div
	})

	hasTab := strings.Contains(text, "\t")

	if hasTab {
		fieldLength.Set(4)

		if padForIcon {
			fieldLength.Set(3)
		}

		parts := strings.SplitN(text, "\t", 2)
		parts[1] = strings.Replace(parts[1], "\t", "", -1)

		if padForIcon {
			fieldLength.Insert(1, 1)
			alignment.Insert(1, int(TextAlignmentLeft))

			parts = append(parts, "")
			copy(parts[2:], parts[1:])
			parts[1] = ""
		}

		for i := range parts {
			buffer.WriteString(s.textConversion(parts[i], fieldLength[i], TextAlignment(alignment[i])))

			if i < len(parts)-1 {
				buffer.WriteByte(0)
			}
		}
	} else {
		buffer.WriteString(s.textConversion(text, fieldLength[0], TextAlignment(alignment[0])))
		buffer.WriteByte(0)
	}

	textLength := uint16(len(buffer.Bytes()))

	if !padForIcon {
		textLength += 256
	}

	color = color % 32
	line = line % 6

	if line == 0 {
		line++
	}

	// TODO encode

	index := uint16((line-1)*256) + uint16(color)
	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlDisplayTextOnLine, textLength, index, buffer.Bytes())

	if err != nil {
		return err
	}

	length := uint16(0)

	if hasTab {
		length = fieldLength.Sum()
	} else {
		length = fieldLength[0]
	}

	lengthPercentage := float32(length) * float32(s.charsPerIcon) / 22.0
	duration := float32(s.maxDisplayTextWait) * lengthPercentage / 1000
	time.Sleep(time.Duration(duration))

	return nil
}

func (s *Display) DisplayTextAnywhere(x, y uint16, text string, color TextColor) error {
	x = x % 320
	y = y % 320
	y2 := y + 40

	b := bytes.Buffer{}
	binary.Write(&b, binary.BigEndian, x>>8)
	binary.Write(&b, binary.BigEndian, x&0xFF)
	binary.Write(&b, binary.BigEndian, y>>8)
	binary.Write(&b, binary.BigEndian, y&0xFF)

	binary.Write(&b, binary.BigEndian, 319>>8)
	binary.Write(&b, binary.BigEndian, 319&0xFF)
	binary.Write(&b, binary.BigEndian, y2>>8)
	binary.Write(&b, binary.BigEndian, y2&0xFF)
	b.WriteString(text)

	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlDisplayTextAnywhere, uint16(len(b.Bytes())), uint16(color), b.Bytes())

	if err != nil {
		return err
	}

	time.Sleep(s.maxDisplayTextWait / 1000)

	return nil
}

func (s *Display) DimWhenIdle(dim bool) error {
	var err error

	if dim {
		_, err = s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlDimWhenIdle, 1, 0, []byte{})
	} else {
		_, err = s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlDimWhenIdle, 0, 266, []byte{})
	}

	return err
}

func (s *Display) ClearLines(lines uint16, color BackgroundColor) error {
	lines = lines % 63

	if lines == 0 {
		lines++
	}

	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlClearLines, lines, uint16(color), []byte{})

	if err != nil {
		return err
	}

	duration := float32(s.clearLineWait) * float32(s.countBitsSet(lines)/6) * 0.8 / 1000
	time.Sleep(time.Duration(duration))

	return nil
}

func (s *Display) SendCommandToFlash(address, command uint16) error {
	_, err := s.device.Control(usb.REQUEST_TYPE_VENDOR, ctrlSendCommandToFlash, address, command, []byte{})

	return err
}

// TODO
func (s *Display) WriteRawImageToFlash(sector uint16, raw []byte, checkSizes bool) error {
	return nil
}

// TODO
func (s *Display) WriteImageToFlash(sector uint16, bitmap []byte, checkSizes bool) error {
	return nil
}

// TODO
func (s *Display) GetDeviceInfo() (*DeviceInfo, error) {
	return nil, nil
}

func (s *Display) Release() error {
	return s.device.Close()
}

func (s *Display) sleepByIconNumber(iconNumber uint16) {
	if _, ok := largeImageIndexes[iconNumber]; ok {
		time.Sleep(s.maxDisplayIconWait)
	} else {
		time.Sleep(s.minDisplayIconWait)
	}
}

func (s *Display) countBitsSet(value uint16) uint16 {
	count := uint16(0)

	for value > 0 {
		value &= value - 1
		count++
	}

	return count
}
