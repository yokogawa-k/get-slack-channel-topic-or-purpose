package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	flags "github.com/jessevdk/go-flags"
	"github.com/k0kubun/pp"
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

func findPublicChannelIDByName(c *slack.Client, name string) (string, error) {

	convParams := &slack.GetConversationsParameters{
		Limit: 300,
	}

	channels, _, err := c.GetConversations(convParams)
	if err != nil {
		return "", err
	}

	if len(channels) == 0 {
		return "", fmt.Errorf("No channel found")
	}

	for _, c := range channels {
		if c.Name == name {
			log.Debugf("Channel Name [%s] is %s\n", name, c.ID)
			log.Debug(pp.Sprint(c))
			return c.ID, nil
		}
	}

	return "", fmt.Errorf("%s 's channelID not found", name)
}

func findTeamName(c *slack.Client) (string, error) {
	teamInfo, err := c.GetTeamInfo()
	log.Debug(pp.Sprint(teamInfo))
	teamName := teamInfo.Domain
	return teamName, err
}

// discard nanoseconds
func slackTimestampToTime(slackTimestamp string) (time.Time, error) {
	ts := strings.Split(slackTimestamp, ".")
	ti, err := strconv.ParseInt(ts[0], 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.Unix(ti, 0), err
}

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {
	var opts struct {
		Debug        bool   `short:"d" long:"debug" description:"Show debug infomation"`
		SlackToken   string `long:"slack_token" description:"Set Slack API Token" required:"true" env:"SLACK_TOKEN"`
		SlackChannel string `long:"slack_channel" description:"Set target Slack channel name" required:"true" env:"SLACK_CHANNEL"`
		Subtype      string `long:"subtype" description:"Set target Slack message subtype(channel_topic, channel_purpose, channel_join etc...)" required:"true" default:"channel_topic" env:"SLACK_MSG_SUBTYPE"`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Debug {
		log.Info("Set Debug Level")
		log.SetLevel(log.DebugLevel)
	}

	log.Infof("Loglevel is %s", log.GetLevel())
	c := slack.New(opts.SlackToken)

	channelID, err := findPublicChannelIDByName(c, opts.SlackChannel)
	if err != nil {
		panic(err)
	}

	teamName, err := findTeamName(c)
	if err != nil {
		panic(err)
	}

	count := 0
	for cursor := ""; ; {
		histParams := &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Limit:     200,
			Cursor:    cursor,
		}

		r, err := c.GetConversationHistory(histParams)
		if err != nil {
			panic(err)
		}

		for _, msg := range r.Messages {
			if msg.Msg.SubType == opts.Subtype {
				t, err := slackTimestampToTime(msg.Timestamp)
				if err != nil {
					log.Error(err)
				} else {
					ts := strings.Replace(msg.Timestamp, ".", "", 1)
					url := fmt.Sprintf("https://%s.slack.com/archives/%s/p%s", teamName, channelID, ts)
					fmt.Printf("%s: %s\n", t, msg.Text)
					fmt.Printf(" url: %s\n", url)
				}
				switch count {
				case 0:
					count++
				case 1:
					os.Exit(0)
				default:
					// unreachable
					log.Errorf("bug: count is %d\n", count)
					os.Exit(2)
				}
				log.Debug(pp.Sprint(msg))
			}
		}
		if r.ResponseMetaData.NextCursor != "" {
			cursor = r.ResponseMetaData.NextCursor
			log.Info("search next 200 history")
		} else {
			break
		}
	}

	fmt.Printf("Log data with target subtype [%s] found less than one\n", opts.Subtype)
	os.Exit(0)
}
