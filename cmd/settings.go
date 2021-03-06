package cmd

import (
	"log"
	"net/http"

	"github.com/salihciftci/liman/pkg/tool"
)

func settingsHandler(w http.ResponseWriter, r *http.Request) {
	err := parseSessionCookie(w, r)
	if err != nil {
		return
	}

	if r.Method == "POST" {
		pass := r.FormValue("cpass")

		match := tool.CheckPass(pass, userPassword)
		if !match {
			http.Redirect(w, r, "/settings", http.StatusFound)
			return
		}

		nPass := r.FormValue("npass")
		cNPass := r.FormValue("cnpass")

		if nPass != cNPass {
			http.Redirect(w, r, "/settings", http.StatusFound)
			return
		}

		match = tool.CheckPass(nPass, userPassword)

		if match {
			http.Redirect(w, r, "/settings", http.StatusFound)
			return
		}

		bNPass, err := tool.HashPasswordAndSave(nPass)
		if err != nil {
			log.Println(err)
			return
		}

		userPassword = string(bNPass)
		http.Redirect(w, r, "/logout", http.StatusFound)
	}

	bn, _ := getNotification()

	version, err := tool.Version()
	if err != nil {
		log.Println(err)
	}

	var data []interface{}
	data = append(data, bn)
	data = append(data, version)
	data = append(data, apiKey)

	err = tpl.ExecuteTemplate(w, "settings.tmpl", data)
	if err != nil {
		log.Println(r.Method, r.URL.Path, err)
	}
	log.Println(r.Method, r.URL.Path)
}
