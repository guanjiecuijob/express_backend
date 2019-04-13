package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kr/pretty"
	"googlemaps.github.io/maps"
	"log"
	"math"
	"net/http"

)

type Order struct {
	//UserID     string  `json:"user_id"`
	Size       string  `json:"size"`
	Weight     float64 `json:"weight"`
	//Arrival    string  `json:"arrival"`
	PickupLoc  string  `json:"pickup"`
	DropoffLoc string  `json:"dropoff"`
}

type DeliverMethod struct {
	RobotTime     float64 `json:"robot_time"`
	RobotDistance float64 `json:"robot_distance"`
	RobotPrice    float64 `json:"robot_price"`
	DroneTime     float64 `json:"drone_time"`
	DroneDistance float64 `json:"drone_distance"`
	DronePrice    float64 `json:"drone_price"`
	FormatPickup  string  `json:"pickup"`
	FormatDropoff string  `json:"dropoff"`
}

func SearchPath(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received one order request")
	w.Header().Set("Content-Type", "text/palin")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == "OPTIONS" {
		return
	}

	decoder := json.NewDecoder(r.Body)
	var order Order
	// handle err by informal user type
	if err := decoder.Decode(&order); err != nil {
		http.Error(w, "cannot decode order data from client", http.StatusBadRequest)
		fmt.Printf("cannot decode order data from client %v.\n", err)
		return
	}

	deliverMethod, err := processOrder(order)
	if err != nil {
		http.Error(w, "Failed to processOrder using googleMapAPI", http.StatusInternalServerError)
		fmt.Printf("Failed to processOrder using googleMapAPI %v.\n", err)
		return
	}

	js, err := json.Marshal(deliverMethod)
	if err != nil {
		http.Error(w, "Failed to parse posts into JSON format", http.StatusInternalServerError)
		fmt.Printf("Failed to parse posts into JSON format %v.\n", err)
		return
	}

	w.Write(js)
}

func processOrder(order Order) (DeliverMethod, error) {

	c, err := maps.NewClient(maps.WithAPIKey("googleMapApi")) // 这里的密钥要改
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r1 := &maps.DistanceMatrixRequest{
		Origins:      []string{order.PickupLoc},
		Destinations: []string{order.DropoffLoc},
	}

	distanceMatrixResponse, err := c.DistanceMatrix(context.Background(), r1)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	var FormatPickup string = distanceMatrixResponse.OriginAddresses[0]
	var FormatDropoff string = distanceMatrixResponse.DestinationAddresses[0]

	var RobotDistance float64 = float64(distanceMatrixResponse.Rows[0].Elements[0].Distance.Meters) / 1000 // km
	var RobotTime float64 = float64(distanceMatrixResponse.Rows[0].Elements[0].Duration) / 1000000000 / 60 //  minute
	var RobotPrice float64 = float64(RobotDistance * 0.1 * order.Weight)                                   // dollar (assume $0.1 per km per kg)

	r2 := &maps.GeocodingRequest{
		Address: order.PickupLoc,
	}

	pickupLatLong, err := c.Geocode(context.Background(), r2)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	r3 := &maps.GeocodingRequest{
		Address: order.DropoffLoc,
	}

	dropoffLatLong, err := c.Geocode(context.Background(), r3)
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}
	var DroneDistance float64 = float64(straightDistance(pickupLatLong[0].Geometry.Location.Lat, pickupLatLong[0].Geometry.Location.Lng, dropoffLatLong[0].Geometry.Location.Lat, dropoffLatLong[0].Geometry.Location.Lng)) // km
	var DroneTime float64 = float64(DroneDistance / 200 * 60)                                                                                                                                                               // minute (assume drone speed 200 km/hour)
	var DronePrice float64 = float64(DroneDistance * 0.3 * order.Weight)                                                                                                                                                    // dollar (assume $0.3 per km per kg)

	deliverMethod := DeliverMethod{RobotTime, RobotDistance, RobotPrice, DroneTime, DroneDistance, DronePrice, FormatPickup, FormatDropoff}

	pretty.Println(RobotTime, RobotDistance, RobotPrice, DroneTime, DroneDistance, DronePrice)
	return deliverMethod, nil
}

func straightDistance(lat1, lon1, lat2, lon2 float64) float64 {
	var p float64 = 0.017453292519943295 // Math.PI / 180
	var a float64 = 0.5 - math.Cos((lat2-lat1)*p)/2 +
		math.Cos(lat1*p)*math.Cos(lat2*p)*
			(1-math.Cos((lon2-lon1)*p))/2

	return 12742 * math.Sin(math.Sqrt(a)) // 2 * R; R = 6371 km
}
