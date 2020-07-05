package web

import (
	"fmt"
	"net/http"

	"github.com/IMS+/EHR/web/controllers"
)

// Serve : start the web app and server
func Serve(app *controllers.Application) {
	fs := http.FileServer(http.Dir("web/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))

	http.HandleFunc("/home.html", app.HomeHandler)
	http.HandleFunc("/request.html", app.RequestHandler)
	http.HandleFunc("/createEHR.html", app.CreateEHRhandler)
	http.HandleFunc("/getEHR.html", app.GetEHRhandler)
	http.HandleFunc("/updateEHR.html", app.UpdateEHRhandler)
	http.HandleFunc("/createPatient.html", app.CreatePatienthandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/createEHR.html", http.StatusTemporaryRedirect)
	})

	fmt.Println("Listening (http://localhost:3000/) ...")
	http.ListenAndServe(":3000", nil)
}
