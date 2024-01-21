package variable

import (
	"context"
	"fmt"
	"log"

	wsMod "github.com/costa86/tformer/workspace"

	"github.com/costa86/tformer/helper"
	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func update(client tfe.Client, workspaceName string, org string, id string, variable helper.Variable) {
	ctx := context.Background()
	ws := wsMod.GetByName(client, org, workspaceName)

	result, err := client.Variables.Update(ctx, ws.ID, id, tfe.VariableUpdateOptions{
		Type:        "vars",
		Key:         tfe.String(variable.Key),
		Value:       tfe.String(variable.Value),
		Description: tfe.String(variable.Description),
		Category:    tfe.Category(tfe.CategoryType(variable.Category)),
		HCL:         tfe.Bool(variable.HCL),
		Sensitive:   tfe.Bool(variable.Sensitive)})

	helper.HandleError(err)
	fmt.Printf("Variable updated: %s", result.Key)
}

func listVariables(client tfe.Client, workspaceName, org string) *tfe.VariableList {
	ctx := context.Background()
	ws := wsMod.GetByName(client, org, workspaceName)

	result, err := client.Variables.List(ctx, ws.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	return result
}

func Create(client tfe.Client, workspaceName string, org string, variable helper.Variable) {
	ctx := context.Background()
	ws := wsMod.GetByName(client, org, workspaceName)
	variables := listVariables(client, workspaceName, org)

	for _, v := range variables.Items {
		if v.Key == variable.Key {
			update(client, workspaceName, org, v.ID, variable)
			return
		}
	}

	result, err := client.Variables.Create(ctx, ws.ID, tfe.VariableCreateOptions{
		Type:        "vars",
		Key:         tfe.String(variable.Key),
		Value:       tfe.String(variable.Value),
		Description: tfe.String(variable.Description),
		Category:    tfe.Category(tfe.CategoryType(variable.Category)),
		HCL:         tfe.Bool(variable.HCL),
		Sensitive:   tfe.Bool(variable.Sensitive)})

	helper.HandleError(err)
	fmt.Printf("Variable created: %s", result.Key)
}

func Delete(client tfe.Client, workspaceName, orgNme, variableName string) {
	ctx := context.Background()
	ws := wsMod.GetByName(client, orgNme, workspaceName)
	vars := listVariables(client, workspaceName, orgNme)

	for _, v := range vars.Items {
		if v.Key == variableName {
			err := client.Variables.Delete(ctx, ws.ID, v.ID)
			helper.HandleError(err)
			fmt.Printf("Variable deleted: %s", variableName)
			return
		}
	}
}

func List(client *tfe.Client, name, org string) {
	list := listVariables(*client, name, org)
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "KEY", "VALUE", "CATEGORY", "DESCRIPTION", "HCL", "SENSITIVE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Key, v.Value, v.Category, v.Description, v.HCL, v.Sensitive)
	}

	tbl.Print()
}
