package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	Token_File_Name = "rulepolicym_token.toml"
	Token_File_Path = "."
)

var now = time.Now()

func init() {
	filename := strings.Split(Token_File_Name, ".")[0]
	extension := strings.Split(Token_File_Name, ".")[1]
	viper.SetConfigName(filename)
	viper.SetConfigType(extension)
	viper.AddConfigPath(Token_File_Path)
}

func ReadToken() (string, string, error) {
	// use viper to read token from file

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; write the token file
			err = CreateTokenFile()
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", fmt.Errorf("not config file not found err:%w", err)
		}

	}
	tokenfileinfo, err := os.Stat(filepath.Join(Token_File_Path, Token_File_Name))
	if err != nil {
		return "", "", err
	}
	// if the tokenfile's modification is 14 minutes ago, write new token file
	if now.Sub(tokenfileinfo.ModTime()) > 14*time.Minute {
		err = UpdateToken()
		if err != nil {
			return "", "", err
		}
	}
	// get token from viper
	accesstoken := viper.GetString("accesstoken")
	authtoken := viper.GetString("authtoken")
	if accesstoken == "" || authtoken == "" {
		return "", "", fmt.Errorf("token not found in config file")
	}
	return accesstoken, authtoken, nil

}

func UpdateToken() error {
	accesstokenresp, err := GetAccessToken()
	if err != nil {
		return err
	}
	authtokenresp, err := GetAuthToken(accesstokenresp.Access_Token)
	if err != nil {
		return err
	}
	// write token to file
	viper.Set("accesstoken", accesstokenresp.Access_Token)
	viper.Set("authtoken", authtokenresp.Token)
	return viper.WriteConfig()
}

func CreateTokenFile() error {
	file, err := os.Create(filepath.Join(Token_File_Path, Token_File_Name))
	if err != nil {
		return err
	}
	defer file.Close()
	return UpdateToken()
}
