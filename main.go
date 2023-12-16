package main

/*

This application has the follow endpoints:
* GET /        - Show home page with hostname and IP address of host/container.
* GET /css     - Show the CSS template.
* GET /health  - Show status of application.
* GET /metrics - Show the metrics in format supported by Prometheus.

This project will be comprised of two Go files:
* main.go - Contains a parent router, which sub-routers mount to and its handlers.

References:
* https://gobyexample.com/http-server
* https://github.com/cod3rcursos/curso-go/tree/master/http
* https://github.com/weaveworks/scope/blob/v1.13.2/common/hostname/hostname.go
* https://github.com/richardpct/go-hostname/blob/master/hostname02/hostname.go
* https://adlerhsieh.com/blog/rendering-dynamic-data-in-go-http-template
* https://github.com/paulbouwer/hello-kubernetes
* https://golangdocs.com/templates-in-golang
* https://hackthedeveloper.com/golang-html-template-parsefiles-and-execute/
* https://blog.gopheracademy.com/advent-2017/using-go-templates/
* https://philipptanlak.com/mastering-html-templates-in-go-the-fundamentals/
* https://pkg.go.dev/html/template#example-Template-Parsefiles
* https://stackoverflow.com/questions/61057271/how-to-run-html-with-css-using-golang
* https://stackoverflow.com/questions/28793619/golang-what-to-use-http-servefile-or-http-fileserver/28798174#28798174
* https://www.tumblr.com/golang-examples/99458329439/get-local-ip-addresses
* https://pkg.go.dev/embed
* https://charly3pins.dev/blog/learn-how-to-use-the-embed-package-in-go-by-building-a-web-page-easily/
* https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/
* https://blog.pvincent.io/2017/12/prometheus-blog-series-part-1-metrics-and-labels/
* https://medium.com/@matryer/the-http-handler-wrapper-technique-in-golang-updated-bc7fbcffa702
* https://github.com/gin-gonic/gin
* https://www.tabnine.com/blog/golang-gin/
* https://blog.logrocket.com/gin-binding-in-go-a-tutorial-with-examples/
* https://semaphoreci.com/community/tutorials/building-go-web-applications-and-microservices-using-gin
* https://blog.logrocket.com/rest-api-golang-gin-gorm/
* https://www.twilio.com/blog/build-restful-api-using-golang-and-gin
* https://chenyitian.gitbooks.io/gin-tutorials/content/gin/8.html
* https://hoohoo.top/blog/20210530112304-golang-tutorial-introduction-gin-html-template-and-how-integration-with-bootstrap/
* https://b-nova.com/en/home/content/fully-featured-golang-with-gin-web-framework
* https://github.com/penglongli/gin-metrics
* https://blog.petehouston.com/monitor-gin-gonic-application-with-prometheus-metrics/
* https://dev.to/kishanbsh/capturing-custom-last-request-time-metrics-using-prometheus-in-gin-36d6

*/

import (
	"embed"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

const (
	htmlExtension = "templates/*.html"
)

var (
	// Embed static content for using with web server.
	//
	//go:embed templates/*
	filesEmbed embed.FS

	returnCodeError bool   = false
	cssDir          string = "css"
	protocol               = "http"
	address                = "0.0.0.0"
	port                   = "3000"
	ginMode                = "release"
)

// Show status of application
func status(context *gin.Context) {

	context.JSON(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Pass the data that the page uses (in this case, 'title')
		gin.H{
			"message": "It' works",
		},
	)
}

// Get hostname
func hostname() (name string) {

	name, err := os.Hostname()

	if err != nil {
		log.Fatal(err)
	}

	return
}

// Get interface name and IPAddress
func getAddress() (address map[string]string) {

	address = make(map[string]string)

	listIfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range listIfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		if addrs != nil {
			address[iface.Name] = addrs[0].String()
		}
	}

	return
}

// Parse of HTML content
func parseTemplate(context *gin.Context) {
	// Call the HTML method of the Context to render a template
	context.HTML(
		// Set the HTTP status to 200 (OK)
		http.StatusOK,
		// Use the index.html template
		"index.html",
		// Pass the data that the page uses
		gin.H{
			"Title":     "kube-pires",
			"Hostname":  hostname(),
			"Addresses": getAddress(),
		},
	)
}

func main() {
	// Reading environment variables
	if protocolEnv := os.Getenv("PROTOCOL"); protocolEnv != "" {
		protocol = strings.ToLower(protocolEnv)
	}

	if addressEnv := os.Getenv("ADDRESS"); addressEnv != "" {
		address = strings.ToLower(addressEnv)
	}

	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port = portEnv
	}

	if ginModeEnv := os.Getenv("GIN_MODE"); ginModeEnv != "" {
		// Gin run in "debug" mode by default. Can be switch to "release" mode in production.
		// using env: export GIN_MODE=release
		// or
		// using code: gin.SetMode(gin.ReleaseMode)
		//gin.SetMode(gin.ReleaseMode)
		ginMode = ginModeEnv
	}
	// Setting ginMode
	gin.SetMode(ginMode)

	// Defining listen and baseurl
	listen := address + ":" + port
	baseurl := protocol + "://" + address + ":" + port

	log.Printf("Starting up on %s", baseurl)

	// Set the router as the default one provided by Gin
	// Default With the Logger and Recovery middleware already attached
	router := gin.Default()

	// Get global Monitor object
	metrics := ginmetrics.GetMonitor()

	// Set metric path, default /debug/metrics
	metrics.SetMetricPath("/metrics")

	// Set slow time, default 5s
	metrics.SetSlowTime(10)

	// Set request duration, default {0.1, 0.3, 1.2, 5, 10}
	// used to p95, p99
	metrics.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})

	// Set middleware for gin
	metrics.Use(router)

	// HTTP Handlers
	// Exposing CSS files without embed directive
	router.StaticFS("/css", http.Dir(cssDir))
	router.StaticFile("/favicon.ico", "./images/favicon.ico")

	// Process the templates at the start so that they don't have to be loaded
	// from the disk again. This makes serving HTML pages very fast.
	router.LoadHTMLGlob("templates/*")

	// Define the route for the index page and display the index.html template
	// To start with, we'll use an inline route handler. Later on, we'll create
	// standalone functions that will be used as route handlers.
	router.GET("/", parseTemplate)
	router.GET("/health", status)

	// Init webserver
	errServer := router.Run(listen)
	if errServer != nil {
		log.Fatal("[ERROR] Starting the HTTP Server :", errServer)
		return
	}
}
