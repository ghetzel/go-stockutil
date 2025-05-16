package fileutil

import (
	"os"
	"time"
)

// An os.FileInfo-compatible wrapper that allows for individual values to be overridden.
type FileInfo struct {
	os.FileInfo
	name  string
	size  int64
	mode  *os.FileMode
	mtime time.Time
	dir   *bool
	sys   any
}

func NewFileInfo(wrap ...os.FileInfo) *FileInfo {
	if len(wrap) > 0 && wrap[0] != nil {
		return &FileInfo{
			FileInfo: wrap[0],
		}
	} else {
		return new(FileInfo)
	}
}

func (info *FileInfo) Name() string {
	if info.name != `` {
		return info.name
	} else if info.FileInfo != nil {
		return info.FileInfo.Name()
	} else {
		return ``
	}
}

func (info *FileInfo) Size() int64 {
	if info.size != 0 {
		return info.size
	} else if info.FileInfo != nil {
		return info.FileInfo.Size()
	} else {
		return 0
	}
}

func (info *FileInfo) Mode() os.FileMode {
	if info.mode != nil {
		return *info.mode
	} else if info.FileInfo != nil {
		return info.FileInfo.Mode()
	} else {
		return 0
	}
}

func (info *FileInfo) ModTime() time.Time {
	if !info.mtime.IsZero() {
		return info.mtime
	} else if info.FileInfo != nil {
		return info.FileInfo.ModTime()
	} else {
		return time.Time{}
	}
}

func (info *FileInfo) IsDir() bool {
	if info.dir != nil {
		return *info.dir
	} else if info.FileInfo != nil {
		return info.FileInfo.IsDir()
	} else {
		return false
	}
}

func (info *FileInfo) Sys() any {
	if info.sys != nil {
		return info.sys
	} else if info.FileInfo != nil {
		return info.FileInfo.Sys()
	} else {
		return nil
	}
}

func (info *FileInfo) SetName(name string) {
	info.name = name
}

func (info *FileInfo) SetSize(sz int64) {
	info.size = sz
}

func (info *FileInfo) SetMode(mode os.FileMode) {
	info.mode = &mode
}

func (info *FileInfo) SetModTime(mtime time.Time) {
	info.mtime = mtime
}

func (info *FileInfo) SetIsDir(isDir bool) {
	info.dir = &isDir
}

func (info *FileInfo) SetSys(iface any) {
	info.sys = iface
}
