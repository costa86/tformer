package configversion

import (
	"context"
	"encoding/json"
	"log"

	"github.com/hashicorp/go-tfe"
)

func List(client *tfe.Client, workspaceId string) {

	ctx := context.Background()

	res, err := client.ConfigurationVersions.List(ctx, workspaceId, &tfe.ConfigurationVersionListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	resJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", resJSON)

}

func Download(client *tfe.Client, cvId string) {

	ctx := context.Background()

	res, err := client.ConfigurationVersions.Download(ctx, cvId)
	if err != nil {
		log.Fatal(err)
	}
	resJSON, err := json.MarshalIndent(res, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", resJSON)

}
