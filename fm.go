package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

func collapseAll(root *tview.TreeNode) {
	for _, v := range root.GetChildren() {
		collapseAll(v)
		v.SetExpanded(false)
	}
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
func showInfo(n *tview.TreeNode, fmFlex *tview.Flex) {
	f, err := os.Open(n.GetReference().(string))
	if err != nil {
		return
	}
	stat, err := f.Stat()
	if err != nil {
		return
	}
	if preview != nil {
		fmFlex.RemoveItem(preview)
		preview = nil
	}
	preview = tview.NewTextView().SetText(fmt.Sprintf("Name: %s\nSize: %d bytes\nModified: %s",
		stat.Name(), stat.Size(), stat.ModTime().Format(time.RFC822)))

	fmFlex.AddItem(preview, 0, 1, false)
}

func previewFile(n *tview.TreeNode, fmFlex *tview.Flex) {
	if preview != nil {
		fmFlex.RemoveItem(preview)
		preview = nil
	}
	f, err := os.Open(n.GetReference().(string))
	defer f.Close()
	if err != nil {
		return
	}
	stat, err := f.Stat()
	if err != nil {
		return
	}
	if stat.IsDir() {
		return
	}
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
		preview = tview.NewImage().SetImage(img)
		//fmFlex.AddItem(tview.NewImage().SetImage(img), 0, 1, false)
	} else {
		content, err := os.ReadFile(n.GetReference().(string))
		if err != nil {
			return
		}
		preview = tview.NewTextView().SetText(string(content))
	}

	fmFlex.AddItem(preview, 0, 1, false)

}

func deleteFile(tree *tview.TreeView, root *tview.TreeNode) {
	p := tree.GetCurrentNode().GetReference().(string)
	os.RemoveAll(p)
	le := strings.Split(p, "/")
	pt := strings.Join(le[:len(le)-1], "/")
	n := findNode(root, pt)
	if n == nil {
		n = root
	}
	n.ClearChildren()
	readDir(n, pt)
	tree.SetCurrentNode(n)
	n.Expand()

}

func moveFile(tree *tview.TreeView, root *tview.TreeNode) {
	if len(mv[0]) == 0 || len(mv[1]) == 0 {
		return
	}
	f, err := os.Open(mv[1])
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
	f.Close()

	cmd := exec.Command("mv", "", mv[0], mv[1])
	cmd.Run()

	n := findNode(root, mv[1])

	if n == nil {
		n = root
	}
	n.ClearChildren()
	readDir(n, mv[1])
	tree.SetCurrentNode(n)
	n.Expand()

	n = findNode(root, mv[0])

	le := strings.Split(mv[0], "/")
	pt := strings.Join(le[:len(le)-1], "/")
	n = findNode(root, pt)

	if n == nil {
		n = root
	}
	n.ClearChildren()
	readDir(n, pt)
	mv[0] = ""
	mv[1] = ""

}

func copyFile(tree *tview.TreeView, root *tview.TreeNode) {
	if len(cp[0]) == 0 || len(cp[1]) == 0 {
		return
	}
	f, err := os.Open(cp[1])
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
	f.Close()
	cmd := exec.Command("cp", "-r", cp[0], cp[1])
	if err := cmd.Run(); err != nil {
		panic(err.Error())
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
	cp[0] = ""
	cp[1] = ""
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
