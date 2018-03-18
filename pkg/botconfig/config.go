package botconfig

import (
	"godiscordbot/pkg/botresponse"
)

//Config is a wrapper for the autoresposnses and mutedServers objects needed by the bot
type Config struct {
	Ars          []*botresponse.AutoResponse
	MutedServers []string
	Token        string
}

//NewConfig is a server cstor
func NewConfig() *Config {
	sc := new(Config)
	sc.Ars = []*botresponse.AutoResponse{}
	sc.MutedServers = []string{}
	sc.Token = "NDE1NzMxMDQ0MDgwNDE4ODE2.DW6Lpw.LYY0ZKfCKvC3YMJbCe1u0XpIF7M" //dev bot
	//sc.token = "NDA4ODQ1OTY4MDEyNzM4NTYw.DVV_9w.wRs92tvW30aAmW8JOMgqB2GzFQY" //main bot

	return sc
}
