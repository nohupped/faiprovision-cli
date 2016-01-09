package faimodules

import (
	"os"
	"fmt"
	"syscall"
	"bufio"
	"regexp"
)
// CheckForError check will check for errors and will panic
func CheckForError(e error) {
	if e != nil {
		Error.Println(e)
		panic(e)
	}
}

// GetLock accepts file of type *os.File returned by os.Open,
// and the locktype of type int, which is a constant defined in syscall.LOCK*
func GetLock(file *os.File, locktype int )  {
	fmt.Println("Acquiring exclusive lock on ", file)
	Info.Println("Acquiring lock on", file)
	syscall.Flock(int(file.Fd()), locktype)
	Info.Println("Acquired exclusive lock on ", file)
	fmt.Println("Acquired filelock")
}

// UngetLock accepts file of type *os.File, and removes the lock set by syscall.Flock.
// Call this to avoid race conditions.
func UngetLock(file *os.File)  {
	syscall.Flock(int(file.Fd()), syscall.LOCK_UN);
}

// ReadDhcpRO function Input path to main dhcp conf and returns all include files.
func ReadDhcpRO(conf string) (includeconf []string)  {
	includefiles := make([]string, 0, 5)

	file, err := os.Open(conf)
	CheckForError(err)
//	GetLock(file, syscall.LOCK_EX)


	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^include\ "(.*)";`)
	for scanner.Scan() {
		if re.MatchString(scanner.Text()) {
			includefiles = append(includefiles, re.FindStringSubmatch(scanner.Text())[1])
		}else {
			continue
		}
	}

//	UngetLock(file)
	//syscall.Flock(int(file.Fd()), syscall.LOCK_UN);
	return includefiles
}

// getIpFromInclude will iterate through the slice and return the IPs.
func getIpFromInclude(includefile []string)  {

}



