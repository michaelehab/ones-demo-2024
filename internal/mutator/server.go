package mutator

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/project-alvarium/alvarium-sdk-go/pkg/interfaces"
	"github.com/project-alvarium/ones-demo-2024/internal/config"
	"github.com/project-alvarium/ones-demo-2024/internal/db"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// HttpServer contains references to dependencies required by the http server implementation.
type HttpServer struct {
	config    config.EndpointInfo
	dbMongo   *db.MongoProvider
	logger    interfaces.Logger
	router    *mux.Router
	chPublish chan []byte
}

// NewHttpServer is a factory method that returns an initialized HttpServer receiver struct.
func NewHttpServer(router *mux.Router, pub chan []byte, config config.EndpointInfo, dbMongo *db.MongoProvider, logger interfaces.Logger) *HttpServer {
	return &HttpServer{
		config:    config,
		dbMongo:   dbMongo,
		logger:    logger,
		router:    router,
		chPublish: pub,
	}
}

// BootstrapHandler fulfills the BootstrapHandler contract.  It creates two go routines -- one that executes ListenAndServe()
// and another that waits on closure of a context's done channel before calling Shutdown() to cleanly shut down the
// http server.
func (b *HttpServer) BootstrapHandler(
	ctx context.Context,
	wg *sync.WaitGroup) bool {

	// this allows env override to explicitly set the value used
	// for ListenAndServe as needed for different deployments
	addr := ":" + strconv.Itoa(b.config.Service.Port)

	timeout := time.Millisecond * 10000
	server := &http.Server{
		Addr:         addr,
		Handler:      b.router,
		WriteTimeout: timeout,
		ReadTimeout:  timeout,
	}

	b.logger.Write(slog.LevelInfo, "Web server starting ("+addr+")")

	wg.Add(1)
	go func() {
		defer wg.Done()

		if len(b.config.Certificate) > 0 {
			_ = server.ListenAndServeTLS(b.config.Certificate, b.config.Key)
		} else {
			_ = server.ListenAndServe()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-ctx.Done()
		//DEBUG
		b.logger.Write(slog.LevelInfo, "Web server shutting down")
		_ = server.Shutdown(ctx)
		b.dbMongo.Close(ctx)
		close(b.chPublish)
		//DEBUG
		b.logger.Write(slog.LevelInfo, "Web server shut down")
	}()

	return true
}
