package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/togglr-io/togglr"
	"go.uber.org/zap"
)

const tokenURL = "https://github.com/login/oauth/access_token"
const userURL = "https://api.github.com/user"
const emailURL = "https://api.github.com/user/emails"

type oauthResponse struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type userResponse struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type emailResponse struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}

func auth(code, clientID, secret string) (oauthResponse, error) {
	var oauthRes oauthResponse
	client := http.DefaultClient

	req, err := http.NewRequest(http.MethodPost, tokenURL, nil)
	if err != nil {
		return oauthRes, err
	}

	query := req.URL.Query()
	query.Add("client_id", clientID)
	query.Add("client_secret", secret)
	query.Add("code", code)
	req.URL.RawQuery = query.Encode()

	req.Header.Set("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return oauthRes, err
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return oauthRes, err
	}

	if err := json.Unmarshal(body, &oauthRes); err != nil {
		return oauthRes, err
	}

	return oauthRes, nil
}

func getUser(accessToken string) (userResponse, error) {
	var user userResponse
	req, err := http.NewRequest(http.MethodGet, userURL, nil)
	if err != nil {
		return user, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return user, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return user, err
	}
	defer res.Body.Close()

	log.Printf("user response: %s", string(body))
	if err := json.Unmarshal(body, &user); err != nil {
		return user, err
	}

	return user, nil
}

func getPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, emailURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var emails []emailResponse
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.Primary {
			return email.Email, nil
		}
	}

	return "", errors.New("no primary email found")
}

func HandleGithubRedirect(log *zap.Logger, us togglr.UserService) http.HandlerFunc {
	log = log.With(zap.String("handler", "HandleGithubRedirect"))
	clientID := os.Getenv("TOGGLR_GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("TOGGLR_GITHUB_CLIENT_SECRET")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		// first, use authcode with client creds to get an accessToken
		res, err := auth(code, clientID, clientSecret)
		if err != nil {
			log.Error("failed to authenticate", zap.Error(err))
			forbidden(w, "could not authenticate")
			return
		}

		// then, use the accessToken to pull the GH user info
		ghUser, err := getUser(res.AccessToken)
		if err != nil {
			log.Error("failed to fetch user", zap.Error(err))
			forbidden(w, "could not authenticate")
			return
		}

		// we need an email for the user, so we'll pull their primary
		// using the accessToken
		email, err := getPrimaryEmail(res.AccessToken)
		if err != nil {
			log.Error("failed to fetch emails", zap.Error(err))
			forbidden(w, "could not authenticate")
			return
		}

		// check and see if this user already exists in the system
		user, err := us.FetchUserByIdentity(r.Context(), fmt.Sprint(ghUser.ID), togglr.IdentityTypeGithub)
		if err != nil {
			log.Error("failed to fetch existing user")
			forbidden(w, "could not authenticate")
			return
		}

		// if they don't exist, go ahead and create a new one using
		// the GH user data
		if user.ID.IsNull() {
			newUser := togglr.User{
				Name:         ghUser.Name,
				Identity:     fmt.Sprint(ghUser.ID),
				IdentityType: togglr.IdentityTypeGithub,
				Email:        email,
			}

			if _, err := us.CreateUser(r.Context(), newUser); err != nil {
				log.Error("failed to create new user")
				serverError(w, "could not create new user")
				return
			}
		}

		// generate a JWT for the Togglr user

		// add JWT to an HttpOnly cookie

		// respond with noContent
		noContent(w)
	})
}
