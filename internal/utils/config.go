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
	appName = "reamcli"
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func getBaseDir(xdgVar, fallbackSubdir, winEnvVar string) (string, error) {
	if runtime.GOOS == "windows" {
		base := os.Getenv(winEnvVar)
		if base == "" {
			return "", fmt.Errorf("%s non défini", winEnvVar)
		}
		return filepath.Join(base, appName), nil
	}

	base := os.Getenv(xdgVar)
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, fallbackSubdir)
	}
	return filepath.Join(base, appName), nil
}

func GetConfigDir() (string, error) {
	return getBaseDir("XDG_CONFIG_HOME", ".config", "APPDATA")
}

func GetLogPath() (string, error) {
	base, err := getBaseDir("XDG_STATE_HOME", ".local/state", "LOCALAPPDATA")
	if err != nil {
		return "", err
	}
	logDir := filepath.Join(base, "log")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(logDir, "app.log"), nil
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
