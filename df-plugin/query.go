package plugin

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/gameparrot/dfquery"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"net"
	"strconv"
	"strings"
)

type QueryInfo map[string]string

func (q QueryInfo) HostName() string {
	return q[dfquery.QueryKeyHostName]
}

func (q QueryInfo) SetHostName(s string) {
	q[dfquery.QueryKeyHostName] = s
}

func (q QueryInfo) GameType() string {
	return q[dfquery.QueryKeyGameType]
}

func (q QueryInfo) SetGameType(s string) {
	q[dfquery.QueryKeyGameType] = s
}

func (q QueryInfo) GameID() string {
	return q[dfquery.QueryKeyGameID]
}

func (q QueryInfo) SetGameID(s string) {
	q[dfquery.QueryKeyGameID] = s
}

func (q QueryInfo) Version() string {
	return q[dfquery.QueryKeyVersion]
}

func (q QueryInfo) SetVersion(s string) {
	q[dfquery.QueryKeyVersion] = s
}

func (q QueryInfo) ServerEngine() string {
	return q[dfquery.QueryKeyServerEngine]
}

func (q QueryInfo) SetServerEngine(s string) {
	q[dfquery.QueryKeyServerEngine] = s
}

func (q QueryInfo) Plugins() string {
	return q[dfquery.QueryKeyPlugins]
}

func (q QueryInfo) SetPlugins(s string) {
	q[dfquery.QueryKeyPlugins] = s
}

func (q QueryInfo) Map() string {
	return q[dfquery.QueryKeyMap]
}

func (q QueryInfo) SetMap(s string) {
	q[dfquery.QueryKeyMap] = s
}

func (q QueryInfo) PlayersAmount() int {
	i, err := strconv.Atoi(q[dfquery.QueryKeyNumPlayers])
	if err != nil {
		return -1
	}

	return i
}

func (q QueryInfo) SetPlayersAmount(a int) {
	q[dfquery.QueryKeyNumPlayers] = strconv.Itoa(a)
}

func (q QueryInfo) MaxPlayers() int {
	i, err := strconv.Atoi(q[dfquery.QueryKeyMaxPlayers])
	if err != nil {
		return -1
	}

	return i
}

func (q QueryInfo) SetMaxPlayers(a int) {
	q[dfquery.QueryKeyMaxPlayers] = strconv.Itoa(a)
}

func (q QueryInfo) Whitelist() bool {
	return q[dfquery.QueryKeyWhitelist] == "true"
}

func (q QueryInfo) SetWhitelist(a bool) {
	q[dfquery.QueryKeyWhitelist] = strconv.FormatBool(a)
}

type QueryHandler interface {
	HandleRequestInfo(addr net.Addr, info QueryInfo, players *[]string)
}

type NopQueryHandler struct{}

func (n NopQueryHandler) HandleRequestInfo(net.Addr, QueryInfo, *[]string) {}

func handleQuery(srv *server.Server, m *Manager, h QueryHandler) func(addr net.Addr) (map[string]string, []string) {
	return func(addr net.Addr) (map[string]string, []string) {
		var plugins []string
		for name := range m.Plugins() {
			plugins = append(plugins, name)
		}

		info := map[string]string{
			dfquery.QueryKeyHostName:     m.srvConf.Name,
			dfquery.QueryKeyGameType:     "MINECRAFTPE",
			dfquery.QueryKeyVersion:      "v" + protocol.CurrentVersion,
			dfquery.QueryKeyServerEngine: "Dragonfly",
			dfquery.QueryKeyPlugins:      strings.Join(plugins, ", "),
			dfquery.QueryKeyNumPlayers:   strconv.Itoa(srv.PlayerCount()),
			dfquery.QueryKeyMaxPlayers:   strconv.Itoa(srv.MaxPlayerCount()),
			dfquery.QueryKeyWhitelist:    "false",
		}

		var players []string
		for p := range srv.Players(nil) {
			players = append(players, p.Name())
		}

		if h != nil && addr != nil && info != nil && players != nil {
			h.HandleRequestInfo(addr, info, &players)
		}

		return info, players
	}
}
