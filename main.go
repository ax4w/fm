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
	sideWindow  tview.Primitive
	previewNode *tview.TreeNode
	infoNode    *tview.TreeNode
	app         *tview.Application
	newFileWin  *tview.Form
	cp          = [2]string{"", ""}
	mv          = [2]string{"", ""}
	rootDir     *string
)

func main() {
	config := loadConfig()

	lastKey := ' '
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	rootDir = flag.String("dir", home, "Open a specific directory")
	helpStr := ""
	for k, v := range config.KeyBinds {
		helpStr += k + ": " + v + " "
	}
	header := tview.NewTextArea().SetText("Fm - File Manager\n"+helpStr, false).
		SetTextStyle(tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorWhite))

	appFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	appFlex.AddItem(header, 0, 1, false)

	fmFlex := tview.NewFlex().SetDirection(tview.FlexColumn)
	appFlex.AddItem(fmFlex, 0, 15, false)

	app = tview.NewApplication()
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
			action := config.KeyBinds[string(event.Rune())]
			switch strings.ToLower(action) {
			case "new":
				newFileWindow(tree, root, fmFlex)
				break
			case "collapseall":
				collapseAll(root)
				break
			case "delete":
				if lastKey == 'd' {
					deleteFile(tree, root)
					lastKey = ' '
					return event
				}
				break
			case "open":
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
			case "preview":
				n := tree.GetCurrentNode()
				previewFile(tree, n, fmFlex, root)
				break
			case "info":

				showInfo(tree, tree.GetCurrentNode(), fmFlex)
				break
			case "copy":

				if tree.GetCurrentNode().GetReference() != nil {
					cp[0] = tree.GetCurrentNode().GetReference().(string)
				}
				break
			case "paste":
				if tree.GetCurrentNode().GetReference() != nil {
					cp[1] = tree.GetCurrentNode().GetReference().(string)

				} else {
					cp[1] = *rootDir
				}
				copyFile(tree, root)
				break
			case "move":
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

	tree.SetSelectedFunc(func(node *tview.TreeNode) {

		reference := node.GetReference()
		if reference == nil {
			return
		}

		if !isDir(reference.(string)) {
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
	if err := app.SetRoot(appFlex, true).EnableMouse(false).SetFocus(tree).Run(); err != nil {
		panic(err)
	}
}
