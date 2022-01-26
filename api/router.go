package api

import (
	"github.com/julienschmidt/httprouter"
)

func GetRouter() *httprouter.Router {
	router := httprouter.New()
	router.GET("/merchant", handleGetMerchants)
	router.POST("/merchant", handleAddMerchants)
	router.DELETE("/merchant/:code", handleDeleteMerchant)
	router.POST("/merchant/:code/addmember", handleAddMember)
	router.DELETE("/merchant/:code/:email", handleDeleteMember)
	router.GET("/merchant/:code/members", handleGetMember)
	return router
}
