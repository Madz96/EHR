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
		// firstName := r.FormValue("firstName")
		// lastName := r.FormValue("lastName")
		// gender := "gender"
		// birthday := r.FormValue("birthday")
		// address := r.FormValue("address")
		// contactNo := r.FormValue("contactNo")

		txid, err := app.Fabric.CreatePatient("firstName", "lastName", "contactNo", "gender", "2006-01-01", "address")
		if err != nil {
			http.Error(w, "Unable to invoke createPatient in the blockchain : "+err.Error(), 500)
		}
		data.TransactionID = txid
		data.Success = true
		data.Response = true
	}
	renderTemplate(w, r, "createPatient.html", data)
}
