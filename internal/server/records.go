package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/tim3-p/gophkeeper/internal/common"
	"github.com/tim3-p/gophkeeper/internal/store"
)

func listRecords(w http.ResponseWriter, r *http.Request) {
	log.Print("listRecords")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	records, err := serverStore.ListRecords(user)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(records)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func listRecordsByType(w http.ResponseWriter, r *http.Request) {
	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	log.Print("listRecordsByType " + recordType)

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	records, err := serverStore.ListRecordsByType(user, recordType)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(records)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func getRecordByID(w http.ResponseWriter, r *http.Request) {
	log.Print("getRecordByID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
		return
	}

	record, err := serverStore.GetRecordByID(user, int64(id))
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record id %d not found", id)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(record)

	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func getRecordID(w http.ResponseWriter, r *http.Request) {
	log.Print("getRecordID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	var resp common.StoreRecordResponse
	var err error
	resp.ID, err = serverStore.GetRecordID(user, recordType, recordName)
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record %s of type %s not found",
			recordName, recordType)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(resp)

	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func getRecordByTypeName(w http.ResponseWriter, r *http.Request) {
	log.Print("getRecordTypeName")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	record, err := serverStore.GetRecordByTypeName(user, recordType, recordName)
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record %s of type %s not found",
			recordName, recordType)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
	err = json.NewEncoder(w).Encode(record)

	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func deleteRecordByID(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteRecordID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
		return
	}

	err = serverStore.DeleteRecordByID(user, int64(id))
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record id %d not found", id)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	writeStatus(w, http.StatusOK, "OK")
}

func deleteRecordByTypeName(w http.ResponseWriter, r *http.Request) {
	log.Print("deleteRecordTypeName")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	err := serverStore.DeleteRecordByTypeName(user, recordType, recordName)
	if err == store.ErrNotFound {
		msg := fmt.Sprintf("Record %s of type %s not found",
			recordName, recordType)
		log.Print(msg)
		writeStatus(w, http.StatusNotFound, msg)
		return
	}
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	writeStatus(w, http.StatusOK, "OK")
}

func storeRecord(w http.ResponseWriter, r *http.Request) {
	log.Print("storeRecord")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var resp common.StoreRecordResponse
	resp.Status = "OK"

	var record common.Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}
	resp.Name = record.Name
	resp.ID, err = serverStore.StoreRecord(user, record)

	if err != nil {
		log.Printf("storeRecord() error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Store Record: %v", err),
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func updateRecordByID(w http.ResponseWriter, r *http.Request) {
	log.Print("updateRecordID")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		writeStatus(w, http.StatusBadRequest, "cannot parse 'id' param")
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var resp common.StoreRecordResponse
	resp.Status = "OK"
	var record common.Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}
	resp.Name = record.Name
	err = serverStore.UpdateRecordByID(user, int64(id), record)

	if err != nil {
		log.Printf("update record error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Update Record: %v", err),
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}

func updateRecordByTypeName(w http.ResponseWriter, r *http.Request) {
	log.Print("updateRecordTypeName")

	user, _, ok := r.BasicAuth()
	if !ok {
		writeStatus(w, http.StatusBadRequest, "no basic auth")
		return
	}

	recordType := common.RecordType(chi.URLParam(r, "record_type"))
	recordName := chi.URLParam(r, "record_name")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print(err)
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}

	var resp common.StoreRecordResponse
	resp.Status = "OK"
	var record common.Record
	err = json.Unmarshal(body, &record)
	if err != nil {
		writeStatus(w,
			http.StatusBadRequest,
			fmt.Sprintf("Cannot Parse Body: %v", err),
		)
		return
	}
	resp.Name = record.Name
	err = serverStore.UpdateRecordByTypeName(user, recordType, recordName, record)

	if err != nil {
		log.Printf("update record error: %v", err)
		resp.Status = "error"
		writeStatus(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Cannot Update Record: %v", err),
		)
		return
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeStatus(w,
			http.StatusInternalServerError,
			"Internal Server Error",
		)
		return
	}
}
