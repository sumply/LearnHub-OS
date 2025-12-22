package transport

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/internal/usecase"
)

func decodeJSON(r io.ReadCloser, v any) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}
	return nil
}

func encodeJSON(w io.Writer, v any) error {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	return nil
}

func sendError(w http.ResponseWriter, statusCode int, what string) {
	body := map[string]any{
		"error": what,
	}
	w.WriteHeader(statusCode)
	encodeJSON(w, body)
}

func formatShortName(f string, l string, m string) string {
	if f == "" || l == "" {
		return ""
	}
	name := fmt.Sprintf("%s %v.", f, l[0])
	if m != "" {
		name = fmt.Sprintf("%s %v.", name, m[0])
	}
	return name
}

func toUserRole(s string) (usecase.UserRole, bool) {
	switch s {
	case "root":
		return usecase.Root, true
	case "admin":
		return usecase.Admin, true
	case "teacher":
		return usecase.Teacher, true
	case "student":
		return usecase.Student, true
	default:
		return usecase.Student, false
	}
}
