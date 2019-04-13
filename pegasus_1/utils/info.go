package utils

// define types that can be represented as entries of database
// like user, item ...

type User struct {
	UserID string `json:"user_id"` // why Username should be uppercase ?
	Password string `json:"password"`
	Username string `json:"username"`

}

type Item struct {

}

type Order struct {
	UserID     string  `json:"user_id"`
	Size       string  `json:"size"`
	Weight     float64 `json:"weight"`
	Arrival    string  `json:"arrival"`
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
