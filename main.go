// Demo code for the TreeView primitive.
package main

import (
	"flag"
	"os"
	"path"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	preview       tview.Primitive
	cp            = []string{"", ""}
	mv            = []string{"", ""}
	switchedFocus = false
	rootDir       *string
)

func main() {
	loadConfig()
	lastKey := ' '
	home, _ := os.UserHomeDir()
	rootDir = flag.String("dir", home, "Open a specific directory")

	header := tview.NewTextArea().SetText("Fm - File Manager\n"+
		"CTRL+C: quit, d+d: delete file, ENTER: Select, s: preview file, o: open file, i: file info "+
		"c: copy, p: paste, t: collapse all, m: select | drop", false).
		SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))

	appFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	appFlex.AddItem(header, 0, 1, false)

	fmFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	appFlex.AddItem(fmFlex, 0, 15, false)

	app := tview.NewApplication()
	flag.Parse()
	root := tview.NewTreeNode(*rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)
	readDir(root, *rootDir)

	tree.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); key {
		case tcell.KeyRune:
			switch event.Rune() {
			case 't':
				collapseAll(root)
				break
			case 'd':
				if lastKey == 'd' {
//FixMe: make root not deletable
					deleteFile(tree, root)
					lastKey = ' '
					return event
				}
				break
			case 'o':
				lastKey = 'o'
				ext := strings.ReplaceAll(path.Ext(tree.GetCurrentNode().GetText()), ".", "")
				fpath := tree.GetCurrentNode().GetReference().(string)
				if v, ok := config.OpenInApp[ext]; ok {
					openFile(v, fpath)
				} else {
					if v, ok := config.OpenInApp["default"]; ok {
						openFile(v, fpath)
					}
				}
				break
			case 's':

				n := tree.GetCurrentNode()
				previewFile(n, fmFlex)
				break
			case 'i':

				showInfo(tree.GetCurrentNode(), fmFlex)
				break
			case 'c':

				if tree.GetCurrentNode().GetReference() != nil {
					cp[0] = tree.GetCurrentNode().GetReference().(string)
				}
				break
			case 'p':
				if tree.GetCurrentNode().GetReference() != nil {
					cp[1] = tree.GetCurrentNode().GetReference().(string)

				} else {
					cp[1] = *rootDir
				}
				copyFile(tree, root)
				break
			case 'm':
				if mv[0] == "" {
					if tree.GetCurrentNode().GetReference() == nil {
						return event
					}

					mv[0] = tree.GetCurrentNode().GetReference().(string)
				} else {
					if tree.GetCurrentNode().GetReference() == nil {
						mv[1] = *rootDir
					} else {
						mv[1] = tree.GetCurrentNode().GetReference().(string)
					}
					moveFile(tree, root)
				}
				break
			}
			lastKey = event.Rune()
		}
		return event
	})
	//tree.SetChangedFunc(func(n *tview.TreeNode) { app.Draw() })
	tree.SetSelectedFunc(func(node *tview.TreeNode) {

		reference := node.GetReference()
		if reference == nil {
			return
		}

		f, err := os.Open(reference.(string))
		if err != nil {
			return
		}
		stat, err := f.Stat()
		if err != nil {
			return
		}
		if !stat.IsDir() {
			return
		}

		children := node.GetChildren()
		if len(children) == 0 {
			path := reference.(string)
			readDir(node, path)
			node.SetExpanded(true)
		} else {
			node.SetExpanded(!node.IsExpanded())
			if !node.IsExpanded() {
				node.ClearChildren()
			}
		}
	})

	fmFlex.AddItem(tree, 0, 2, false)
	//fmFlex.AddItem(tview.NewBox().SetTitle("Preview").SetBorder(true), 0, 1, false)

	//	flex.AddItem(tview.NewBox().SetTitle("Preview").SetBorder(true), 0, 1, false)
	if err := app.SetRoot(appFlex, true).EnableMouse(false).SetFocus(tree).Run(); err != nil {
		panic(err)
	}
}
