package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func Lock(client tfe.Client, id string) {
	ctx := context.Background()

	result, err := client.Workspaces.Lock(ctx, id, tfe.WorkspaceLockOptions{})

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Workspace locked %s by %s", result.Name, result.LockedBy.User.Email)

}

func listWorkspaces(client tfe.Client, orgName string) *tfe.WorkspaceList {
	ctx := context.Background()

	result, err := client.Workspaces.List(ctx, orgName, nil)
	if err != nil {
		log.Fatal(err)
	}
	return result

}
func List(client *tfe.Client, name string, asJson bool) {
	list := listWorkspaces(*client, name)

	if asJson {
		userJSON, err := json.MarshalIndent(list.Items, "", "    ")
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%s", userJSON)
		return
	}

	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "NAME", "TERRAFORM VERSION")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Name, v.TerraformVersion)
	}

	tbl.Print()
}

func Create(client tfe.Client, organization, name string) {
	ctx := context.Background()

	result, err := client.Workspaces.Create(ctx, organization, tfe.WorkspaceCreateOptions{Name: tfe.String(name)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Workspace created: %s", result.Name)

}

func Delete(client tfe.Client, id string) {
	ctx := context.Background()

	err := client.Workspaces.DeleteByID(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Workspace deleted: %s", id)

}

func Purge(client tfe.Client, orgName string) {
	ctx := context.Background()
	workspaces := listWorkspaces(client, orgName)
	for _, v := range workspaces.Items {
		err := client.Workspaces.DeleteByID(ctx, v.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("All Workspace deleted")

}
