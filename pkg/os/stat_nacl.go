// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package os

import "syscall"

func isSymlink(stat *syscall.Stat_t) bool {
	return stat.Mode&syscall.S_IFMT == syscall.S_IFLNK
}

func fileInfoFromStat(name string, fi *FileInfo, lstat, stat *syscall.Stat_t) *FileInfo {
	fi.Dev = uint64(stat.Dev)
	fi.Ino = uint64(stat.Ino)
	fi.Nlink = uint64(stat.Nlink)
	fi.Mode = stat.Mode
	fi.Uid = int(stat.Uid)
	fi.Gid = int(stat.Gid)
	fi.Rdev = uint64(stat.Rdev)
	fi.Size = int64(stat.Size)
	fi.Blksize = int64(stat.Blksize)
	fi.Blocks = int64(stat.Blocks)
	fi.Atime_ns = int64(stat.Atime) * 1e9
	fi.Mtime_ns = int64(stat.Mtime) * 1e9
	fi.Ctime_ns = int64(stat.Ctime) * 1e9
	for i := len(name) - 1; i >= 0; i-- {
		if name[i] == '/' {
			name = name[i+1:]
			break
		}
	}
	fi.Name = name
	if isSymlink(lstat) && !isSymlink(stat) {
		fi.FollowedSymlink = true
	}
	return fi
}
