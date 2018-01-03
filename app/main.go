package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/gorilla/websocket"
)

const (
	sessionName = "sicuro-auth"
)

var (
	port         = os.Getenv("PORT")
	appDIR       = filepath.Join(os.Getenv("ROOT_DIR"), "app")
	sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
	templates    = template.Must(template.ParseFiles(fetchTemplates()...))
	upgrader     = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func init() {
	setupGithubOauthCfg()
}

func main() {
	http.HandleFunc("/run", runCIHandler())
	http.HandleFunc("/show", showPageHandler())
	http.HandleFunc("/index", indexPageHandler())
	http.HandleFunc("/dashboard", dashboardPageHandler())
	http.HandleFunc("/ci/", ciPageHandler())
	http.HandleFunc("/gh/subscribe", githubSubscriptionHandler())

	http.HandleFunc("/ws/", wsHandler)

	http.HandleFunc("/gh/auth", ghAuth)
	http.HandleFunc("/gh/callback", ghAuthCallback)
	http.HandleFunc("/gh/webhook", githubWebhookHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", 303)
	})

	fmt.Printf("Starting server on port: %s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
