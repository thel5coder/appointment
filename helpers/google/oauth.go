package google

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const OauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func GetGoogleProfile(token string) (res map[string]interface{},err error) {
	response, err := http.Get(OauthGoogleUrlAPI + token)
	if err != nil {
		return res, err
	}

	defer response.Body.Close()
	responseBody, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(responseBody))
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(responseBody,&res)
	if err != nil {
		return res,err
	}
	fmt.Print(res)

	return res,err
}

func getUser(token string) (res []byte,err error) {
	response, err := http.Get(OauthGoogleUrlAPI + token)
	if err != nil {
		fmt.Println("error token")
		return res, err
	}
	defer response.Body.Close()
	res, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("error read body")

		return nil, err
	}

	return res, nil
}