package main

func main() {
	server := Server{
		Port:         env("PORT", "9501"),
		Store:        InitStore("/config"),
		DiscordID:    env("DISCORD_ID", ""),
		DiscordToken: env("DISCORD_TOKEN", ""),
	}
	defer server.Store.DB.Close()
	server.Start()
}
