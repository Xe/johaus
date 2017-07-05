package main

import (
	"bytes"
	"flag"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"context"

	"github.com/eaburns/johaus/parser"
	"github.com/eaburns/johaus/pretty"
	prettyprint "github.com/eaburns/pretty"
	"github.com/velour/chat"
	"github.com/velour/chat/irc"

	// Register all supported dialects.
	_ "github.com/eaburns/johaus/parser/alldialects"
)

var (
	nick    = flag.String("n", "", "The bot's IRC nickname")
	pass    = flag.String("p", "", "The bot's IRC password")
	server  = flag.String("s", "irc.freenode.net:6697", "The IRC server")
	channel = flag.String("c", "#velour-test", "The IRC channel")
)

func main() {
	flag.Parse()

	ctx := context.Background()
	cl, err := irc.DialSSL(ctx, *server, *nick, *nick, *pass, false)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := cl.Close(ctx); err != nil {
			panic(err)
		}
	}()
	ch, err := cl.Join(ctx, *channel)
	if err != nil {
		panic(err)
	}
	if err := run(ctx, ch); err != nil {
		panic(err)
	}
}

func run(ctx context.Context, ch chat.Channel) error {
	for {
		ev, err := ch.Receive(ctx)
		if err != nil {
			return err
		}
		prettyprint.Print(ev)
		fmt.Println("")
		msg, ok := ev.(chat.Message)
		if !ok {
			continue
		}
		fixRelayedMessage(&msg)
		fmt.Println("fixed")
		prettyprint.Print(ev)
		fmt.Println("")
		text := strings.TrimSpace(msg.Text)
		if isGreeting(text) {
			if err := sayHi(ctx, msg); err != nil {
				return err
			}
			continue
		}
		if isParseRequest(text) {
			parseText(ctx, msg)
		}
	}
}

func isGreeting(text string) bool {
	const pat = `^coi(\s+|\.*)la(\s+|\.*)*johaus`
	ok, err := regexp.MatchString(pat, text)
	if err != nil {
		panic(err.Error())
	}
	return ok
}

func sayHi(ctx context.Context, msg chat.Message) error {
	to := msg.From.DisplayName
	return send(ctx, msg.Origin(),
		"coi la'o zoi. "+to+" .zoi mi'e sampre po'o"+
			" gi'e se finti la jelca")
}

const parseRequestPrefix = ".jo'au "

func isParseRequest(text string) bool {
	return strings.HasPrefix(text, parseRequestPrefix)
}

func parseText(ctx context.Context, msg chat.Message) error {
	text := strings.TrimSpace(msg.Text)[len(parseRequestPrefix):]
	const breaks = ". \t\n\r?!"
	i := strings.IndexAny(text, breaks)
	if i < 0 {
		return nil
	}
	dialect := text[:i]
	text = strings.TrimLeft(text[i:], breaks)
	ch := msg.Origin()
	if !knownDialects[dialect] {
		// TODO: more informative message, and print the known dialects.
		return send(ctx, ch, "la'o zoi. "+dialect+" .zoi mo")
	}
	const (
		maxReplyBytes = 450
		tooBigMsg     = ".u'u dukse lo ka clani"
	)
	// If the given text is greater than the max reply bytes
	// then all is hopeless, so give up now.
	if len(text) > maxReplyBytes {
		return send(ctx, ch, tooBigMsg)
	}
	tree, err := parser.Parse(dialect, text)
	if err != nil {
		parseErr, ok := err.(*parser.Error)
		if !ok {
			return send(ctx, ch, err.Error())
		}
		goodText := text[parseErr.Byte:]
		if len(goodText) > 10 {
			goodText = goodText[:10] + "…"
		}
		return send(ctx, ch, goodText+" "+err.Error())
	}
	parser.RemoveMorphology(tree)
	parser.AddElidedTerminators(tree)
	parser.RemoveSpace(tree)
	parser.CollapseLists(tree)
	b := bytes.NewBuffer(nil)
	pretty.Braces(b, tree)
	reply := b.String()
	if len(reply) > maxReplyBytes {
		return send(ctx, ch, tooBigMsg)
	}
	return send(ctx, ch, reply)
}

var knownDialects = func() map[string]bool {
	ds := make(map[string]bool)
	for _, d := range parser.Dialects() {
		ds[d.Name] = true
	}
	return ds
}()

// relayBots is the set of all known relay bot names.
var relayBots = map[string]bool{
	"^^^^":   true, // Telegram
	"kahai":  true, // Discord
	"kahai1": true,
}

// The Lojban channels contain various relay bots,
// bridging IRC with Slack, Telegram, and Discord, for example.
// This fixes the From field to replace the relay bot's info
// with that of the user who's message was relayed.
// In the case of a relayed message,
// the From field is set to an empty User,
// except Nick, Fullname, and DisplayName
// are replaced with the relayed user's nick,
// and Channel is set to the Message.Origin.
func fixRelayedMessage(msg *chat.Message) {
	if !relayBots[msg.From.Nick] {
		return
	}
	const nickOpen = '<'
	if r, _ := utf8.DecodeRuneInString(msg.Text); r != nickOpen {
		return
	}
	const nickClose = ">: "
	i := strings.Index(msg.Text, nickClose)
	if i < 0 {
		return
	}
	nick := strings.TrimSpace(msg.Text[1:i])
	msg.From = &chat.User{
		Nick:        nick,
		FullName:    nick,
		DisplayName: nick,
		Channel:     msg.Origin(),
	}
	msg.Text = msg.Text[i+len(nickClose):]
}

func send(ctx context.Context, ch chat.Channel, text string) error {
	msg := chat.Message{Text: text}
	_, err := ch.Send(ctx, msg)
	return err
}
