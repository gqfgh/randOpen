package main

import (
	"log"
	"syscall"
	"unsafe"
)

//goland:noinspection GoSnakeCaseUsage
const (
	FO_DELETE          uint32 = 0x3
	FOF_ALLOWUNDO      uint32 = 0x40 // 允许撤销，删除文件到回收站
	FOF_NOCONFIRMATION uint32 = 0x10 // 无需确认
)

// SHFILEOPSTRUCT shell 文件操作结构体
type SHFILEOPSTRUCT struct {
	hwnd   uintptr // 窗口的 handle
	wFunc  uint32  // 指定操作类型，如 FO_DELETE 删除
	pFrom  *uint16 // 源文件指针
	pTo    *uint16 // 目标文件指针
	fFlags uint32  // 各种选项，如 FOF_ALLOWUNDO, FOF_NOCONFIRMATION
}

// 定义 SHFileOperation 函数
var shell32 = syscall.NewLazyDLL("shell32.dll")
var procSHFileOperation = shell32.NewProc("SHFileOperationW")

// SHFileOperation 返回 0 代表删除成功
func SHFileOperation(shellFileOperationStructPointer *SHFILEOPSTRUCT) int {
	rc, _, _ := procSHFileOperation.Call(uintptr(unsafe.Pointer(shellFileOperationStructPointer)))
	return int(rc)
}

// SendToTrashBin 删除文件到回收站，适用于windows
func SendToTrashBin(file string) (bool, error) {
	pFrom, err := syscall.UTF16PtrFromString(file)
	if err != nil {
		log.Println("源文件指针 pFrom 转换失败", err)
	}
	shellFileOperationStructPointer := &SHFILEOPSTRUCT{
		hwnd:   0,
		wFunc:  FO_DELETE,
		pFrom:  pFrom, // 注意要以两个 \0 结尾
		pTo:    nil,
		fFlags: FOF_ALLOWUNDO | FOF_NOCONFIRMATION,
	}
	ret := SHFileOperation(shellFileOperationStructPointer)
	if ret != 0 {
		// 删除失败，输出错误信息
		err = syscall.Errno(ret)
		log.Println("删除文件失败：", err)
		return false, err
	}
	return true, nil
}
