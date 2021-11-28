package main

import (
	"log"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type AccountData struct {
	Attributes     *AccountAttributes `json:"attributes,omitempty"`
	ID             string             `json:"id,omitempty"`
	OrganisationID string             `json:"organisation_id,omitempty"`
	Type           string             `json:"type,omitempty"`
	Version        *int64             `json:"version,omitempty"`
}

type AccountAttributes struct {
	AccountClassification   *string  `json:"account_classification,omitempty"`
	AccountMatchingOptOut   *bool    `json:"account_matching_opt_out,omitempty"`
	AccountNumber           string   `json:"account_number,omitempty"`
	AlternativeNames        []string `json:"alternative_names,omitempty"`
	BankID                  string   `json:"bank_id,omitempty"`
	BankIDCode              string   `json:"bank_id_code,omitempty"`
	BaseCurrency            string   `json:"base_currency,omitempty"`
	Bic                     string   `json:"bic,omitempty"`
	Country                 *string  `json:"country,omitempty"`
	Iban                    string   `json:"iban,omitempty"`
	JointAccount            *bool    `json:"joint_account,omitempty"`
	Name                    []string `json:"name,omitempty"`
	SecondaryIdentification string   `json:"secondary_identification,omitempty"`
	Status                  *string  `json:"status,omitempty"`
	Switched                *bool    `json:"switched,omitempty"`
}

var accounts map[string]AccountData

// http handler
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	log.Println("Request recived: ", r.Method)

	switch r.Method {
	case http.MethodGet:
		fetch_service(w, r)
		break
	case http.MethodPost:
		add_service(w, r)
		break
    case http.MethodDelete:
        delete_service(w, r)
        break
	default:
		w.WriteHeader(405)
		w.Write([]byte("Method not allowed"))
		break
	}

}

func add_service(w http.ResponseWriter, r *http.Request) {
    account_ids, ok := r.URL.Query()["account_id"]
    
    if !ok || len(account_ids[0]) < 1 {
        log.Println("Url Param 'account_id' is missing")
        return
    }
    account_id := account_ids[0]
	payload, _ := ioutil.ReadAll(r.Body)

	var account AccountData

	err := json.Unmarshal(payload, &account)

	if err != nil || account.ID != account_id {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}
    accounts[account_id] = account

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(201)

	json, _ := json.Marshal(account)
	w.Write(json)

	log.Println("Response returned: ", 201)
}

func delete_service(w http.ResponseWriter, r *http.Request) {
    account_ids, ok := r.URL.Query()["account_id"]
    
    if !ok || len(account_ids[0]) < 1 {
        log.Println("Url Param 'account_id' is missing")
        return
    }
    account_id := account_ids[0]

	if account_id == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}

    account, ok := accounts[account_id]
    if ! ok {
        w.WriteHeader(400)
        w.Write([]byte("Bad Request"))
        return
    }
    delete(accounts, account_id)
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(201)

	json, _ := json.Marshal(account)
	w.Write(json)

	log.Println("Response returned: ", 201)
}

func fetch_service(w http.ResponseWriter, r *http.Request) {
    account_ids, ok := r.URL.Query()["account_id"]
    
    if !ok || len(account_ids[0]) < 1 {
        log.Println("Url Param 'account_id' is missing")
        return
    }
    account_id := account_ids[0]

	if account_id == "" {
		w.WriteHeader(400)
		w.Write([]byte("Bad Request"))
		return
	}
	account, ok := accounts[account_id]
    if ! ok {
        w.WriteHeader(400)
        w.Write([]byte("Bad Request"))
        return
    }
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(200)

	json, _ := json.Marshal(account)
	w.Write(json)

	log.Println("Response returned: ", 200)
}

func main() {
	http.HandleFunc("/v1/organisation/accounts/{account_id}", HandleRequest)

	err := http.ListenAndServe(":8082", nil)
	if err != nil {
		log.Println(err)
		return
	}
}
