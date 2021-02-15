package main

import (
	"flag"
	melonadeClientGo "github.com/devit-tel/melonade-client-go"
	"log"
)

func main() {
	var oldMelonade = flag.String("old-url", "http://melonade.example.com", "old melonade process manager url")
	var newMelonade = flag.String("new-url", "http://melonade-new.example.com", "new melonade process manager url")
	flag.Parse()

	cOld := melonadeClientGo.New(*oldMelonade)
	cNew := melonadeClientGo.New(*newMelonade)

	ts, err := cOld.GetTaskDefinitions()
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range ts {
		err := cNew.SetTaskDefinition(*t)
		if err != nil {
			log.Println(err, *t)
		}
	}

	ws, err := cOld.GetWorkflowDefinitions()
	if err != nil {
		log.Fatal(err)
	}

	for _, w := range ws {
		err := cNew.SetWorkflowDefinition(*w)
		if err != nil {
			log.Println(err, *w)
		}
	}
}
