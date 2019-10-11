package server

import (
	"cinemas-microservices/booking-service/src/api"
	"cinemas-microservices/booking-service/src/routes"
	tracing "cinemas-microservices/booking-service/src/tracing"
	"context"
	"os"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

var e *echo.Echo

// Start ...
func Start(r map[string]interface{}, se chan error) {

	// get server settings from dependecy injection
	ss := r["ss"].(map[string]interface{})

	e = echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.Recover())

	tracer, closer := tracing.InitJaeger("booking-service", ss["tracingUpstream"].(string))
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	e.Use(tracing.TraceWithConfig(tracing.TraceConfig{
		Tracer:  tracer,
		Skipper: nil,
	}))

	app := e.Group("/booking")

	routes.API(app, r["repo"].(api.Repository))
	routes.HealthyAPI(e)

	// Start server
	go func() {
		if err := e.Start(":" + strconv.Itoa(ss["port"].(int))); err != nil {
			log.Info("shutting down the server")
			se <- err
		}
	}()
}

// Shutdown ...
func Shutdown(s *mgo.Session) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	s.Close()
	log.Warn("Server shutdown")
	os.Exit(1)
}

// GetServer ...
func GetServer() *echo.Echo {
	return e
}
