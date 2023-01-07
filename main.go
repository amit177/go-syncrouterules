package main

import (
	"errors"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/vishvananda/netlink"
)

const configFile = "config.toml"

type Table struct {
	SearchNextHop string `toml:"search_next_hop"`
	TargetTable   int    `toml:"target_table"`
	RulePriority  int    `toml:"rule_priority"`
}

type Config struct {
	SleepTime string           `toml:"sleep_time"`
	Tables    map[string]Table `toml:"tables"`
}

func parseConfig() (config Config, err error) {
	_, err = os.Stat(configFile)
	if err != nil {
		return
	}

	_, err = toml.DecodeFile(configFile, &config)
	if err != nil {
		return
	}

	_, err = time.ParseDuration(config.SleepTime)
	if err != nil {
		err = errors.New("Could not parse sleep time: " + err.Error())
		return
	}

	for name, table := range config.Tables {
		if net.ParseIP(table.SearchNextHop) == nil {
			err = errors.New("Invalid search_next_hop IP address in table '" + name + "'")
			return
		}
	}

	return
}

func main() {
	config, err := parseConfig()
	if err != nil {
		LogMessage(FATAL, "main.main", "Could not load config file '"+configFile+"': "+err.Error())
	}

	sleepTime, _ := time.ParseDuration(config.SleepTime)

	for {
		for name, table := range config.Tables {
			LogMessage(INFO, "main.main", "Looking for routes matching table '"+name+"'")
			routes := scanRoutes(table.SearchNextHop)
			LogMessage(INFO, "main.main", "Found "+strconv.Itoa(len(routes))+" routes, checking rules")
			rules := scanRules(table.TargetTable)
			LogMessage(INFO, "main.main", "Found "+strconv.Itoa(len(rules))+" rules")
			syncRouteRules(routes, rules, table.TargetTable, table.RulePriority)
		}

		time.Sleep(sleepTime)
	}
}

// scanRoutes returns a map of routes that have a specific next-hop
func scanRoutes(targetNextHop string) map[string]interface{} {
	matchingRoutes := make(map[string]interface{})

	routes, err := netlink.RouteList(nil, netlink.FAMILY_V4)
	if err != nil {
		LogMessage(ERROR, "main.scanRoutes", "Got error when fetching route list: "+err.Error())
		return matchingRoutes
	}

	for _, route := range routes {
		if route.Dst != nil {
			if _, exists := matchingRoutes[route.Dst.String()]; !exists && route.Gw.String() == targetNextHop {
				matchingRoutes[route.Dst.String()] = nil
			}
		}
	}

	return matchingRoutes
}

// scanRules returns a map of rules that are in a specific table
func scanRules(targetTable int) map[string]interface{} {
	matchingRules := make(map[string]interface{})

	rules, err := netlink.RuleList(netlink.FAMILY_V4)
	if err != nil {
		LogMessage(ERROR, "main.scanRules", "Got error when fetching rule list: "+err.Error())
		return matchingRules
	}

	for _, rule := range rules {
		if rule.Src != nil {
			if _, exists := matchingRules[rule.Src.String()]; !exists && rule.Table == targetTable {
				matchingRules[rule.Src.String()] = nil
			}
		}
	}

	return matchingRules
}

// syncRouteRules creates/deletes rules for the matched routes
func syncRouteRules(routes map[string]interface{}, rules map[string]interface{}, targetTable int, rulePriority int) {
	// sync routes -> add routes that have no rules
	for route := range routes {
		if _, exists := rules[route]; !exists {
			LogMessage(INFO, "main.syncRouteRules", "Adding rule for "+route)

			_, ipv4Net, err := net.ParseCIDR(route)
			if err != nil {
				LogMessage(ERROR, "main.syncRouteRules", "Error parsing '"+route+"': "+err.Error())
				continue
			}

			rule := netlink.NewRule()
			rule.Priority = rulePriority
			rule.Family = netlink.FAMILY_V4
			rule.Table = targetTable
			rule.Src = ipv4Net

			err = netlink.RuleAdd(rule)
			if err != nil {
				LogMessage(ERROR, "main.syncRouteRules", "Error adding rule: "+err.Error())
				continue
			}
		}
	}

	// sync rules -> remove rules that have no routes
	for rule := range rules {
		if _, exists := routes[rule]; !exists {
			LogMessage(INFO, "main.syncRouteRules", "Removing rule for "+rule)

			_, ipv4Net, err := net.ParseCIDR(rule)
			if err != nil {
				LogMessage(ERROR, "main.syncRouteRules", "Error parsing '"+rule+"': "+err.Error())
				continue
			}

			rule := netlink.NewRule()
			rule.Family = netlink.FAMILY_V4
			rule.Table = targetTable
			rule.Src = ipv4Net

			err = netlink.RuleDel(rule)
			if err != nil {
				LogMessage(ERROR, "main.syncRouteRules", "Error deleting rule: "+err.Error())
				continue
			}
		}
	}
}
