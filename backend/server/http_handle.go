package server

import (
	"context"
	"encoding/json"
	"github.com/XHXHXHX/medical_marketing/util/excel"
	"io"
	"net/http"
)

func (s *Server) ImportConsumer() http.HandlerFunc {
	filename := "file"
	return func(w http.ResponseWriter, r *http.Request) {
		buffer, err := excel.GetFileContentFromUploadFile(r, filename, 0)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		ctx := context.Background()

		result, err := s.reportService.Import(ctx, buffer)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		b, err := json.Marshal(result)
		if err != nil {
			_, _ = io.WriteString(w, err.Error())
			return
		}

		_, _ = io.WriteString(w, string(b))
	}
}
