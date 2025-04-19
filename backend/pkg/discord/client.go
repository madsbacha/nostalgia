package discord

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nostalgia/pkg/oauth2"
	"strings"
)

type Client struct {
	ApiEndpoint   string
	LoginEndpoint string
	ClientId      string
	ClientSecret  string
	Scopes        []string
}

func NewClient(clientId, clientSecret string, scopes []string) Client {
	return Client{
		ApiEndpoint:   "https://discord.com/api",
		LoginEndpoint: "https://discord.com/api/oauth2/authorize",
		ClientId:      clientId,
		ClientSecret:  clientSecret,
		Scopes:        scopes,
	}
}

func TokenFromJson(s []byte) (oauth2.Token, error) {
	token := oauth2.Token{}

	err := json.Unmarshal(s, &token)
	if err != nil {
		return oauth2.Token{}, err
	}

	return token, nil
}

func (ctx Client) VerifyCode(code, redirectUri string) (*oauth2.Token, error) {
	return oauth2.ExchangeCode(
		code,
		ctx.ClientId,
		ctx.ClientSecret,
		redirectUri,
		fmt.Sprintf("%s/oauth2/token", ctx.ApiEndpoint))
}

func (ctx Client) GetUser(token oauth2.Token) (User, error) {
	requestUrl := fmt.Sprintf("%s/users/@me", ctx.ApiEndpoint)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return User{}, err
	}

	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return User{}, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("failed to get user: status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)

	user, err := UserFromJson(body)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (ctx Client) GetLoginUrl(redirectUri string) string {
	scope := strings.Join(ctx.Scopes, "%20")
	callbackUrl := url.PathEscape(redirectUri)
	loginUrl := fmt.Sprintf("%s?"+
		"client_id=%s&"+
		"redirect_uri=%s&"+
		"response_type=code&"+
		"scope=%s", ctx.LoginEndpoint, ctx.ClientId, callbackUrl, scope)
	return loginUrl
}
