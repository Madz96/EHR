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
		name := r.FormValue("name")
		contactNo := r.FormValue("contactNo")

		txid, err := app.Fabric.CreatePatient(name, contactNo)
		if err != nil {
			http.Error(w, "Unable to invoke createPatient in the blockchain : "+err.Error(), 500)
		}
		data.TransactionID = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "createPatient.html", data)
}
