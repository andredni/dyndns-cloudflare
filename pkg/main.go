package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/gorilla/schema"
)

func main() {
	configFile := flag.String("file", "", "a string")
	flag.Parse()

	if len(*configFile) > 0 {
		loadConfig(*configFile)
	}

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}

type Message struct {
	Message string
}

type Account struct {
	Password string `schema:"password,required"`
	Username string `schema:"username"`
	Ipv4     string `schema:"ipv4,required"`
	Hostname string `schema:"hostname,required"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	params := Account{}
	decoder := schema.NewDecoder()
	err := decoder.Decode(&params, r.URL.Query())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(fatal(err.Error()))
		return
	}

	myRecords := append(config.Records, params.Hostname)
	uniqueRecords := unique(myRecords)

	api, err := cloudflare.NewWithAPIToken(params.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(fatal(err.Error()))
		return
	}

	ctx := context.Background()
	zones, err := api.ListZones(ctx)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(fatal(err.Error()))
		return
	}

	for _, myrecord := range uniqueRecords {
		for _, zone := range zones {

			match, _ := regexp.MatchString(zone.Name, myrecord)
			if match {

				records, err := api.DNSRecords(ctx, zone.ID, cloudflare.DNSRecord{})
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(fatal(err.Error()))
					return
				}

				for _, record := range records {
					match, _ := regexp.Match(myrecord, []byte(record.Name))
					if match {
						err := api.DeleteDNSRecord(ctx, zone.ID, record.ID)
						if err != nil {
							w.WriteHeader(http.StatusBadRequest)
							w.Write(fatal(err.Error()))
							return
						}
					}
				}

				new_record := cloudflare.DNSRecord{Name: myrecord, Content: params.Ipv4, Type: "A"}
				_, err = api.CreateDNSRecord(ctx, zone.ID, new_record)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(fatal(err.Error()))
					return
				}
			}
		}
	}

	message := Message{Message: "Records " + strings.Join(uniqueRecords, ", ") + " was successful upated with IP " + params.Ipv4}
	log.Println(message.Message)
	jsonResp, _ := json.Marshal(message)
	w.Write(jsonResp)
}
