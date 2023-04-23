package handler

import (
	"encoding/json"
	"github.com/HaCaK/pse-bee-gobooking/src/property/model"
	"github.com/HaCaK/pse-bee-gobooking/src/property/service"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getPropertyFromRequest(r *http.Request) (*model.Property, error) {
	var property model.Property
	err := json.NewDecoder(r.Body).Decode(&property)
	if err != nil {
		log.Errorf("Can't decode request body to property struct: %v", err)
		return nil, err
	}
	return &property, nil
}

func CreateProperty(w http.ResponseWriter, r *http.Request) {
	property, err := getPropertyFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := service.CreateProperty(property); err != nil {
		log.Errorf("Error calling service CreateProperty: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sendJson(w, property)
}

func GetProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := service.GetProperties()
	if err != nil {
		log.Errorf("Error calling service GetProperties: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	sendJson(w, properties)
}

func GetProperty(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	property, err := service.GetProperty(id)
	if err != nil {
		log.Errorf("Failure retrieving property with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if property == nil {
		http.Error(w, "404 property not found", http.StatusNotFound)
		return
	}
	sendJson(w, property)
}

func UpdateProperty(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	property, err := getPropertyFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	property, err = service.UpdateProperty(id, property)
	if err != nil {
		log.Errorf("Failure updating property with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if property == nil {
		http.Error(w, "404 property not found", http.StatusNotFound)
		return
	}
	sendJson(w, property)
}

func DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id, err := getId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	property, err := service.DeleteProperty(id)
	if err != nil {
		log.Errorf("Failure deleting property with ID %v: %v", id, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if property == nil {
		http.Error(w, "404 property not found", http.StatusNotFound)
		return
	}
	sendJson(w, result{Success: "OK"})
}
