package main

import (
	"log"
	"net/http"
	"os/exec"
)

type ASync struct {
	cmdPrefix string
	cmds      chan string
}

func NewASync(cmdPrefix string, cmdsChanSize int) *ASync {
	a := ASync{
		cmdPrefix: cmdPrefix,
		cmds:      make(chan string, cmdsChanSize),
	}
	go func() {
		for {
			select {
			case cmd := <-a.cmds:
				go func() {
					if a.cmdPrefix != "" {
						cmd = a.cmdPrefix + " " + cmd
					}
					log.Printf("RUN \"%s\"", cmd)
					out, err := exec.Command("sh", "-c", cmd).Output()
					if err != nil {
						log.Printf("RUN \"%s\" error: %s", cmd, err)
					} else {
						log.Printf("RUN \"%s\" output: %s", cmd, out)
					}
					// exec.Command("sh", "-c", cmd).Run()
				}()
			}
		}
	}()
	return &a
}

func (a *ASync) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "post only", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	cmd := r.PostFormValue("cmd")
	a.cmds <- cmd
	log.Printf("POST \"%s\"", cmd)
}
