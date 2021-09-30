package console

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var ColorMap = map[string]string{
	"black":   "0",
	"red":     "1",
	"green":   "2",
	"yellow":  "3",
	"blue":    "4",
	"magenta": "5",
	"cyan":    "6",
	"white":   "7",
	"default": "9",
}

var BrightColor = map[string]string{
	"gray":           "0",
	"bright-red":     "1",
	"bright-green":   "2",
	"bright-yellow":  "3",
	"bright-blue":    "4",
	"bright-magenta": "5",
	"bright-cyan":    "6",
	"bright-white":   "7",
}

var OptionsMap = map[string]map[string]string{
	"bold": {
		"set":   "1",
		"unset": "22",
	},
	"underscore": {
		"set":   "4",
		"unset": "24",
	},
	"blink": {
		"set":   "5",
		"unset": "25",
	},
	"reverse": {
		"set":   "7",
		"unset": "27",
	},
	"conceal": {
		"set":   "8",
		"unset": "28",
	},
}

type Color struct {
	foreground string
	background string
	options    []string
}

func (c *Color) Set() string {
	var setCodes []string

	if len(c.foreground) > 0 {
		setCodes = append(setCodes, c.foreground)
	}
	if len(c.background) > 0 {
		setCodes = append(setCodes, c.background)
	}

	for _, option := range c.options {
		val := OptionsMap[option]["set"]
		setCodes = append(setCodes, val)
	}

	if len(setCodes) == 0 {
		return ""
	}

	return fmt.Sprintf("\033[%sm", strings.Join(setCodes, ";"))
}

func (c *Color) Unset() string {
	var unsetCodes []string

	if len(c.foreground) > 0 {
		unsetCodes = append(unsetCodes, "39")
	}

	if len(c.background) > 0 {
		unsetCodes = append(unsetCodes, "49")
	}

	for _, option := range c.options {
		val := OptionsMap[option]["unset"]
		unsetCodes = append(unsetCodes, val)
	}

	if 0 == len(unsetCodes) {
		return ""
	}

	return fmt.Sprintf("\033[%sm", strings.Join(unsetCodes, ";"))
}

func (c *Color) Apply(text string) string {
	return c.Set() + text + c.Unset()
}

func NewColor(foreground string, background string) *Color {
	foreground, _ = parseColor(foreground, false)
	background, _ = parseColor(background, true)
	color := &Color{
		foreground: foreground,
		background: background,
	}
	return color
}

func NewColorWithO(
	foreground string,
	background string,
	options []string,
) *Color {

	foreground, _ = parseColor(foreground, false)
	background, _ = parseColor(background, true)

	color := &Color{
		foreground: foreground,
		background: background,
		options:    options,
	}
	return color
}

func parseColor(color string, background bool) (string, error) {
	if len(color) == 0 {
		return "", nil
	}

	prefix := "3"
	if background {
		prefix = "4"
	}

	split := strings.Split(color, "")

	if "#" == split[0] {
		color = color[1:]

		if 3 == len(color) {
			s := strings.Split(color, "")
			color = s[0] + s[0] + s[1] + s[1] + s[2] + s[2]
		}

		if 6 != len(color) {
			return "", errors.New(fmt.Sprintf(`invalid "%s" Color.`, color))
		}

		val, err := convertHexColorToAnsi(color)
		if err != nil {
			return "", err
		}
		return prefix + val, nil
	}

	if val, ok := ColorMap[color]; ok {
		return prefix + val, nil
	}

	if val, ok := BrightColor[color]; ok {
		prefix = "9"
		if background {
			prefix = "10"
		}
		return prefix + val, nil
	}

	errMsg := fmt.Sprintf(
		`invalid "%s" Color name`,
		color,
	)
	return "", errors.New(errMsg)
}

func convertHexColorToAnsi(color string) (string, error) {
	hexdec, err := strconv.ParseInt(color, 16, 0)
	if nil != err {
		return "", err
	}

	r := (hexdec >> 16) & 255
	g := (hexdec >> 8) & 255
	b := hexdec & 255

	term := os.Getenv("COLORTERM")
	if "truecolor" != term {
		return degradeHexColorToAnsi(r, g, b)
	}

	return fmt.Sprintf("8;2;%d;%d;%d", r, g, b), nil
}

func degradeHexColorToAnsi(r int64, g int64, b int64) (string, error) {
	sat := float64(calcSaturation(r, g, b)) / 50
	if 0 == sat {
		return "", nil
	}

	r1 := math.Round(float64(b) / 255)
	r2 := math.Round(float64(g) / 255)
	r3 := math.Round(float64(r) / 255)
	retVal := int(r1)<<2 | int(r2)<<1 | int(r3)
	return strconv.Itoa(retVal), nil
}

func calcSaturation(r int64, g int64, b int64) int {
	var fr, fg, fb float64
	fr = float64(r) / 255
	fg = float64(g) / 255
	fb = float64(b) / 255
	v := math.Max(fr, fg)
	v = math.Max(v, fb)

	min := math.Min(fr, fg)
	min = math.Min(min, fb)
	diff := v - min

	if 0 == diff {
		return 0
	}

	return int(diff * 100 / v)
}
