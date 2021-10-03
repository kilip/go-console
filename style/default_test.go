package style

import (
	qt "github.com/frankban/quicktest"
	"testing"
)

func TestNewDefaultStyle(t *testing.T) {
	c := qt.New(t)
	buffMock := new(buffMock)
	ds := NewDefaultStyle(buffMock, buffMock)

	ds.Write("hello world")
	c.Assert(buffMock.Output, qt.Equals, "hello world")
}

type outputTestCase struct {
	Name,
	Expected,
	Input string
	Style func(testCase outputTestCase, ds *DefaultStyle)
}

func getBlockTestCase() []outputTestCase {
	return []outputTestCase{
		{
			Name: "block",
			Expected: `
X [CUSTOM] Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et      
X          dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea 
X          commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat    
X          nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit    
X          anim id est laborum                                                                                          


`,
			Input: "Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum",
			Style: func(testCase outputTestCase, ds *DefaultStyle) {
				ds.GetFormatter().SetDecorated(false)
				ds.BlockO(
					testCase.Input,
					"CUSTOM",
					"fg=white;bg=green",
					"X ",
					true,
					true,
				)
			},
		},
		{
			Name:  "block",
			Input: "Lopadotemachoselachogaleokranioleipsanodrimhypotrimmatosilphioparaomelitokatakechymenokichlepikossyphophattoperisteralektryonoptekephalliokigklopeleiolagoiosiraiobaphetraganopterygon",
			Expected: `
 § [CUSTOM] Lopadotemachoselachogaleokranioleipsanodrimhypotrimmatosilphioparaomelitokatakechymenokichlepikossyphophat 
 §          toperisteralektryonoptekephalliokigklopeleiolagoiosiraiobaphetraganopterygon                               


`,
			Style: func(testCase outputTestCase, ds *DefaultStyle) {
				ds.GetFormatter().SetDecorated(false)
				ds.BlockO(
					testCase.Input,
					"CUSTOM",
					"fg=blue;bg=blue",
					" § ",
					false,
					true,
				)
			},
		},
	}
}

func TestDefaultStyle_Output(t *testing.T) {
	cases := getBlockTestCase()

	for _, tCase := range cases {
		t.Run(tCase.Name, func(t *testing.T) {
			c := qt.New(t)
			buffMock := new(buffMock)
			ds := NewDefaultStyle(buffMock, buffMock)
			tCase.Style(tCase, ds)
			c.Assert(buffMock.Output, qt.Equals, tCase.Expected)
		})
	}
}
