package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const(
	Reset = "\033[0m" 
	Red = "\033[31m"
	Green = "\033[32m" 
	Yellow = "\033[33m" 
	Blue = "\033[34m" 
	Magenta = "\033[35m" 
	Cyan = "\033[36m" 
	Gray = "\033[37m" 
	White = "\033[97m"
	Bold      = "\033[1m"
)
func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}


func GenerateID() string {
    timestamp := time.Now().UnixNano()
    randomBytes := make([]byte, 4)
    rand.Read(randomBytes)
    data := append([]byte(fmt.Sprintf("%d", timestamp)), randomBytes...)
    hash := sha1.Sum(data)
    return hex.EncodeToString(hash[:])[:12]
}

func PrintInitMessage(directory string) {
    width := 60
    dirName := filepath.Base(directory)

    topBottom := fmt.Sprintf("%s╔%s╗%s", Cyan, strings.Repeat("═", width-2), Reset)
    empty := fmt.Sprintf("%s║%s║%s", Cyan, strings.Repeat(" ", width-2), Reset)
    
    mainMsg := "Nora Version Control Initialized!"
    initMsg := fmt.Sprintf("Initialized in: %s", dirName)
    timeMsg := time.Now().Format("2006-01-02 15:04:05")

    mainMsgPadding := (width - 2 - len(mainMsg)) / 2
    initMsgPadding := (width - 2 - len(initMsg)) / 2
    timeMsgPadding := (width - 2 - len(timeMsg)) / 2

    fmt.Println()
    fmt.Println(topBottom)
    fmt.Println(empty)
    fmt.Printf("%s║%s%s%s%s%s║%s\n", 
        Cyan, 
        strings.Repeat(" ", mainMsgPadding),
        Magenta + Bold + mainMsg + Reset,
        strings.Repeat(" ", width-2-mainMsgPadding-len(mainMsg)),
        Reset,
        Cyan,
        Reset)
    fmt.Println(empty)
    fmt.Printf("%s║%s%s%s%s%s║%s\n",
        Cyan,
        strings.Repeat(" ", initMsgPadding),
        Green + initMsg + Reset,
        strings.Repeat(" ", width-2-initMsgPadding-len(initMsg)),
        Reset,
        Cyan,
        Reset)
    fmt.Println(empty)
    fmt.Printf("%s║%s%s%s%s%s║%s\n",
        Cyan,
        strings.Repeat(" ", timeMsgPadding),
        Yellow + timeMsg + Reset,
        strings.Repeat(" ", width-2-timeMsgPadding-len(timeMsg)),
        Reset,
        Cyan,
        Reset)
    fmt.Println(empty)
    fmt.Printf("%s╚%s╝%s\n", Cyan, strings.Repeat("═", width-2), Reset)
    fmt.Println()
}