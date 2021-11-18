package main

func main() {
	engine := Engine{}

	engine.InitCmdParams()
	engine.LoadConfig()
	engine.ParseConfig()

	engine.ProcessCommands()

	engine.PrintResult()
}