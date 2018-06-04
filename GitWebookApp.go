package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/google/go-github/github"
	"os/exec"
)

func handleWebHook(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("headers: %v\n", r.Header)

	payload, err := github.ValidatePayload(r, []byte("02447eac-6779-11e8-adc0-fa7ae01bbebc"))
	if err != nil {
		log.Printf("error validating request body: err=%s\n", err)
		return
	}
	defer r.Body.Close()

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Printf("could not parse webhook: err=%s\n", err)
		return
	}

	switch e := event.(type) {
	case *github.PushEvent:
		 executeDeployScript(e.Repo.GetName(), e.Repo.GetCloneURL(), e.GetAfter())
		 log.Println(e)
	default:
		log.Printf("unknown event type %s\n", github.WebHookType(r))
		return
	}

}

func executeDeployScript(repo string, repo_url string,commit string) {
	log.Printf("calling deploy script %s, repo_url %s, commit %s", repo,repo_url,commit)
	cmdStr := "/home/ec2-user/deploy.sh "
	cmd := exec.Command("/bin/sh", "-c", cmdStr, repo, repo_url,commit)
	_, err := cmd.Output()
	if err != nil {
		log.Println("Error executing the deploy script:")
		log.Println(err.Error())
		return
	}
}

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w , "Healthy")
}

func main() {
	log.Println("server started")
	http.HandleFunc("/webhook",handleWebHook)
	http.HandleFunc("/health",health)
	log.Fatal(http.ListenAndServe(":9000",nil))
}
