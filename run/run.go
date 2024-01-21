package run

import (
	"context"
	"errors"
	"fmt"
	"log"

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
		if err != nil {
			return nil, err
		}
		fmt.Println(" ID:", color.BlueString(cv.ID))

		err = client.ConfigurationVersions.Upload(ctx, cv.UploadURL, tfDirectory)
		if err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Using existing Config Version ...", cvID)
		cv, err = client.ConfigurationVersions.Read(ctx, cvID)
		if err != nil {
			return nil, err
		}
		fmt.Println(" ID:", color.BlueString(cv.ID), " Status: ", color.BlueString(string(cv.Status)))
		if cv.Status != "uploaded" {
			return nil, errors.New("provider configuration version is not allowed")
		}
	}

	return cv, nil
}

func CreateOrDestroy(client tfe.Client, workspaceId, message string, dir string, autoApply bool, isDestroy bool) *tfe.Run {
	ctx := context.Background()
	ws, err := client.Workspaces.ReadByID(ctx, workspaceId)

	if err != nil {
		log.Fatal(err)
	}

	cv, err := createOrReadConfigurationVersion(ctx, &client, workspaceId, "", dir, false)

	if err != nil {
		log.Fatal(err)
	}

	result, err := client.Runs.Create(ctx, tfe.RunCreateOptions{
		AutoApply:            tfe.Bool(autoApply),
		Workspace:            ws,
		ConfigurationVersion: cv,
		IsDestroy:            tfe.Bool(isDestroy),
		Message:              tfe.String(message)})

	if err != nil {
		log.Fatal(err)
	}
	return result

}

func listRun(client tfe.Client, workspaceId string) *tfe.RunList {
	ctx := context.Background()

	result, err := client.Runs.List(ctx, workspaceId, nil)
	if err != nil {
		log.Fatal(err)
	}
	return result

}
func List(client *tfe.Client, workspaceName, org string) {
	ws := workspace.GetByName(*client, org, workspaceName)
	list := listRun(*client, ws.ID)
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "STATUS", "CREATED AT", "MESSAGE", "SOURCE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Status, v.CreatedAt, v.Message, v.Source)
	}

	tbl.Print()
}
