package main

import (
	"strconv"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/Esseh/notorious-dev/PM"
	"github.com/Esseh/notorious-dev/CORE"
	"github.com/Esseh/notorious-dev/CONTEXT"
)

func INIT_PM_HANDLERS(r *httprouter.Router) {
	r.GET("/pm", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		ctx := CONTEXT.NewContext(res,req)
		pageNumber, _ := strconv.ParseInt(req.FormValue("Page"),10,64)
		CORE.ServeTemplateWithParams(res, "private-message", struct {
			HeaderData	CONTEXT.HeaderData
			ErrorResponse, RedirectURL string
			PageNumber int64
			Messages []PM.PrivateMessage
		}{
			HeaderData:    *MakeHeader(ctx),
			RedirectURL:   req.FormValue("redirect"),
			ErrorResponse: req.FormValue("ErrorResponse"),
			PageNumber:   pageNumber,
			Messages: 	   PM.RetrieveMessages(ctx,3,int(pageNumber)),
		})
	})
	r.POST("/pm", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		ctx := CONTEXT.NewContext(res,req)
		PM.SendMessage(ctx,req.FormValue("Receiver"),req.FormValue("Title"),req.FormValue("Body"))
		ctx.Redirect("/pm")
	})
}