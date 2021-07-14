package main

import (
	"encoding/json"
	"flag"
	"fmt"
	melonadeClientGo "github.com/devit-tel/melonade-client-go"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var melonadeUrl = flag.String("url", "http://melonade.example.com", "melonade process manager endpoint")
	var mode = flag.String("mode", "", "dump or restore")
	var path = flag.String("path", "./backup", "files location")
	flag.Parse()

	switch *mode {
	case "dump":
		dump(melonadeClientGo.New(*melonadeUrl), *path)
	case "restore":
		restore(melonadeClientGo.New(*melonadeUrl), *path)
	case "upgrade":
		upgrade(melonadeClientGo.New(*melonadeUrl), *path)
	case "clean":
		clean(melonadeClientGo.New(*melonadeUrl), *path)
	default:
		fmt.Println(`unknown mode: only support "dump" and "restore"`)
	}
}

func dump(c melonadeClientGo.Service, p string) error {
	ts, err := c.GetTaskDefinitions()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("downloaded task definitions")

	os.MkdirAll(fmt.Sprintf("%s/tasks", p), os.ModePerm)

	for _, t := range ts {
		b, err := json.Marshal(t)
		if err != nil {
			log.Println(t.Name, err)
		}
		err = ioutil.WriteFile(fmt.Sprintf("%s/tasks/%s.json", p, t.Name), b, 0644)
		if err != nil {
			log.Println(t.Name, err)
		}
	}

	ws, err := c.GetWorkflowDefinitions()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("downloaded workflow definitions")
	for _, w := range ws {
		b, err := json.Marshal(w)
		if err != nil {
			log.Println(w.Name, w.Rev, err)
		}

		os.MkdirAll(fmt.Sprintf("%s/workflows/%s", p, w.Name), os.ModePerm)

		err = ioutil.WriteFile(fmt.Sprintf("%s/workflows/%s/%s.json", p, w.Name, w.Rev), b, 0644)
		if err != nil {
			log.Println(w.Name, err)
		}
	}

	return nil
}

func restore(c melonadeClientGo.Service, p string) error {
	err := filepath.Walk(fmt.Sprintf("%s/tasks", p),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				if strings.HasSuffix(path, ".json") {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					var t melonadeClientGo.TaskDefinition
					err = json.Unmarshal(b, &t)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}

					err = c.CreateTaskDefinition(t)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}

					fmt.Printf("restored task %s\n", t.Name)
				}
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}

	return filepath.Walk(fmt.Sprintf("%s/workflows", p),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				if strings.HasSuffix(path, ".json") {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					var w melonadeClientGo.WorkflowDefinition
					err = json.Unmarshal(b, &w)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}

					err = c.CreateWorkflowDefinition(w)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					fmt.Printf("restored workflow %s:%s\n", w.Name, w.Rev)
				}
			}
			return nil
		})
}

func upgrade(c melonadeClientGo.Service, p string) error {
	err := filepath.Walk(fmt.Sprintf("%s/tasks", p),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				if strings.HasSuffix(path, ".json") {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					var t melonadeClientGo.TaskDefinition
					err = json.Unmarshal(b, &t)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}

					err = c.UpdateTaskDefinition(t)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}

					fmt.Printf("restored task %s\n", t.Name)
				}
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}

	return filepath.Walk(fmt.Sprintf("%s/workflows", p),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				if strings.HasSuffix(path, ".json") {
					b, err := ioutil.ReadFile(path)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					var w melonadeClientGo.WorkflowDefinition
					err = json.Unmarshal(b, &w)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}

					err = c.UpdateWorkflowDefinition(w)
					if err != nil {
						fmt.Println(err.Error())
						return nil
					}
					fmt.Printf("restored workflow %s:%s\n", w.Name, w.Rev)
				}
			}
			return nil
		})
}

func clean(c melonadeClientGo.Service, p string) error {
	ws, err := c.GetWorkflowDefinitions()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("downloaded workflow definitions")
	for _, w := range ws {
		err := c.DeleteWorkflowDefinition(w.Name, w.Rev)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Printf("Deleted %s:%s", w.Name, w.Rev)
		}
	}

	return nil
}
