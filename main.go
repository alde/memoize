package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var cacheDir = strings.TrimRight(userHomeDir(), "/") + "/.memoize_cache"

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
}

func setupCache() {
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0755)
	}
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: memoize [clear | full command to be memoized]")
		os.Exit(1)
	}
	if os.Args[1] == "clear" {
		fmt.Print("clearing cache....")
		if err := os.RemoveAll(cacheDir); err != nil {
			fmt.Println("error\n", err)
			os.Exit(1)
		}
		fmt.Print("done\n")
		os.Exit(0)

	}
	setupCache()

	command := strings.Join(os.Args[1:], " ")
	hasher := sha256.New()
	hasher.Write([]byte(command))
	shaSum := hex.EncodeToString(hasher.Sum(nil))

	if !isCached(shaSum) {
		cache, err := exec.Command(os.Args[1], os.Args[2:]...).Output()
		if err != nil {
			log.Fatal(err)
		}
		if err := writeCache(shaSum, cache); err != nil {
			log.Fatal(err)
		}
	}
	output, err := readCache(shaSum)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
}

func isCached(shaSum string) bool {
	cacheFile := filepath.FromSlash(fmt.Sprintf("%s/%s", cacheDir, shaSum))
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		return false
	}
	return true
}

func writeCache(shaSum string, output []byte) error {
	cacheFile := filepath.FromSlash(fmt.Sprintf("%s/%s", cacheDir, shaSum))
	return ioutil.WriteFile(cacheFile, output, 0644)
}

func readCache(shaSum string) ([]byte, error) {
	cacheFile := filepath.FromSlash(fmt.Sprintf("%s/%s", cacheDir, shaSum))
	return ioutil.ReadFile(cacheFile)
}
