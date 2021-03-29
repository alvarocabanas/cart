package rest

import (
	cart "cart/internal"
	"cart/internal/creator"
	"cart/internal/getter"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"error"`
}

type CartHandler struct {
	cartCreator creator.CartCreator
	cartGetter  getter.CartGetter
}

// NewCartHandler creates the handler/controller for the api, in future iterations there could be one handler per
// application service
func NewCartHandler(cartCreator creator.CartCreator, cartGetter getter.CartGetter) CartHandler {
	return CartHandler{
		cartCreator: cartCreator,
		cartGetter:  cartGetter,
	}
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body creator.AddItemDTO

	inData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = json.Unmarshal(inData, &body)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = h.cartCreator.AddItem(r.Context(), body)
	if err != nil {
		switch err {
		case cart.ErrItemNotFound:
			h.writeErrorResponse(w, http.StatusNotFound, err)
		case cart.ErrWrongAddItemQuantity:
			h.writeErrorResponse(w, http.StatusBadRequest, err)
		default:
			h.writeErrorResponse(w, http.StatusInternalServerError, err)
		}
		return
	}

	h.writeResponse(w, http.StatusCreated, nil)
}

func (h *CartHandler) GetCartStatus(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	cartStatusDTO := h.cartGetter.Get(r.Context())
	cartStatusResponse, _ := json.Marshal(cartStatusDTO)

	h.writeResponse(w, http.StatusOK, cartStatusResponse)
}

func (h *CartHandler) writeResponse(w http.ResponseWriter, status int, data []byte) {
	if data != nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.WriteHeader(status)
	if data != nil {
		w.Write(data)
	}
}

func (h *CartHandler) writeErrorResponse(w http.ResponseWriter, status int, err error) {
	errorResponse := ErrorResponse{
		Status:  status,
		Message: err.Error(),
	}

	r, _ := json.Marshal(errorResponse)

	h.writeResponse(w, status, r)
}
