package pkg

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// Token_File_Name = "rulepolicym_token.toml"
	Token_File_Path = "."
)

var now = time.Now()

var tokenfile_name string

func init() {

	dceurl := os.Getenv(Dce_Url_Env_Name)
	tokenfile_name = fmt.Sprintf("rulepolicym_token"+"_%s.toml", strings.TrimSuffix(strings.TrimPrefix(dceurl, "http://"), "/"))
	// split the tokenfile_name by the last dot
	lastIndex := strings.LastIndex(tokenfile_name, ".")
	var filename, extension string
	if lastIndex == -1 {
		filename = tokenfile_name
		extension = "toml"
	} else {
		filename = tokenfile_name[:lastIndex]
		extension = tokenfile_name[lastIndex+1:]

	}

	//filename := strings.Split(tokenfile_name, ".")[0]
	//extension := strings.Split(tokenfile_name, ".")[1]
	viper.SetConfigName(filename)
	viper.SetConfigType(extension)
	viper.AddConfigPath(Token_File_Path)
}

func ReadToken(client *http.Client) (string, string, error) {
	// use viper to read token from file

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; write the token file
			err = CreateTokenFile(client)
			if err != nil {
				return "", "", err
			}
		} else {
			return "", "", fmt.Errorf("not config file not found err:%w", err)
		}

	}
	tokenfileinfo, err := os.Stat(filepath.Join(Token_File_Path, tokenfile_name))
	if err != nil {
		return "", "", err
	}
	// if the tokenfile's modification is 14 minutes ago, write new token file
	if now.Sub(tokenfileinfo.ModTime()) > 14*time.Minute {
		err = UpdateToken(client)
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

func UpdateToken(client *http.Client) error {
	accesstokenresp, err := GetAccessToken(client)
	if err != nil {
		return err
	}
	authtokenresp, err := GetAuthToken(client, accesstokenresp.Access_Token)
	if err != nil {
		return err
	}
	// write token to file
	viper.Set("accesstoken", accesstokenresp.Access_Token)
	viper.Set("authtoken", authtokenresp.Token)
	return viper.WriteConfig()
}

func CreateTokenFile(client *http.Client) error {
	file, err := os.Create(filepath.Join(Token_File_Path, tokenfile_name))
	if err != nil {
		return err
	}
	defer file.Close()
	return UpdateToken(client)
}
