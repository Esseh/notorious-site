package main
import(
	"net/http"
	"fmt"
	"github.com/Esseh/notorious-dev/CONTEXT"
	"github.com/Esseh/notorious-dev/NOTIFICATION"
	"github.com/julienschmidt/httprouter"
)
func INIT_NOTIFICATION_API(r *httprouter.Router) {
	r.POST("/NOTIFICATION/api/clear", func(res http.ResponseWriter, req *http.Request, params httprouter.Params){
		fmt.Fprint(res,NOTIFICATION.ClearNotificationsAPI(CONTEXT.NewContext(res,req)))
	})
}