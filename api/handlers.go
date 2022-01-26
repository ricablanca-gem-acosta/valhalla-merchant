package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func jsonResponse(w http.ResponseWriter, body interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	p, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(p))
}

func messageResponse(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	type Message struct {
		Message string `json:"message"`
	}
	msg := Message{message}
	p, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, string(p))
}

func handleGetMerchants(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	merchants, err := listMerchants()
	if err != nil {
		messageResponse(w, "Failed to get merchant list", http.StatusInternalServerError)
		return
	}
	codes := []string{}
	for _, m := range merchants {
		codes = append(codes, m.Code)
	}
	jsonResponse(w, codes, http.StatusOK)
}

func handleAddMerchants(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		messageResponse(w, "Invalid input. No body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var merchant Merchant
	err = json.Unmarshal(body, &merchant)
	if err != nil {
		messageResponse(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	gottenMerchant, _ := getMerchant(merchant.Code)
	if gottenMerchant != nil {
		messageResponse(w, "Merchant with the same code is already created.", http.StatusConflict)
		return
	}
	_, err = createMerchant(&merchant.Code)
	if err != nil {
		messageResponse(w, "Failed to add merchant", http.StatusInternalServerError)
		return
	}
	messageResponse(w, "Merchant added.", http.StatusCreated)
}

func handleDeleteMerchant(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	merchantCode := p.ByName("code")
	merchant, _ := getMerchant(merchantCode)
	if merchant == nil {
		messageResponse(w, "Merchant does not exist.", http.StatusNotFound)
		return
	}
	err := deleteMerchant(merchant.Code)
	if err != nil {
		messageResponse(w, "Failed to delete merchant", http.StatusInternalServerError)
		return
	}
	messageResponse(w, "Merchant deleted.", http.StatusOK)
}

func handleAddMember(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		messageResponse(w, "Invalid input. No body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	merchantCode := p.ByName("code")
	merchant, err := getMerchant(merchantCode)
	if err != nil {
		fmt.Println(err)
		messageResponse(w, "Merchant does not exist", http.StatusNotFound)
		return
	}
	var member Member
	err = json.Unmarshal(body, &member)
	if err != nil {
		messageResponse(w, "Invalid input format", http.StatusBadRequest)
		return
	}
	for _, m := range merchant.Members {
		if m.Email == member.Email {
			messageResponse(w, "Member with the same email is already created.", http.StatusConflict)
			return
		}
	}
	err = createMember(merchantCode, member.Email)
	if err != nil {
		messageResponse(w, "Failed to add member", http.StatusInternalServerError)
		return
	}
	messageResponse(w, "Member added.", http.StatusCreated)
}

func handleDeleteMember(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	merchantCode := p.ByName("code")
	merchant, err := getMerchant(merchantCode)
	if err != nil {
		fmt.Println(err)
		messageResponse(w, "Merchant does not exist", http.StatusNotFound)
		return
	}
	foundMember := false
	email := p.ByName("email")
	for _, m := range merchant.Members {
		if m.Email == email {
			foundMember = true
			break
		}
	}
	if !foundMember {
		messageResponse(w, "Member does not exist on this merchant", http.StatusNotFound)
		return
	}
	err = deleteMember(merchantCode, email)
	if err != nil {
		messageResponse(w, "Failed to delete member", http.StatusInternalServerError)
		return
	}
	messageResponse(w, "Member deleted.", http.StatusOK)
}

func handleGetMember(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	itemPerPage := 20
	merchantCode := p.ByName("code")
	merchant, err := getMerchant(merchantCode)
	if err != nil {
		fmt.Println(err)
		messageResponse(w, "Merchant does not exist", http.StatusNotFound)
		return
	}
	type pageResult struct {
		Page       int      `json:"page"`
		TotalPages int      `json:"totalPages"`
		Count      int      `json:"count"`
		Data       []Member `json:"data"`
	}
	page := 1
	pageQ := r.URL.Query().Get("page")
	if pageQ != "" {
		page, err = strconv.Atoi(pageQ)
		if err != nil {
			messageResponse(w, "Invalid input format", http.StatusBadRequest)
			return
		}
	}
	totalItems := len(merchant.Members)
	totalPages := int(math.Ceil(float64(totalItems) / float64(itemPerPage)))
	if page > int(totalPages) || page < 1 {
		messageResponse(w, "Page out of bounds", http.StatusNotFound)
		return
	}
	data := []Member{}
	if totalItems > 0 {
		low := itemPerPage * (page - 1)
		high := low + itemPerPage
		if page == totalPages {
			high = low + (totalPages % itemPerPage)
		}
		data = merchant.Members[low:high]
	}
	count := len(data)
	res := pageResult{
		Page:       page,
		TotalPages: totalPages,
		Data:       data,
		Count:      count,
	}
	jsonResponse(w, res, http.StatusOK)
}
