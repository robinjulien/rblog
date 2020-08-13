package cli

// Config represents the config located in the config.json file
type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	BlogName string `json:"blogName"`
}

// Execute take args, and route it to the specific action function
func Execute(args []string) int {
	if len(args) > 1 {
		switch args[1] {
		case "help":
			help()
		case "init":
			initInstance()
		default:
			help()
		}
	} else {
		help()
	}
	return 0
}

// check is used to check if there is an error that should not happen, and thus that there is no recovery from that error
func check(err error) {
	if err != nil {
		panic(err)
	}
}
