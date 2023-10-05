package transport

import (
	"log"
	"net/http"
	"os"

	"github.com/FadyGamilM/go-websockets/configs"
	"github.com/gin-gonic/gin"
)

func CreateRouter() *gin.Engine {
	return gin.Default()
}

type RestServer struct {
	*http.Server
}

func CreateServer(r *gin.Engine) *RestServer {
	configs, err := configs.LoadServerConfigs("./configs")
	if err != nil {
		return nil
	}
	return &RestServer{
		Server: &http.Server{
			Addr:    configs.Server.Port,
			Handler: r,
		},
	}
}

func (srv *RestServer) Run() {
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("cannot start the server on port : %v \n the error is : %v \n", srv.Addr, err)
		os.Exit(1)
	}
}
