package main

import (
	"fmt"
	"github.com/easysoft/zentaoatf/src/action"
	httpClient "github.com/easysoft/zentaoatf/src/http"
	"github.com/easysoft/zentaoatf/src/mock"
	"github.com/easysoft/zentaoatf/src/utils"
	"github.com/jroimartin/gocui"
	"log"
	"net/http/httptest"
	"strconv"
	"strings"
	"time"
)

const (
	leftWidth          = 32
	labelWidth         = 15
	inputFullLineWidth = 64
	inputNumbWidth     = 25
	buttonWidth        = 10
	space              = 2
)

var server *httptest.Server

func main() {
	server = mock.Server("case-from-prodoct.json")
	defer server.Close()

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true
	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("qickbar", 0, 0, leftWidth, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if v, err := g.SetView("import", 3, 0, 14, 2); err != nil {
		v.Frame = false
		fmt.Fprint(v, "  Import   ")
	}
	if v, err := g.SetView("switch", 19, 0, 31, 2); err != nil {
		v.Frame = false
		fmt.Fprint(v, "  Switch   ")
	}

	if v, err := g.SetView("side", 0, 2, leftWidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}

	if v, err := g.SetView("main", leftWidth, 0, maxX-1, maxY-10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
	}

	if v, err := g.SetView("cmd", leftWidth, maxY-10, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true
		v.Autoscroll = true
	}
	return nil

}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("import", gocui.MouseLeft, gocui.ModNone, importProjectUi); err != nil {
		return err
	}
	if err := g.SetKeybinding("switch", gocui.MouseLeft, gocui.ModNone, switchProjectUi); err != nil {
		return err
	}

	if err := g.SetKeybinding("cmdline", gocui.MouseLeft, gocui.ModNone, setEdit); err != nil {
		return err
	}

	return nil
}

func setEdit(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("cmdline"); err != nil {
		return err
	}

	v.Autoscroll = true
	v.Clear()

	return nil
}

func importProjectUi(g *gocui.Gui, v *gocui.View) error {
	maxX, _ := g.Size()

	slideView, _ := g.View("side")
	slideX, _ := slideView.Size()

	left := slideX + 2
	right := left + labelWidth
	if v, err := g.SetView("urlLabel", left, 1, right, 3); err != nil {
		v.Frame = false
		fmt.Fprint(v, "ZentaoUrl")
	}

	left = right + space
	right = left + inputFullLineWidth
	if v, err := g.SetView("urlInput", left, 1, right, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Clear()
		fmt.Fprint(v, server.URL)

		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("urlInput"); err != nil {
			return err
		}
	}

	left = slideX + 2
	right = left + labelWidth
	if v, err := g.SetView("productLabel", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprint(v, "ProdoctId")
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("productInput", left, 4, right, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprint(v, "1")
	}

	left = right + space
	right = left + labelWidth
	if v, err := g.SetView("taskLabel", left, 4, right, 6); err != nil {
		v.Frame = false
		fmt.Fprint(v, "or TaskId")
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("taskInput", left, 4, right, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprint(v, "1")
	}

	left = slideX + 2
	right = left + labelWidth
	if v, err := g.SetView("languageLabel", left, 7, right, 9); err != nil {
		v.Frame = false
		fmt.Fprint(v, "Language")
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("languageInput", left, 7, right, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Editable = true
		v.Wrap = true

		fmt.Fprint(v, "python")
	}

	left = right + space
	right = left + labelWidth
	if v, err := g.SetView("singleFileLabel", left, 7, right, 9); err != nil {
		v.Frame = false
		fmt.Fprint(v, "SingleFile")
	}

	left = right + space
	right = left + inputNumbWidth
	if v, err := g.SetView("singleFileInput", left, 7, right, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		if err := g.SetKeybinding("singleFileInput", gocui.MouseLeft, gocui.ModNone, changeSingleFile); err != nil {
			return err
		}

		fmt.Fprint(v, "true")
	}

	buttonX := (maxX-leftWidth)/2 + leftWidth - buttonWidth
	if v, err := g.SetView("submit", buttonX, 10, buttonX+buttonWidth, 12); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Highlight = true
		v.BgColor = gocui.ColorGreen
		v.FgColor = gocui.ColorBlack

		fmt.Fprint(v, "  Submit  ")
		if err := g.SetKeybinding("submit", gocui.MouseLeft, gocui.ModNone, importProjectRequest); err != nil {
			return err
		}
	}

	return nil
}

func switchProjectUi(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func importProjectRequest(g *gocui.Gui, v *gocui.View) error {
	urlView, _ := g.View("urlInput")
	productView, _ := g.View("productInput")
	taskView, _ := g.View("taskInput")
	languageView, _ := g.View("languageInput")
	singleFileView, _ := g.View("singleFileInput")

	url := strings.TrimSpace(urlView.ViewBuffer())

	productCode := strings.TrimSpace(productView.Buffer())
	taskId := strings.TrimSpace(taskView.Buffer())
	language := strings.TrimSpace(languageView.Buffer())
	singleFileStr := strings.TrimSpace(singleFileView.Buffer())
	singleFile, e := strconv.ParseBool(singleFileStr)
	if e != nil {
		singleFile = true
	}

	params := make(map[string]string)
	if productCode != "" {
		params["entityType"] = "product"
		params["entityVal"] = productCode
	} else {
		params["entityType"] = "task"
		params["entityVal"] = taskId
	}

	cmdView, _ := g.View("cmd")
	_, _ = fmt.Fprintln(cmdView, fmt.Sprintf("#atf gen -u %s -t %s -v %s -l %s -s %t",
		url, params["entityType"], params["entityVal"], language, singleFile))

	jsonBuf, e := httpClient.GetBuf(url, params)
	if e != nil {
		fmt.Fprint(cmdView, e.Error())
		return nil
	}

	err := action.Generate(jsonBuf, language, singleFile)
	if err == nil {

		fmt.Fprintln(cmdView, fmt.Sprintf("success to generate test scripts in 'xdoc/scripts' at %s",
			utils.DateTimeStr(time.Now())))
	} else {
		fmt.Fprint(cmdView, err.Error())
	}

	return nil
}

func changeSingleFile(g *gocui.Gui, v *gocui.View) error {
	val := strings.TrimSpace(v.Buffer())

	singleFile, e := strconv.ParseBool(val)

	v.Clear()
	if e != nil || !singleFile {
		fmt.Fprint(v, "true")
	} else {
		fmt.Fprint(v, "false")
	}

	return nil
}
