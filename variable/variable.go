package variable

import (
	"context"
	"fmt"
	"log"

	"github.com/costa86/tformer/helper"
	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func Create(client tfe.Client, workspaceId string, variable helper.Variable) {
	ctx := context.Background()

	result, err := client.Variables.Create(ctx, workspaceId, tfe.VariableCreateOptions{
		Type:        "vars",
		Key:         tfe.String(variable.Key),
		Value:       tfe.String(variable.Value),
		Description: tfe.String(variable.Description),
		Category:    tfe.Category(tfe.CategoryType(variable.Category)),
		HCL:         tfe.Bool(variable.HCL),
		Sensitive:   tfe.Bool(variable.Sensitive)})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Variable created: %s", result.Key)
}

func listVariables(client tfe.Client, workspaceId string) *tfe.VariableList {
	ctx := context.Background()

	result, err := client.Variables.List(ctx, workspaceId, nil)
	if err != nil {
		log.Fatal(err)
	}
	return result
}
func Delete(client tfe.Client, workspaceId, variableId string) {
	ctx := context.Background()

	err := client.Variables.Delete(ctx, workspaceId, variableId)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Variable deleted: %s", variableId)
}

func List(client *tfe.Client, name string) {
	list := listVariables(*client, name)
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "KEY", "VALUE", "CATEGORY", "DESCRIPTION", "HCL", "SENSITIVE")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.ID, v.Key, v.Value, v.Category, v.Description, v.HCL, v.Sensitive)
	}

	tbl.Print()
}
