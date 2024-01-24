package configversion

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/costa86/tformer/helper"
	"github.com/costa86/tformer/workspace"
	"github.com/fatih/color"
	"github.com/hashicorp/go-slug"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func Get(client *tfe.Client, id string) {
	ctx := context.Background()
	cv, err := client.ConfigurationVersions.Read(ctx, id)
	helper.HandleError(err)
	cvJSON, err := json.MarshalIndent(cv, "", "    ")
	helper.HandleError(err)
	log.Printf("%s", cvJSON)
}

func List(client *tfe.Client, workspaceName, org string) {

	ctx := context.Background()
	ws := workspace.GetByName(*client, org, workspaceName)

	res, err := client.ConfigurationVersions.List(ctx, ws.ID, &tfe.ConfigurationVersionListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "URL", "STATUS", "SOURCE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range res.Items {
		tbl.AddRow(v.ID, v.UploadURL, v.Status, v.Source)
	}

	tbl.Print()
}

func Download(client *tfe.Client, cvId string) {
	ctx := context.Background()
	cv, err := client.ConfigurationVersions.Download(ctx, cvId)
	helper.HandleError(err)
	reader := bytes.NewReader(cv)
	e := slug.Unpack(reader, cvId)
	helper.HandleError(e)
}
