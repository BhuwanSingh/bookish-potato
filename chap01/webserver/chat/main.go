package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once          // ensures that the template is only parsed once
	filename string             // stores the name of the HTML file
	templ    *template.Template // holds the parsed template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		// Must is a helper that wraps a call to a function returning (*Template, error) and panics if the error is non-nil.
		// It is intended for use in variable initializations such as
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates", t.filename)))
	})
	// holds the template Data, including the request HOST.
	data := map[string]interface{}{
		"Host": r.Host,
	}
	// finds the cookie with "auth"
	// If found, extracts user data and decodes it.
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	// renders the template with the data.
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	gomniauth.SetSecurityKey("PUT your auth key here.")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Could not get the env variables.")
	}

	googleKey := os.Getenv("GOOGLE_KEY")
	googleSecret := os.Getenv("GOOGLE_SECRET")

	if googleKey == "" || googleSecret == "" {
		log.Fatal("Google credentials not set in .env file")
	}

	gomniauth.WithProviders(
		facebook.New("key", "secret", "http://localhost:8080/auth/callback/facebook"),
		github.New("key", "secret", "http://localhost:8080/auth/callback/github"),
		google.New(googleKey, googleSecret, "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom(UseAuthAvatar)
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc(("/logout"), func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	go r.run()
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("Could not start server at :8080", err)
	}
}
