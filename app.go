package main

import (
	"context"
	"fmt"
	"github.com/arduino/go-paths-helper"
	"github.com/electricbubble/go-toast"
	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"math/rand"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type App struct {
	ctx         context.Context
	files       paths.PathList
	currentFile *paths.Path
}

// startup 初始化函数
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx                               // 从这里开始，拿到 context
	runtime.WindowSetPosition(a.ctx, 10, 700) // 设置窗口位置
	a.listenKeyboard()
}

// SelectDir 选择文件夹
func (a *App) SelectDir() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择文件夹",
	})
	if err != nil {
		log.Println("选择文件夹出错：", err)
	}
	return dir
}

// TraversalFiles 遍历文件
func (a *App) TraversalFiles(dir string) {
	p := paths.New(dir)
	files, err := p.ReadDirRecursiveFiltered(nil, func(file *paths.Path) bool {
		if file.Ext() == ".txt" {
			return false
		}
		return true
	})
	if err != nil {
		log.Println("遍历文件失败：", err)
		return
	}
	a.files = files
	toast.Push("遍历完成")
}

// OpenFile 打开文件
func (a *App) OpenFile() {
	if a.files == nil {
		toast.Push("还没遍历文件呢！")
		return
	}
	r := rand.New(rand.NewSource(time.Now().Unix()))
	a.currentFile = a.files[r.Intn(a.files.Len())]
	openViaDefaultProc(a.currentFile.String())
}

// DelFile 删除文件
func (a *App) DelFile() {
	closeVlc()
	ok, err := SendToTrashBin(a.currentFile.String())
	if !ok {
		log.Println(fmt.Sprintf("删除文件%q出错：", a.currentFile.String()), err)
		return
	}
}

// RenameFile 重命名文件
func (a *App) RenameFile() {
	closeVlc()
	ext := a.currentFile.Ext()
	newName, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "重命名文件",
		DefaultDirectory: a.currentFile.Parent().String(),
		DefaultFilename:  a.currentFile.String(),
		Filters: []runtime.FileFilter{{
			DisplayName: "*" + ext,
			Pattern:     "*" + ext,
		}},
	})
	if newName == "" {
		toast.Push("取消重命名")
		return
	}
	if err != nil {
		log.Println("重命名时，获取新文件名出错：", err)
		return
	}
	// 保证文件后缀
	if !strings.Contains(newName, ext) {
		newName += ext
	}
	err = a.currentFile.Rename(paths.New(newName))
	if err != nil {
		log.Println("重命名出错：", err)
	}
}

// CopyFileName 复制文件名
func (a *App) CopyFileName() {
	err := runtime.ClipboardSetText(a.ctx, a.currentFile.String())
	if err != nil {
		log.Println("复制文件名到剪切板出错：", err)
	}
}

// 通过默认程序打开文件
func openViaDefaultProc(file string) {
	cmd := exec.Command("powershell", fmt.Sprintf("start-process %q", file))
	cmdSetToHidden(cmd)
	err := cmd.Run()
	if err != nil {
		log.Println("打开文件失败：", err)
	}
}

func cmdSetToHidden(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}

// 监听键盘
func (a *App) listenKeyboard() {
	hook.Register(hook.KeyDown, []string{"/"}, func(event hook.Event) {
		a.OpenFile()
	})
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "d"}, func(event hook.Event) {
		a.DelFile()
	})
	hook.Register(hook.KeyDown, []string{"ctrl", "alt", "r"}, func(event hook.Event) {
		a.RenameFile()
	})
	eventCh := hook.Start()
	<-hook.Process(eventCh)
}

func closeVlc() {
	cmd := exec.Command("powershell", "stop-process -name vlc")
	cmdSetToHidden(cmd)
	err := cmd.Run()
	if err != nil {
		log.Println("关闭 vlc 失败：", err)
	}
	time.Sleep(time.Millisecond * 500)
}
