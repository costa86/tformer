package workspace

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/costa86/tformer/helper"

	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func GetByName(client tfe.Client, org, name string) *tfe.Workspace {
	ws, err := client.Workspaces.Read(context.Background(), org, name)
	helper.HandleError(err)
	return ws

}

func LockOrUnlock(client tfe.Client, name, org string, lock bool) {
	ctx := context.Background()

	ws := GetByName(client, org, name)
	var result *tfe.Workspace
	var err error
	var action = "locked"

	if lock {
		result, err = client.Workspaces.Lock(ctx, ws.ID, tfe.WorkspaceLockOptions{})
		helper.HandleError(err)
	} else {
		result, err = client.Workspaces.Unlock(ctx, ws.ID)
		helper.HandleError(err)
		action = "unlocked"
	}

	fmt.Printf("Workspace %s %s ", result.Name, action)

}
func LockOrUnlockAll(client tfe.Client, org string, lock bool) {
	ctx := context.Background()

	list := listWorkspaces(client, org)
	var action = "locked"

	if lock {
		for _, v := range list.Items {
			res, err := client.Workspaces.Lock(ctx, v.ID, tfe.WorkspaceLockOptions{})
			helper.HandleError(err)
			fmt.Printf("Workspace %s locked ", res.Name)

		}
	} else {
		action = "unlocked"
		for _, v := range list.Items {
			res, err := client.Workspaces.Unlock(ctx, v.ID)
			helper.HandleError(err)
			fmt.Printf("Workspace %s unlocked ", res.Name)

		}
	}

	fmt.Printf("%d workspaces %s ", len(list.Items), action)

}

func listWorkspaces(client tfe.Client, orgName string) *tfe.WorkspaceList {
	ctx := context.Background()

	result, err := client.Workspaces.List(ctx, orgName, nil)
	helper.HandleError(err)
	return result

}
func List(client *tfe.Client, name string, asJson bool) {
	list := listWorkspaces(*client, name)

	if asJson {
		userJSON, err := json.MarshalIndent(list.Items, "", "    ")
		helper.HandleError(err)

		log.Printf("%s", userJSON)
		return
	}

	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "NAME", "TERRAFORM VERSION", "LOCKED")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Name, v.TerraformVersion, v.Locked)
	}

	tbl.Print()
}

func Create(client tfe.Client, organization, name string) {
	ctx := context.Background()

	result, err := client.Workspaces.Create(ctx, organization, tfe.WorkspaceCreateOptions{Name: tfe.String(name)})
	helper.HandleError(err)
	fmt.Printf("Workspace created: %s", result.Name)

}

func Delete(client tfe.Client, name, org string) {
	ctx := context.Background()
	ws := GetByName(client, org, name)

	err := client.Workspaces.DeleteByID(ctx, ws.ID)
	helper.HandleError(err)
	fmt.Printf("Workspace deleted: %s", ws.Name)

}

func Purge(client tfe.Client, orgName string) {
	ctx := context.Background()
	workspaces := listWorkspaces(client, orgName)
	for _, v := range workspaces.Items {
		err := client.Workspaces.DeleteByID(ctx, v.ID)
		helper.HandleError(err)
	}
	fmt.Printf("All Workspace deleted")

}
