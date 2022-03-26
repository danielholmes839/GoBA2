package auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type discordUser struct {
	ID            string      `json:"id"`
	Email         string      `json:"email"`
	Username      string      `json:"username"`
	Avatar        string      `json:"avatar"`
	Locale        string      `json:"locale"`
	Discriminator string      `json:"discriminator"`
	Token         string      `json:"token"`
	Verified      bool        `json:"verified"`
	MFAEnabled    bool        `json:"mfa_enabled"`
	Banner        string      `json:"banner"`
	AccentColor   int         `json:"accent_color"`
	Bot           bool        `json:"bot"`
	PublicFlags   interface{} `json:"public_flags"`
	PremiumType   int         `json:"premium_type"`
	System        bool        `json:"system"`
	Flags         int         `json:"flags"`
}

func DiscordIdentity(client *http.Client) (*Identity, error) {
	// get discord identity values
	res, err := client.Get("https://discord.com/api/users/@me")
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// unmarshal the discord user
	user := &discordUser{}
	if err := json.Unmarshal(data, user); err != nil {
		return nil, err
	}

	// parse the color
	hex := strconv.FormatInt(int64(user.AccentColor), 16) // 15 -> f
	color := fmt.Sprintf("#%0*s", 6-len(hex), hex)        // f -> #00000f

	identity := &Identity{
		Provider: "discord",
		UserID:   user.ID,
		Username: user.Username,
		AvatarID: user.Avatar,
		Color:    color,
	}

	return identity, nil
}
