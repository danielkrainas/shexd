package mods

type ChannelMap map[string]*Channel

type Channel struct {
	Alias    string `json:"alias"`
	Protocol string `json:"protocol"`
	Endpoint string `json:"endpoint"`
}
