package utils

import (
	"log"
	"time"
)

func GetISTTimestamp() string {
	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Println("Timezone error:", err)
		return time.Now().In(time.FixedZone("IST", 5*60*60+30*60)).
			Format(time.RFC3339)
	}

	return time.Now().In(ist).Format(time.RFC3339)
}
