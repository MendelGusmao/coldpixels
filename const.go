package coldpixels

type (
	TextColor       uint
	BackgroundColor uint16
	TextAlignment   int16
	TextLine        int
)

const (
	vendorID, deviceID uint16 = 0x16c0, 0x05dc

	TextColorGreen      TextColor = 1
	TextColorYellow     TextColor = 2
	TextColorRed        TextColor = 3
	TextColorWhite      TextColor = 5
	TextColorCyan       TextColor = 6
	TextColorGrey       TextColor = 7
	TextColorBlack      TextColor = 13
	TextColorBrown      TextColor = 15
	TextColorBrickRed   TextColor = 16
	TextColorDarkBlue   TextColor = 17
	TextColorLightBlue  TextColor = 18
	TextColorOrange     TextColor = 21
	TextColorPurple     TextColor = 22
	TextColorPink       TextColor = 23
	TextColorPeach      TextColor = 24
	TextColorGold       TextColor = 25
	TextColorLavender   TextColor = 26
	TextColorOrangeRed  TextColor = 27
	TextColorMagenta    TextColor = 28
	TextColorNavy       TextColor = 30
	TextColorLightGreen TextColor = 31

	BackgroundColorBlack     BackgroundColor = 0x0000
	BackgroundColorBlue      BackgroundColor = 0x001f
	BackgroundColorGreen     BackgroundColor = 0x07e0
	BackgroundColorCyan      BackgroundColor = 0x07ff
	BackgroundColorBrown     BackgroundColor = 0x79e0
	BackgroundColorDarkGrey  BackgroundColor = 0x7bef
	BackgroundColorLightGrey BackgroundColor = 0xbdf7
	BackgroundColorRed       BackgroundColor = 0xf800
	BackgroundColorPurple    BackgroundColor = 0xf81f
	BackgroundColorOrange    BackgroundColor = 0xfbe0
	BackgroundColorGold      BackgroundColor = 0xfd20
	BackgroundColorYellow    BackgroundColor = 0xffe0
	BackgroundColorWhite     BackgroundColor = 0xffff

	TextAlignmentNone   TextAlignment = -1
	TextAlignmentCentre TextAlignment = 0
	TextAlignmentLeft   TextAlignment = 1
	TextAlignmentRight  TextAlignment = 2

	TextLine1   TextLine = 1 << 0
	TextLine2   TextLine = 1 << 1
	TextLine3   TextLine = 1 << 2
	TextLine4   TextLine = 1 << 3
	TextLine5   TextLine = 1 << 4
	TextLine6   TextLine = 1 << 5
	TextLineAll TextLine = TextLine1 + TextLine2 + TextLine3 + TextLine4 + TextLine5 + TextLine6
)
