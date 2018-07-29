package stopengo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const (
	//OpenIDProvider is the provider url to authenticate at
	OpenIDProvider = "https://steamcommunity.com/openid/login"
	//OpenIDNS is the providers spec namespace
	OpenIDNS = "http://specs.openid.net/auth/2.0"
	//OpenIDIdentifier is the used identifier by the provider and client
	OpenIDIdentifier = "http://specs.openid.net/auth/2.0/identifier_select"
	//OpenIDClaimedID is the used identifier by the provider and client (in case of Steam's OpenID)
	OpenIDClaimedID = "http://specs.openid.net/auth/2.0/identifier_select"
	//OpenIDModeCheckIDSetup is the mode to redirect to Steam
	OpenIDModeCheckIDSetup = "checkid_setup"
	//OpenIDModeCheckAuthentication checks if the redirect from Steam to the given redirectURL was successful
	OpenIDModeCheckAuthentication = "check_authentication"
)

var (
	//Just FYI regexp's were taken from https://github.com/solovev/steam_go/blob/master/auth.go#L19
	validationRegexp = regexp.MustCompile(`^(https|http):\/\/steamcommunity.com\/openid\/id\/[0-9]{15,25}$`)
	extractionRegexp = regexp.MustCompile("\\D+")

	errInvalidOpenIDNS       = errors.New("invalid open id ns")
	errInvalidOpenIDRequest  = errors.New("invalid open id request")
	errInvalidSteamIDPattern = errors.New("invalid steam id pattern")
)

//RedirectURL creates a new Steam OpenID authentication URL
func RedirectURL(realm, returnTo *url.URL) (string, error) {
	u, err := url.Parse(OpenIDProvider)
	if err != nil {
		return "", err
	}

	vals := u.Query()
	vals.Set("openid.claimed_id", OpenIDClaimedID)
	vals.Set("openid.identity", OpenIDIdentifier)
	vals.Set("openid.mode", OpenIDModeCheckIDSetup)
	vals.Set("openid.ns", OpenIDNS)
	vals.Set("openid.realm", fmt.Sprintf(
		"%s://%s:%s",
		realm.Scheme,
		realm.Host,
		realm.Port(),
	))
	vals.Set("openid.return_to", returnTo.String())

	u.RawQuery = vals.Encode()

	return u.String(), nil
}

//Validate checks an incoming request from Steam for validity
func Validate(r *http.Request) error {
	u, err := url.Parse(OpenIDProvider)
	if err != nil {
		return err
	}

	vals := u.Query()
	vals.Set("openid.mode", OpenIDModeCheckAuthentication)
	vals.Set("openid.assoc_handle", r.FormValue("openid.assoc_handle"))
	vals.Set("openid.signed", r.FormValue("openid.signed"))
	vals.Set("openid.sig", r.FormValue("openid.sig"))
	vals.Set("openid.ns", r.FormValue("openid.ns"))

	for _, v := range strings.Split(r.FormValue("openid.signed"), ",") {
		vals.Set("openid."+v, r.FormValue("openid."+v))
	}

	response, err := http.PostForm(u.String(), vals)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")

	if lines[0] != "ns:"+OpenIDNS {
		return errInvalidOpenIDNS
	}

	if strings.HasSuffix(lines[1], "false") {
		return errInvalidOpenIDRequest
	}

	return nil
}

//SteamID64 returns the steamID64 from an incoming Steam request
func SteamID64(r *http.Request) (string, error) {
	steamIDURL := r.FormValue("openid.claimed_id")
	if !validationRegexp.MatchString(steamIDURL) {
		return "", errInvalidSteamIDPattern
	}

	return extractionRegexp.ReplaceAllString(steamIDURL, ""), nil
}
