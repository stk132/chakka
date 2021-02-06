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
	from = saveSetting.Flag("from", "fireworq host").Required().Short('f').String()
	saveQueueFile = saveSetting.Flag("queue-file", "queue setting save file").Default(defaultQueueFile).Short('q').String()
	saveRoutingFile = saveSetting.Flag("routing-file", "routing setting save file").Default(defaultRoutingFile).Short('r').String()

	applySetting = kingpin.Command("apply", "apply queue and routing settings")
	to = applySetting.Flag("to", "setting apply target fireworq host").Required().Short('t').String()
	applyQueueFile = applySetting.Flag("queueFile", "apply queue setting file").Default(defaultQueueFile).Short('q').String()
	applyRoutingFile = applySetting.Flag("routingFile", "apply routing setting file").Default(defaultRoutingFile).Short('r').String()

)

func main() {
	app := kingpin.UsageTemplate(kingpin.CompactUsageTemplate).Version("0.1").Author("stk132")
	kingpin.CommandLine.Help = "fireworq queue and routing setting save, apply cli"
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
