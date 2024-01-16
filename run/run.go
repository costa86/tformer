package run

import (
	"context"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func Create(client tfe.Client, workspaceId string, targetDir []string) *tfe.Run {
	ctx := context.Background()
	ws, err := client.Workspaces.ReadByID(ctx, workspaceId)

	if err != nil {
		fmt.Println("ff")
		log.Fatal(err)
	}
	fmt.Println(ws.CreatedAt)

	result, err := client.Runs.Create(ctx, tfe.RunCreateOptions{
		TargetAddrs:          targetDir,
		Workspace:            ws,
		ConfigurationVersion: &tfe.ConfigurationVersion{AutoQueueRuns: true},
		IsDestroy:            tfe.Bool(false),
		Message:              tfe.String("hello")})

	if err != nil {
		fmt.Println("pp")

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
func List(client *tfe.Client, workspaceId string) {
	list := listRun(*client, workspaceId)
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "STATUS", "CREATED AT", "MESSAGE", "SOURCE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Status, v.CreatedAt, v.Message, v.Source)
	}

	tbl.Print()
}
