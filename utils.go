package main

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func updateTree(tree *tview.TreeView, root *tview.TreeNode, p string) {
	n := findNode(root, p)
	if n == nil {
		n = root
	}
	n.ClearChildren()
	readDir(n, p)
	tree.SetCurrentNode(n)
	n.Expand()

}

func findNode(t *tview.TreeNode, path string) *tview.TreeNode {
	for _, v := range t.GetChildren() {
		if v.GetReference().(string) == path {
			return v
		}
		if f := findNode(v, path); f != nil {
			return f
		}
	}
	return nil
}

func readDir(target *tview.TreeNode, path string) {
	files, err := os.ReadDir(path)
	if err != nil {
		return
	}
	for _, file := range files {
		node := tview.NewTreeNode(file.Name()).
			SetReference(filepath.Join(path, file.Name()))
		if file.IsDir() {
			node.SetColor(tcell.ColorGreen)
			node.SetText("🗁 " + file.Name())
		}
		target.AddChild(node)
	}

}

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

func closeSideWindows(tree *tview.TreeView, fmFlex *tview.Flex) {
	if sideWindow != nil {
		fmFlex.RemoveItem(sideWindow)
		sideWindow = nil
	}
	if newFileWin != nil {
		fmFlex.RemoveItem(newFileWin)
		newFileWin = nil
	}
	app.SetFocus(tree)
}
