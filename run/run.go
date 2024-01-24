package run

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/costa86/tformer/helper"
	"github.com/costa86/tformer/workspace"
	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func createOrReadConfigurationVersion(ctx context.Context, client *tfe.Client, workspaceID string, cvID string, tfDirectory string, speculative bool) (*tfe.ConfigurationVersion, error) {
	var err error
	var cv *tfe.ConfigurationVersion

	if cvID == "" {
		fmt.Print("Creating new Config Version ...")
		cv, err = client.ConfigurationVersions.Create(ctx, workspaceID, tfe.ConfigurationVersionCreateOptions{
			AutoQueueRuns: tfe.Bool(false),
			Speculative:   tfe.Bool(speculative),
		})
		helper.HandleError(err)
		fmt.Println(" ID:", color.BlueString(cv.ID))
		err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, tfDirectory)
		helper.HandleError(err)

	} else {
		fmt.Println("Using existing Config Version:", color.BlueString(cvID))
		cv, err = client.ConfigurationVersions.Read(ctx, cvID)
		helper.HandleError(err)
		fmt.Println(" ID:", color.BlueString(cv.ID), " Status: ", color.BlueString(string(cv.Status)))

		if cv.Status != "uploaded" {
			return nil, errors.New("provider configuration version is not allowed")
		}
	}

	return cv, nil
}

func CreateOrDestroy(client tfe.Client, workspaceId, message string, dir string, autoApply bool, isDestroy bool, cvId string) *tfe.Run {
	ctx := context.Background()
	ws, err := client.Workspaces.ReadByID(ctx, workspaceId)
	helper.HandleError(err)

	var configurationId string

	if cvId == "latest" {
		run := getLatest(&client, ws.ID)
		configurationId = run.ConfigurationVersion.ID
		fmt.Println("Using Configuration Version from the latest Run: ", color.BlueString(string(run.ID)))

	} else {
		configurationId = cvId
	}

	cv, err := createOrReadConfigurationVersion(ctx, &client, workspaceId, configurationId, dir, false)
	helper.HandleError(err)

	result, err := client.Runs.Create(ctx, tfe.RunCreateOptions{
		AutoApply:            tfe.Bool(autoApply),
		Workspace:            ws,
		ConfigurationVersion: cv,
		IsDestroy:            tfe.Bool(isDestroy),
		Message:              tfe.String(message)})

	helper.HandleError(err)
	return result

}

func Get(client tfe.Client, id string) {
	ctx := context.Background()
	result, err := client.Runs.Read(ctx, id)
	helper.HandleError(err)
	userJSON, err := json.MarshalIndent(result, "", "    ")
	helper.HandleError(err)
	log.Printf("%s", userJSON)
}

func listRun(client tfe.Client, workspaceId string) *tfe.RunList {
	ctx := context.Background()

	result, err := client.Runs.List(ctx, workspaceId, nil)
	helper.HandleError(err)

	return result

}
func List(client *tfe.Client, workspaceName, org string) {
	ws := workspace.GetByName(*client, org, workspaceName)
	list := listRun(*client, ws.ID)
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "STATUS", "CREATED AT", "MESSAGE", "SOURCE", "CV ID")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Status, v.CreatedAt, v.Message, v.Source, v.ConfigurationVersion.ID)
	}

	tbl.Print()
}

func getLatest(client *tfe.Client, wsId string) tfe.Run {
	list := listRun(*client, wsId)

	var run tfe.Run
	var latestCreatedAt time.Time

	for _, v := range list.Items {
		if v.CreatedAt.After(latestCreatedAt) {
			latestCreatedAt = v.CreatedAt
			run = *v
		}
	}
	return run

}
