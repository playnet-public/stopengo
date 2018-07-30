package stopengo_test

import (
	"bytes"
	"fmt"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/playnet-public/stopengo"
)

func TestRedirectURL(t *testing.T) {
	realm, _ := url.Parse("https://somedomain.sometld")
	returnTo, _ := url.Parse("https://localhost:666?testValueToTransport=Something")

	u, err := stopengo.RedirectURL(realm, returnTo)
	if err != nil {
		t.Fatal(err.Error())
	}

	nUrl, _ := url.Parse(stopengo.OpenIDProvider)
	vals := nUrl.Query()
	vals.Set("openid.claimed_id", stopengo.OpenIDClaimedID)
	vals.Set("openid.identity", stopengo.OpenIDIdentifier)
	vals.Set("openid.mode", stopengo.OpenIDModeCheckIDSetup)
	vals.Set("openid.ns", stopengo.OpenIDNS)
	vals.Set("openid.realm", fmt.Sprintf(
		"%s://%s:%s",
		realm.Scheme,
		realm.Host,
	))
	vals.Set("openid.return_to", returnTo.String())
	nUrl.RawQuery = vals.Encode()

	if u != nUrl.String() {
		t.Fatal("the urls have to be the same!")
	}
}

func TestValidate(t *testing.T) {
	//test will follow ASAP
}

func TestSteamID64(t *testing.T) {
	nUrl, _ := url.Parse(stopengo.OpenIDProvider)
	vals := nUrl.Query()
	vals.Set("openid.claimed_id", "https://steamcommunity.com/openid/id/76561198040411592")
	nUrl.RawQuery = vals.Encode()

	buf := bytes.Buffer{}

	req := httptest.NewRequest("GET", nUrl.String(), &buf)

	steamid64, err := stopengo.SteamID64(req)
	if err != nil {
		t.Fatal(err.Error())
	}

	if steamid64 != "76561198040411592" {
		t.Fatal("unequal steamid64")
	}
}
