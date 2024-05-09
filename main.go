package main

import (
	"html/template"
	"net/http"
	"os"
	"strconv"
	"fmt"
)

type FormData struct {
	Name          string
	Email         string
	Subscribe     bool
	ContactMethod string
	Comments      string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form data", http.StatusInternalServerError)
			return
		}

		formData := FormData{
			Name:          r.FormValue("name"),
			Email:         r.FormValue("email"),
			Subscribe:     r.FormValue("subscribe") == "on",
			ContactMethod: r.FormValue("contact_method"),
			Comments:      r.FormValue("comments"),
		}

		err = writeFormDataToFile(formData)
		if err != nil {
			http.Error(w, "Error writing form data to file", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./index.html")
	if err != nil {
		http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func writeFormDataToFile(formData FormData) error {
	file, err := os.OpenFile("formdata.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString("Name: " + formData.Name + "\n" +
		"Email: " + formData.Email + "\n" +
		"Subscribe: " + strconv.FormatBool(formData.Subscribe) + "\n" +
		"Contact Method: " + formData.ContactMethod + "\n" +
		"Comments: " + formData.Comments + "\n\n")
	if err != nil {
		return err
	}

	return nil
}

func main() {
	http.HandleFunc("/", formHandler)

	fmt.Println("Listening on:")
	fmt.Println("http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
