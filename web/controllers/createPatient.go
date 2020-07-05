package controllers

import (
	"net/http"
)

// CreatePatienthandler : controller to createPatient
func (app *Application) CreatePatienthandler(w http.ResponseWriter, r *http.Request) {
	data := &struct {
		TransactionID string
		Success       bool
		Response      bool
	}{
		TransactionID: "",
		Success:       false,
		Response:      false,
	}
	if r.FormValue("submitted") == "true" {
		fullName := r.FormValue("firstName")
		gender := r.FormValue("gender")
		birthday := r.FormValue("birthday")

		txid, err := app.Fabric.CreateEHR(fullName, fullName, gender, birthday)
		if err != nil {
			http.Error(w, "Unable to invoke createEHR in the blockchain : "+err.Error(), 500)
		}
		data.TransactionID = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "createPatient.html", data)
}
