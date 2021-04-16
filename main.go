package main

func main() {
	server := Server{
		Port:  env("SERVER_PORT", "9501"),
		Store: InitStore("/config"),
	}
	defer server.Store.DB.Close()
	server.Start()
}
