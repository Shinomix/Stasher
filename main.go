package main

import (
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

type options struct {
	Message      string
	Duration     int
	IsReminder   bool
	ConfFilePath string
}

func registerParams() *options {
	var message = flag.String("message", "", "Optional message to add to the Slack message")
	var isReminder = flag.Bool("reminder", false, "Transform simple message to reminder")
	var duration = flag.Int("duration", 600, "Optional duration for the reminder (in seconds)")
	var confFilePath = flag.String("conf", os.Getenv("HOME")+"/.stasher", "Folder config file is located in (default to ~/.stasher/)")
	var showUsage = flag.Bool("help", false, "Display this usage")
	flag.Parse()

	if *showUsage {
		flag.Usage()
		os.Exit(0)
	}

	return &options{
		Message:      *message,
		Duration:     *duration,
		IsReminder:   *isReminder,
		ConfFilePath: *confFilePath,
	}
}

func simpleMessage(options *options, config *Config, message string) bool {
	params := url.Values{}
	params.Add("token", config.AccessToken)
	params.Add("channel", config.Username)
	params.Add("text", message)

	resp, err := http.Get(config.BaseURL + "/chat.postMessage?" + params.Encode())
	if err != nil {
		fmt.Println("[simpleMessage] Failed to call /chat.postMessage: ", err)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Println("[simpleMessage] /chat.postMessage returned with an error code: ", resp.StatusCode)
		return false
	}

	fmt.Println("> Message sent with success!")
	return true
}

func reminderMessage(options *options, config *Config, message string) bool {
	params := url.Values{}
	params.Add("token", config.AccessToken)
	params.Add("text", message)
	params.Add("user", GetSelfID(config))
	params.Add("time", fmt.Sprintf("%d", options.Duration))

	resp, err := http.Get(config.BaseURL + "/reminders.add?" + params.Encode())
	if err != nil {
		fmt.Println("[reminderMessage] Failed to call /reminders.add: ", err)
		return false
	}
	if resp.StatusCode != 200 {
		fmt.Println("[reminderMessage] /reminders.add returned with an error code: ", resp.StatusCode)
		return false
	}

	fmt.Println("> Reminder set with success!")
	return true
}

func sendNotification(options *options, config *Config, defaultMessage string) {
	var message string
	if options.Message != "" {
		message = options.Message
	} else {
		message = defaultMessage
	}

	if options.IsReminder && options.Duration > 0 {
		reminderMessage(options, config, message)
	} else {
		simpleMessage(options, config, message)
	}
}

func stashAndNotify(options *options, config *Config) {
	out, err := exec.Command("git", "stash").Output()
	if err != nil {
		panic(err.Error())
	}
	outString := fmt.Sprintf("%s", out)

	fmt.Println("> ", outString)
	if strings.Contains(outString, "No local changes to save") || strings.Contains(outString, "fatal:") {
		return
	}
	sendNotification(options, config, outString)
}

func main() {
	options := registerParams()
	config := Config{}
	if !LoadConfig(&config, options) {
		return
	}
	stashAndNotify(options, &config)
}
