package main

import (
	"github.com/josepuga/goini"
	"testing"
)

var content = `
# This is a common comment
; This is a not so common comment

info text=This is a key/value inside [] empty section

[gui settings]
width = 1920
height=720
scale factor=1.33
scale factor2=THIS IS NOT A FLOAT VALUE
valid themes=dark,light,awaita,classic,aqua

[theme]
# You can use 0/1, true/false
use system theme=0
accent color= 0xff00ff
`

func TestFromByte(t *testing.T) {
	ini := goini.NewIni()
	ini.LoadFromBytes([]byte(content))

	sf := ini.GetFloat("gui settings", "scale factor", 2.5)
	t.Logf("Scale factor should be 1.33: %f\n", sf)
	sf = ini.GetFloat("gui settings", "scale factor2", 2.5)
	t.Logf("Scale factor2 should be 2.5: %f\n", sf)
	t.Logf("info text is %s\n", ini.GetString("", "info text", "default value"))
	if ini.GetBool("theme", "use system theme", true) == true {
		t.Error("use system theme must be false")
	}
	if ini.GetInt("theme", "accent color", 1000) != 16_711_935 {
		t.Error("color must be 0xff00ff")
	}
	themes := ini.GetStringSlice("gui settings", "valid themes", "NULL", ",")
	if len(themes) != 5 {
		t.Logf("5 Valid Themes %s\n", themes)
		t.Error("Only 5 themes are valid")
	}

    sections := ini.GetSectionValues()
    // "", "gui settins", "theme"
    if len(sections) != 3 {
        t.Logf("Sections... %s\n", sections)
        t.Error("Sections len is not 3")
    }
}

func TestFromFile(t *testing.T) {
	ini := goini.NewIni()
	err := ini.LoadFromFile("./test.ini")
	if err != nil {
		t.Logf("Error loading file: %s", err)
		t.Fatal("")
	}

	sf := ini.GetFloat("gui settings", "scale factor", 2.5)
	t.Logf("Scale factor should be 1.33: %f\n", sf)
	sf = ini.GetFloat("gui settings", "scale factor2", 2.5)
	t.Logf("Scale factor2 should be 2.5: %f\n", sf)
	t.Logf("info text is %s\n", ini.GetString("", "info text", "default value"))
	if ini.GetBool("theme", "use system theme", true) == true {
		t.Error("use system theme must be false")
	}
	if ini.GetInt("theme", "accent color", 1000) != 16_711_935 {
		t.Error("color must be 0xff00ff")
	}
}
