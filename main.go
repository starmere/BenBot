package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    "time"

    // To use uptimerobot to ping
    "net/http"

    "github.com/bwmarrin/discordgo"
    "github.com/replit/database-go"
    "github.com/joho/godotenv"
)

type void struct {}
var member void
var bLset = map[string]void {}
var size int

func main() {
    godotenv.Load()

    //set := map[string] void {} // New empty set
    //set["Foo"] = member // Add
    //for i := range set { // Loop
    //    fmt.Println(k)
    //}
    //delete(set, "Foo") // Delete
    //size := len(set) // Size
    //_, exists := set["Foo"] // Membership
    //database.Set("bL0", "")

    keys, _ := database.ListKeys("bL")
    for _, key := range keys {
        value, _ := database.Get(key)
        bLset[value] = member
    }
    value, _ := database.Get("size")
    size, _ = strconv.Atoi(value)

    web_server()

    // Create a new Discord session using the bot token from .env
    bot, err := discordgo.New(fmt.Sprintf("Bot %v", os.Getenv("TOKEN")))
    if err != nil {
        panic(err)
    }

    // register events
    bot.AddHandler(ready)
    bot.AddHandler(messageCreate)

    err = bot.Open()
    if err != nil {
        fmt.Println("Error opening Discord session: ", err)
    }
    // Wait here until Ctrl-C or other term signal is received.
    fmt.Println("Bot is now running. Press Ctrl-C to exit.")
    for {}

    // Cleanly close down the Discord session.
    bot.Close()
}

// Uses uptimerobot to ping
func web_server() {
    server_msg := "Hello world!"
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, server_msg)
        fmt.Println("Received HTTP request from web server")
    })

    go http.ListenAndServe(":8080", nil)
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
    //t := UpdateStatusData{Status: "what my creator wills (in Go)"}
    //UpdateStatusComplex(s, event)
    s.UpdateGameStatus(0, "!ben")
    fmt.Printf("Logged in as\n%s\n%s\n--------\n", s.State.User.Username, s.State.User.ID)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    if m.Author.ID != s.State.User.ID {
        msg := strings.ToLower(m.Content)
        footer := &discordgo.MessageEmbedFooter {
            Text:         "BenBot by starmere#7058",
            IconURL:      "https://www.jing.fm/clipimg/full/138-1380087_angry-bear-png-free-angry-cartoon-bear.png",
        }

        if msg == "!ben" {
            embed := &discordgo.MessageEmbed {
                Title:     "BenBot",
                Description: "A moderating bot by <@590933358637350916> (starmere#7058)\n\nUse `!help` to view available commands",
                Color:       0x964b00, // Brown
                Thumbnail: &discordgo.MessageEmbedThumbnail {
                    URL: "https://www.jing.fm/clipimg/full/138-1380087_angry-bear-png-free-angry-cartoon-bear.png",
                },
                Footer: footer,
                Fields: []*discordgo.MessageEmbedField{
                    &discordgo.MessageEmbedField {
                        Name:   ":link: Links",
                        Value:  "**Source code**: https://github.com/starmere/BenBot\n**Official Discord server**: https://discord.gg/J3nJGhm2wh",
                        Inline: false,
                    },
                },
            }
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
        } else if msg == "!help" {
            embed := &discordgo.MessageEmbed {
                Title:     "Commands <@590933358637350916>",
                Description: "<@590933358637350916>\n!help : Get commands\n!bl : Get blacklist\n!blc : Clear database\n!bls : Get size of blacklist\n!bla [s] : Add to blacklist\n!bld [s] : Delete from blacklist\n\n**Mod commands**\n\n**Future commands**",
                Color:       0x00ff00, // Red
                Footer: footer,
            }
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
        } else if msg == "!bl" {
            if len(bLset) == 0 {
                s.ChannelMessageSend(m.ChannelID, "Blacklist: []")
            } else {
                keys := []string{}
                for key := range bLset {
                    keys = append(keys, key)
                }
                s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Blacklist: ['%s']", strings.Join(keys, "' '")))
            }
        } else if msg == "!blc" {
            keys, _ := database.ListKeys("bL")
            for _, key := range keys {
                database.Delete(key)
            }
            bLset = make(map[string] void)
            database.Set("size", "0")
            size = 0
            s.ChannelMessageSend(m.ChannelID, "Database cleared")
        } else if msg == "!bls" {
            s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Size of blacklist: %d", size))
        } else if msg == "!emb" {
            embed := &discordgo.MessageEmbed {
                Title:     "I am an Embed",
                Description: "This is a discordgo embed",
                Timestamp: time.Now().Format(time.RFC3339), // Discord wants ISO8601; RFC3339 is an extension of ISO8601 and should be completely compatible.
                Color:       0x00ff00, // Green
                //Author:      &discordgo.MessageEmbedAuthor {},
                Footer: footer,
                Image: &discordgo.MessageEmbedImage{
                    URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
                },
                Thumbnail: &discordgo.MessageEmbedThumbnail{
                    URL: "https://cdn.discordapp.com/avatars/119249192806776836/cc32c5c3ee602e1fe252f9f595f9010e.jpg?size=2048",
                },
                Fields: []*discordgo.MessageEmbedField{
                    &discordgo.MessageEmbedField {
                        Name:   "I am a field",
                        Value:  "I am a value",
                        Inline: false,
                    },
                    &discordgo.MessageEmbedField {
                        Name:   "I am a second field",
                        Value:  "I am a value",
                        Inline: false,
                    },
                },
            }
            s.ChannelMessageSendEmbed(m.ChannelID, embed)
        } else if len(msg) > 5 {
            if msg[:5] == "!bla " {
                e := msg[5:] // !bla jfdkslaf -> ['', 'jfdkslaf']
                _, exists := bLset[e]
                if exists {
                    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' already in blacklist", e))
                } else {
                    // add the new string to bLset
                    bLset[e] = member
                    // add the string to database
                    database.Set(fmt.Sprintf("bL%d", size), e)
                    size++
                    database.Set("size", fmt.Sprintf("%d", size))
                    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' added to blacklist", e))
                }
            } else if msg[:5] == "!bld " {
                e := msg[5:]
                _, exists := bLset[e]
                if exists {
                    delete(bLset, e)
                    keys, _ := database.ListKeys("bL")
                    for _, key := range keys {
                        database.Delete(key)
                    }
                    for key := range bLset {
                        database.Set(fmt.Sprintf("bL%d", size), key)
                    }
                    size--
                    database.Set("size", fmt.Sprintf("%d", size))
                    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' deleted from blacklist", e))
                } else {
                    s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("'%s' not in blacklist", e))
                }
            }
        } else {
            for key := range bLset {
                if strings.Contains(msg, key) {
                    s.ChannelMessageDelete(m.ChannelID, m.ID)
                    //s.ChannelMessageSend(m.ChannelID, "Message deleted")
                    embed := &discordgo.MessageEmbed {
                        Description: "Message deleted",
                        Footer: footer,
                    }
                    s.ChannelMessageSendEmbed(m.ChannelID, embed)
                    break
                }
            }
        }

        if msg == "!ping" {
            s.ChannelMessageSend(m.ChannelID, "pong")
        }
        if msg == "!hello" {
            s.ChannelMessageSend(m.ChannelID, "Stop right there, criminal scum!")
        }
}
