package main

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
)


func INIT_FORUM_HANDLERS(r *httprouter.Router) {
	r.GET("/forum/category/view",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})
	r.GET("/forum/category/make",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})
	r.POST("/forum/category/make",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})
	r.GET("/forum/forums/view/:ForumID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})
	r.GET("/forum/forums/make/:CategoryID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})
	r.POST("/forum/forums/make/:CategoryID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})
	r.GET("/forum/thread/view/:ThreadID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})	
	r.GET("/forum/thread/make/:ForumID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})	
	r.POST("/forum/thread/make/:ForumID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})	
	r.GET("/forum/post/make/:ThreadID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})	
	r.POST("/forum/post/make/:ThreadID",func(res http.ResponseWriter,req *http.Request, p httprouter.Params){})	
}