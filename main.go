package main

import (
	"flag"
	"log"
	"net/http"
	"os" // Import the os package

	"github.com/jonesrussell/sober/content"
	"github.com/jonesrussell/sober/handlers"
	"github.com/jonesrussell/sober/server"
	"github.com/jonesrussell/sober/services"

	"github.com/labstack/echo/v4"
)

func main() {
	// Define and parse the "generate" flag
	generate := flag.Bool("generate", false, "Generate static site")
	flag.Parse()

	// Get the base path from the environment variable
	basePath := os.Getenv("BASE_PATH")

	pageService := services.NewPageService()
	staticSiteService := services.NewStaticSiteService(pageService)

	if *generate {
		// Generate static site
		staticSiteService.Generate(basePath) // Pass the base path to the Generate method
		return
	}

	// Start Echo server
	srv := server.NewServer(pageService)

	// Set the NotFoundHandler to render your custom 404 page
	srv.Echo.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		if code == http.StatusNotFound {
			// Use your DefaultHandler to render the 404 page
			handler := &handlers.DefaultHandler{
				PageService: pageService,
				BasePath:    basePath,
			}
			if err := handler.RenderPage(c, content.NotFound(), "404"); err != nil {
				c.Logger().Error(err)
			}
		} else {
			srv.Echo.DefaultHTTPErrorHandler(err, c)
		}
	}

	// Start server
	log.Fatal(srv.Echo.Start(":8080"))
}
