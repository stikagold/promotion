package Controllers

import (
	"cpool/Configurator"
	"cpool/Helpers"
	"cpool/Response"
	"cpool/Services/Promotion"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type PromotionsController struct {
	promotionService Promotion.Promotion
	cfg              *Configurator.Cfg
}

func (c *PromotionsController) Initial(cfg *Configurator.Cfg) error {
	c.cfg = cfg
	c.promotionService.Initial(cfg)
	return nil
}

func (c *PromotionsController) Run(router *mux.Router) error {
	router.HandleFunc("/"+Helpers.PROMOTIONS_ENTITY+"/{id}", c.Get).Methods("GET")
	return nil
}

func (c *PromotionsController) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	var response Response.ApiResponse

	if err != nil {
		response.Code = http.StatusBadRequest
		response.Message = "Id should be an integer"
		Helpers.WriteAsResponse(response, http.StatusBadRequest, w)
		return
	}
	resp, err := c.promotionService.GetPromotion(id)
	if err != nil {
		response.Code = http.StatusNotFound
		response.Message = "Promotion not found"
		Helpers.WriteAsResponse(response, http.StatusNotFound, w)
		return
	}
	response.Code = http.StatusOK
	response.Data = resp.GetForResponse()
	Helpers.WriteAsResponse(response, http.StatusOK, w)
}
