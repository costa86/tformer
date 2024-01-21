package organization

import (
	"context"
	"fmt"

	"github.com/costa86/tformer/helper"
	"github.com/fatih/color"
	"github.com/hashicorp/go-tfe"
	"github.com/rodaine/table"
)

func listOrganizations(client tfe.Client) *tfe.OrganizationList {
	ctx := context.Background()

	result, err := client.Organizations.List(ctx, nil)
	helper.HandleError(err)
	return result

}

func List(client *tfe.Client) {
	list := listOrganizations(*client)
	headerFmt := color.New(color.FgMagenta, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("NAME", "EMAIL", "EXTERNAL ID", "CREATED AT")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, v := range list.Items {
		tbl.AddRow(v.Name, v.Email, v.ExternalID, v.CreatedAt)
	}

	tbl.Print()
}

func Create(client tfe.Client, name, email string) {
	ctx := context.Background()

	result, err := client.Organizations.Create(ctx, tfe.OrganizationCreateOptions{Name: tfe.String(name), Email: tfe.String(email)})
	helper.HandleError(err)
	fmt.Printf("Organization created: %s", result.Name)

}

func Delete(client tfe.Client, name string) {
	ctx := context.Background()

	err := client.Organizations.Delete(ctx, name)
	helper.HandleError(err)
	fmt.Printf("Organization delete: %s", name)

}
