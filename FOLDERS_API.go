package main
import(
	"net/http"
	"fmt"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/FOLDERS"
	"github.com/julienschmidt/httprouter"
)
func INIT_FOLDERS_API(r *httprouter.Router) {
	r.POST("/folder/api/newfolder", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.NewFolder(CONTEXT.NewContext(res,req)))
	})
	r.POST("/folder/api/deletefolder", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.DeleteFolder(CONTEXT.NewContext(res,req)))
	})
	r.POST("/folder/api/openfolder", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.OpenFolder(CONTEXT.NewContext(res,req)))
	})
	r.POST("/folder/api/addnote", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.AddNote(CONTEXT.NewContext(res,req)))
	})
	r.POST("/folder/api/removenote", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.RemoveNote(CONTEXT.NewContext(res,req)))
	})
	r.POST("/folder/api/initializeroot", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.InitializeRoot(CONTEXT.NewContext(res,req)))
	})
	r.POST("/folder/api/renamefolder", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,FOLDERS.RenameFolder(CONTEXT.NewContext(res,req)))
	})
}