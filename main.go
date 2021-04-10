package main

func main() {
	server := Server{
		Port:  env("PORT", "9501"),
		Store: InitStore("/config"),
	}
	defer server.Store.DB.Close()
	server.Start()
}
