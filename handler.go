package webhook

import (
	"encoding/json"
	"fmt"
	"github.com/aiteung/atapi"
	"github.com/aiteung/atmessage"
	"github.com/aiteung/module/model"
	"github.com/whatsauth/wa"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

func PostBalasan(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)
	link := "https://medium.com/@sidiqfaisal30/apa-itu-cybser-security-8993e4605cc3"
	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		if msg.Message == "loc" || msg.Message == "Loc" || msg.Message == "lokasi" || msg.LiveLoc {
			location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
			if err != nil {
				// Handle the error (e.g., log it) and set a default location name
				location = "Unknown Location"
			}

			reply := fmt.Sprintf("hallo akang/eteh faibit tebak pasti lagi di %s \n Koordinatenya : %s - %s\n Cara Penggunaan WhatsAuth Ada di link dibawah ini"+
				"yaa kak %s\n", location,
				strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)), link)
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: reply,
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		} else {
			randm := []string{
				"Hallo Akang/eteh punteun " + msg.Alias_name + "\nfaisal lagi gaadaa \n aku faibit salam kenall yaaaa \n Cara penggunaan WhatsAuth ada di link berikut ini ya kak...\n" + link,
				"IHHH jangan SPAAM berisik tau faibit lagi tidur",
				"Kamu tampan tau",
				"masyaalllah kamu cantik banget",
				"bray, lebih baik beliin aku nasgor",
				"Jangan galak galak dong kak, aku catut tauu",
				"orang lain sudah jauh ini masih diem bae, MAKANYA BANGUN PAGI KAK",
				"bro hari ini faibit senang banget",
				"beyanda beyanda",
			}
			dt := &wa.TextMessage{
				To:       msg.Phone_number,
				IsGroup:  false,
				Messages: GetRandomString(randm),
			}
			resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://api.wa.my.id/api/send/message/text")
		}
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}

func ReverseGeocode(latitude, longitude float64) (string, error) {
	// OSM Nominatim API endpoint
	apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?format=json&lat=%f&lon=%f", latitude, longitude)

	// Make a GET request to the API
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Decode the response
	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	// Extract the place name from the response
	displayName, ok := result["display_name"].(string)
	if !ok {
		return "", fmt.Errorf("unable to extract display_name from the API response")
	}

	return displayName, nil
}

func Liveloc(w http.ResponseWriter, r *http.Request) {
	var msg model.IteungMessage
	var resp atmessage.Response
	json.NewDecoder(r.Body).Decode(&msg)

	// Reverse geocode to get the place name
	location, err := ReverseGeocode(msg.Latitude, msg.Longitude)
	if err != nil {
		// Handle the error (e.g., log it) and set a default location name
		location = "Unknown Location"
	}

	reply := fmt.Sprintf("kamu lagi kamu lagi kenapa selalu ada kamu sih. aku tebak pasti kamu lagi di %s \n Koordinatenya : %s - %s\n", location,
		strconv.Itoa(int(msg.Longitude)), strconv.Itoa(int(msg.Latitude)))

	if r.Header.Get("Secret") == os.Getenv("SECRET") {
		dt := &wa.TextMessage{
			To:       msg.Phone_number,
			IsGroup:  false,
			Messages: reply,
		}
		resp, _ = atapi.PostStructWithToken[atmessage.Response]("Token", os.Getenv("TOKEN"), dt, "https://cloud.wa.my.id/api/send/message/text")
	} else {
		resp.Response = "Secret Salah"
	}
	fmt.Fprintf(w, resp.Response)
}

func GetRandomString(strings []string) string {
	randomIndex := rand.Intn(len(strings))
	return strings[randomIndex]
}
