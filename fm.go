package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/rivo/tview"
)

func collapseAll(root *tview.TreeNode) {
	for _, v := range root.GetChildren() {
		collapseAll(v)
		v.SetExpanded(false)
	}
}

func showInfo(tree *tview.TreeView, n *tview.TreeNode, fmFlex *tview.Flex) {
	f, err := os.Open(n.GetReference().(string))
	if err != nil {
		return
	}
	stat, err := f.Stat()
	if err != nil {
		return
	}
	closeSideWindows(tree, fmFlex)
	if infoNode == n {
		infoNode = nil
		return
	}
	previewNode = nil
	infoNode = n
	sideWindow = tview.NewTextView().SetText(fmt.Sprintf("Name: %s\nSize: %d bytes\nModified: %s",
		stat.Name(), stat.Size(), stat.ModTime().Format(time.RFC822)))

	fmFlex.AddItem(sideWindow, 0, 1, false)
}

func newFileWindow(tree *tview.TreeView, root *tview.TreeNode, fmFlex *tview.Flex) {
	closeSideWindows(tree, fmFlex)

	current := tree.GetCurrentNode()
	ref := *rootDir
	if current != root {
		ref = current.GetReference().(string)
	}
	if !isDir(ref) {
		return
	}
	newFileWin = tview.NewForm().
		AddInputField("Name", "", 30, nil, nil).
		AddTextView("Notice", "Put a / after the name to create a folder", 30, 20, false, false).
		AddButton("Create", func() {
			txt := newFileWin.GetFormItem(0).(*tview.InputField).GetText()
			p := ref + "/" + txt
			if txt[len(txt)-1] == '/' {
				os.Mkdir(p, os.ModePerm)
			} else {
				os.Create(p)
			}
			closeSideWindows(tree, fmFlex)

			current.ClearChildren()
			readDir(current, ref)
			tree.SetCurrentNode(current)
			current.Expand()
		}).
		AddButton("Cancel", func() {
			closeSideWindows(tree, fmFlex)
		})
	fmFlex.AddItem(newFileWin, 0, 1, false)
	app.SetFocus(newFileWin)
}

func previewFile(tree *tview.TreeView, n *tview.TreeNode, fmFlex *tview.Flex, root *tview.TreeNode) {
	if n == root {
		return
	}
	closeSideWindows(tree, fmFlex)
	if n == previewNode {
		previewNode = nil
		return
	}
	previewNode = n
	infoNode = nil
	if isDir(n.GetReference().(string)) {
		return
	}
	f, err := os.Open(n.GetReference().(string))
	if err != nil {
		return
	}
	defer f.Close()
	ext := strings.ReplaceAll(path.Ext(n.GetText()), ".", "")
	if ext == "png" || ext == "jpg" || ext == "jpeg" {
		var (
			err error
			img image.Image
		)
		if ext == "png" {
			img, err = png.Decode(f)
		} else {
			img, err = jpeg.Decode(f)
		}
		if err != nil {
			return
		}
		sideWindow = tview.NewImage().SetImage(img)
	} else {
		content, err := os.ReadFile(n.GetReference().(string))
		if err != nil {
			return
		}
		sideWindow = tview.NewTextView().SetText(string(content))
	}

	fmFlex.AddItem(sideWindow, 0, 1, false)

}

func deleteFile(tree *tview.TreeView, root *tview.TreeNode) {
	p := tree.GetCurrentNode().GetReference().(string)
	os.RemoveAll(p)
	le := strings.Split(p, "/")
	pt := strings.Join(le[:len(le)-1], "/")
	updateTree(tree, root, pt)

}

func moveFile(tree *tview.TreeView, root *tview.TreeNode) {
	if len(mv[0]) == 0 || len(mv[1]) == 0 {
		return
	}
	if !isDir(mv[1]) {
		return
	}

	if err := exec.Command("mv", mv[0], mv[1]).Run(); err != nil {
		return
	}

	le := strings.Split(mv[0], "/")
	pt := strings.Join(le[:len(le)-1], "/")
	updateTree(tree, root, pt)
	updateTree(tree, root, mv[1])
	mv = [2]string{}

}

func copyFile(tree *tview.TreeView, root *tview.TreeNode) {

	if len(cp[0]) == 0 || len(cp[1]) == 0 {
		return
	}
	if !isDir(cp[1]) {
		return
	}

	if err := exec.Command("cp", "-r", cp[0], cp[1]).Run(); err != nil {
		return
	}
	n := root
	if cp[1] != *rootDir {
		n = findNode(root, cp[1])
	}
	if n == nil {
		n = root
	}
	n.ClearChildren()
	readDir(n, cp[1])
	tree.SetCurrentNode(n)
	n.Expand()
	cp = [2]string{}
}

func openFile(cmd, path string) {
	cmd = strings.ReplaceAll(cmd, "PATH", path)
	args := strings.Split(cmd, " ")
	switch len(args) {
	case 0:
		return
	case 1:
		cmd := exec.Command(args[0])
		cmd.Run()
		break
	default:
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Run()
	}
}
