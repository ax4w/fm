package main

import (
	"os"

	"github.com/rivo/tview"
)

func isDir(ref string) bool {

	f, err := os.Open(ref)
	if err != nil {
		return false
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}
	return stat.IsDir()

}

func closeSideWindows(fmFlex *tview.Flex) {
	if preview != nil {
		fmFlex.RemoveItem(preview)
		preview = nil
	}
	if newFileWin != nil {
		fmFlex.RemoveItem(newFileWin)
		newFileWin = nil
	}
	app.SetFocus(tree)
}
