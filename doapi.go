package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/unrolled/render"
)

var (
	rdr *render.Render
)

func getCancelledClasses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var katei, kouchi int

	location := r.URL.Query().Get("location")

	switch location {
	case "1":
		katei = 1
		kouchi = 1
	case "2":
		katei = 1
		kouchi = 2
	case "3":
		katei = 3
		kouchi = 3
	default:
		rdr.JSON(w, http.StatusBadRequest, &Error{
			Message: fmt.Sprintf("Unknown location \"%s\". The location parameter should be 1, 2, or 3.", location),
		})
		return
	}

	var notifications []Notification

	for youbi := 1; youbi <= 5; youbi++ {
		body, err := Get(katei, kouchi, youbi)
		if err != nil {
			log.Println(err)
			continue
		}

		n, err := Parse(body)

		if err != nil {
			log.Println(err)
			continue
		}

		notifications = append(notifications, *n)
	}

	if len(notifications) != 5 {
		rdr.JSON(w, http.StatusInternalServerError, notifications)
	} else {
		w.Header().Set("Cache-Control", "max-age:1800, public")
		w.Header().Set("Expires", time.Now().Add(30*time.Hour).Format(http.TimeFormat))
		rdr.JSON(w, http.StatusOK, notifications)
	}
}

func main() {
	// Render
	rdr = render.New()

	// Routing
	router := httprouter.New()
	router.GET("/v1/duet/cancelled_classes", getCancelledClasses)

	//CORS
	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
