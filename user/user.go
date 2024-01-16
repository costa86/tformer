package user

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-tfe"
)

func ReadCurrent(client *tfe.Client) {

	ctx := context.Background()

	user, err := client.Users.ReadCurrent(ctx)
	if err != nil {
		log.Fatal(err)
	}
	userJSON, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", userJSON)

}
