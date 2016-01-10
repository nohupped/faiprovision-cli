package faimodules

import (
//	"fmt"
	"os"
	"syscall"
	"fmt"
	"time"
	"strconv"
	"io"
)

func WriteIncludeToMainConf(includefile string, mainfile string)  {
	file, err := os.OpenFile(mainfile, os.O_APPEND|os.O_RDWR, 0600)
	CheckForError(err)
	GetLock(file, syscall.LOCK_EX)
	defer UngetLock(file)
	includefilepath := fmt.Sprintf("include \"%s\";\n", includefile)
	fmt.Println("Adding the entry", includefilepath, "to", mainfile)
	data, err := file.WriteString(includefilepath)
	CheckForError(err)
	Info.Println(data, "bytes written to ", mainfile)
	fmt.Println(data, "bytes written to", mainfile)

}

//WriteToIncludeConf accepts parameters required to populate the dhcp include file such as
// the conf file path name, struct Host, and the next server's IP address which is usually the tftp server
func WriteToIncludeConf(includefile string, h *Host, nextserver string)  {
	file, err := os.OpenFile(includefile, os.O_EXCL|os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()
	CheckForError(err)
	GetLock(file, syscall.LOCK_EX)
	defer UngetLock(file)
	includeconf := fmt.Sprintf("next-server %s;\nhost %s {\n\toption host-name \"%s\";\n\thardware ethernet %s;\n\tfixed-address %s;\n\tfilename \"pxelinux.0\";\n\toption routers %s;\n}\n", nextserver, h.GetHostname(), h.GetHostname(),h.GetMacID(), h.GetIP(), h.GetRoute())
	file.WriteString(includeconf)
	fmt.Println("Creating include file. Contents: \n\n", includeconf)
	Info.Println("Created IncludeFile. Contents: \n\n", includeconf)
}

func TakeBackup(file string) string{
	var backupfiles string
	fmt.Println("Backing up ", file)
	Info.Println("Taking Backup of ", file)
	src, err := os.Open(file)
	CheckForError(err)
	GetLock(src, syscall.LOCK_EX)
	defer UngetLock(src)
	now := strconv.Itoa(int(time.Now().Unix()))
	tmp := file + now
	dst, err := os.Create(tmp)
	CheckForError(err)
	io.Copy(dst, src)
	backupfiles = dst.Name()
	Info.Println("dhcp conf backed up to ", backupfiles)
	fmt.Println("dhcp.conf backed up to ", backupfiles)
	return backupfiles
}

func CopyFiles(filesource, filedest string) {
	fmt.Println("Restoring file  ", filesource, "->", filedest)
	Info.Println("Restoring file  ", filesource, "->", filedest)
	src, err := os.Open(filesource)
	CheckForError(err)
	GetLock(src, syscall.LOCK_EX)
	defer UngetLock(src)
	dst, err := os.Create(filedest)
	CheckForError(err)
	io.Copy(dst, src)
	Info.Println(src.Name(), "->", dst.Name())
	fmt.Println(src.Name(), "->", dst.Name())
}
