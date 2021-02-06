package main

import (
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

var (
	defaultQueueFile = "./queues.json"
	defaultRoutingFile = "./routings.json"

	saveSetting = kingpin.Command("save", "save queue and routing settings")
	from = saveSetting.Arg("from", "fireworq host").Required().String()
	saveQueueFile = saveSetting.Arg("queueFile", "queue setting save file").Default(defaultQueueFile).String()
	saveRoutingFile = saveSetting.Arg("routingFile", "routing setting save file").Default(defaultRoutingFile).String()

	applySetting = kingpin.Command("apply", "apply queue and routing settings")
	to = applySetting.Arg("to", "setting apply target fireworq host").Required().String()
	applyQueueFile = applySetting.Arg("queueFile", "apply queue setting file").Default(defaultQueueFile).String()
	applyRoutingFile = applySetting.Arg("routingFile", "apply routing setting file").Default(defaultRoutingFile).String()

)

func main() {
	app := kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("stk132")

	switch kingpin.Parse() {
	case "save":
		if err := saveQueueData(*saveQueueFile, *from); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := saveRoutingData(*saveRoutingFile, *from); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "apply":
		if err := applyQueue(*applyQueueFile, *to); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if err := applyRouting(*applyRoutingFile, *to); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	default:
		fmt.Println(app.Help)
	}
}
