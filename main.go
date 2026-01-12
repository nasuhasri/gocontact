package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Contact struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

const dataFile = "contacts.json"

// ---------- Helpers ----------
func readContacts() ([]Contact, error) {
	file, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)
	var contacts []Contact
	json.Unmarshal(bytes, &contacts)
	return contacts, nil
}

func writeContacts(contacts []Contact) error {
	bytes, _ := json.MarshalIndent(contacts, "", "  ")
	return ioutil.WriteFile(dataFile, bytes, 0644)
}

// ---------- Handlers ----------
func getContacts(w http.ResponseWriter, r *http.Request) {
	contacts, _ := readContacts()
	json.NewEncoder(w).Encode(contacts)
}

func getContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	contacts, _ := readContacts()

	for _, c := range contacts {
		if c.ID == id {
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.NotFound(w, r)
}

func addContact(w http.ResponseWriter, r *http.Request) {
	var contact Contact
	json.NewDecoder(r.Body).Decode(&contact)
	contact.ID = uuid.New().String()

	contacts, _ := readContacts()
	contacts = append(contacts, contact)
	writeContacts(contacts)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(contact)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updated Contact
	json.NewDecoder(r.Body).Decode(&updated)

	contacts, _ := readContacts()
	for i, c := range contacts {
		if c.ID == id {
			updated.ID = id
			contacts[i] = updated
			writeContacts(contacts)
			json.NewEncoder(w).Encode(updated)
			return
		}
	}
	http.NotFound(w, r)
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	contacts, _ := readContacts()

	for i, c := range contacts {
		if c.ID == id {
			contacts = append(contacts[:i], contacts[i+1:]...)
			writeContacts(contacts)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.NotFound(w, r)
}

// ---------- Main ----------
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/contacts", getContacts).Methods("GET")
	r.HandleFunc("/contacts/{id}", getContact).Methods("GET")
	r.HandleFunc("/contacts", addContact).Methods("POST")
	r.HandleFunc("/contacts/{id}", updateContact).Methods("PUT")
	r.HandleFunc("/contacts/{id}", deleteContact).Methods("DELETE")

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
