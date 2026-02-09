package handler

import (
	"net/http"
	"server/archive/internal/dto"
	"server/archive/internal/logger"
	"server/archive/internal/rest/transport"
	"server/archive/internal/usecase"
)

type Progress struct {
	handler
	u usecase.ProgressInterface
}

func NewProgress(u usecase.ProgressInterface) (*Progress, error) {
	return &Progress{u: u}, nil
}

func (h *Progress) Get(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	data, err := h.u.Get(r.Context(), identity)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	resp := make([]*dto.ProgressShortResp, len(data))
	for i := range resp {
		resp[i] = dto.NewProgressShortResp(data[i])
	}

	log := logger.FromCtx(r.Context()).With(
		logger.NewTracedField("response", resp),
		logger.NewTracedField("data", data),
	)
	log.Debug("debug")

	if err := transport.EncodeJSON(w, &resp); err != nil {
		h.sendEncodeError(w)
		return
	}
}

func (h *Progress) PatchAnswer(w http.ResponseWriter, r *http.Request) {
	var req dto.AnswerPatchReq
	if err := transport.DecodeJSON(r.Body, &req); err != nil {
		h.sendDecodeError(w)
		return
	}
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	answerID, err := h.getParamAnswerID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	progressID, err := h.getParamProgressID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	err = h.u.UpdateAnswer(r.Context(), identity, &req, progressID, answerID)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Progress) PostStart(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	progressID, err := h.getParamProgressID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	err = h.u.Start(r.Context(), identity, progressID)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Progress) PostFinish(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	progressID, err := h.getParamProgressID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	err = h.u.Finish(r.Context(), identity, progressID)
	if err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Progress) PostAnswerCorrect(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	answerID, err := h.getParamAnswerID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	progressID, err := h.getParamProgressID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	if err := h.u.ReviewAnswer(r.Context(), identity, true, progressID, answerID); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Progress) PostAnswerIncorrect(w http.ResponseWriter, r *http.Request) {
	identity, ok := transport.NewIdentityFromCtx(r.Context())
	if !ok {
		transport.SendAuthDataError(w)
		return
	}
	answerID, err := h.getParamAnswerID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}
	progressID, err := h.getParamProgressID(r)
	if err != nil {
		h.sendParamError(w, err.Error())
		return
	}

	if err := h.u.ReviewAnswer(r.Context(), identity, false, progressID, answerID); err != nil {
		h.sendUsecaseError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
