package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

const (
	cfgDirName = "reamcli"
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func GetConfigDir() string {
	if runtime.GOOS == "windows" {
		return filepath.Join(os.Getenv("APPDATA"), cfgDirName)
	}
	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg == "" {
		home, _ := os.UserHomeDir()
		xdg = filepath.Join(home, ".config")
	}
	return filepath.Join(xdg, cfgDirName)
}

func EnsureConfigDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0700); err != nil {
			return fmt.Errorf("failed to create config dir %s: %w", dir, err)
		}
	}
	return nil
}

func GetCallbackPort() int {
	return viper.GetInt("callback_port")
}

func RandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic("crypto/rand failed: " + err.Error())
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}
