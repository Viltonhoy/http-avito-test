package server

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"http-avito-test/internal/generated"
	"http-avito-test/internal/storage"
	"io/ioutil"
	"net/http"
	"os"

	"go.uber.org/zap"
)

const reportLink = "http://localhost:4000"

func (h *Handler) MonthlyReport(w http.ResponseWriter, r *http.Request) {
	var hand *generated.MonthlyReportRequest

	body, _ := ioutil.ReadAll(r.Body)

	err := json.Unmarshal(body, &hand)
	if err != nil {
		http.Error(w, "malformed request body", http.StatusBadRequest)
		return
	}

	report, err := h.Store.MonthlyReport(r.Context(), int64(hand.Year), int64(hand.Month))
	if err != nil {
		if errors.Is(err, storage.ErrNoRecords) {
			http.Error(w, "records do not exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "error reading monthly report", http.StatusInternalServerError)
		return
	}

	fileFormat := fmt.Sprintf(`../../file_storage/consolidated_report%d-%d.csv`, hand.Year, hand.Month)

	file, err := os.Create(fileFormat)
	if err != nil {
		h.Logger.Error("openinп CSV file error", zap.Error(err))
		http.Error(w, "failed to open CSV file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	write := csv.NewWriter(file)
	err = write.WriteAll(report)
	if err != nil {
		h.Logger.Error("writing to CSV file error", zap.Error(err))
		http.Error(w, "cannot write to CSV file", http.StatusInternalServerError)
		return
	}

	result := generated.MonthlyReportResponse{
		Result: struct {
			Link string "json:\"link\""
		}{
			Link: reportLink,
		},
		Status: "ok",
	}

	marshalledRequest, err := json.Marshal(result)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, writeErr := w.Write(marshalledRequest)
	if err != nil {
		h.Logger.Error("failed to write connection", zap.Error(writeErr))
		return
	}
}
