package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/costa86/tformer/helper"
	orgMod "github.com/costa86/tformer/organization"
	varMod "github.com/costa86/tformer/variable"

	configVerMod "github.com/costa86/tformer/config_version"
	runMod "github.com/costa86/tformer/run"
	userMod "github.com/costa86/tformer/user"
	wsMod "github.com/costa86/tformer/workspace"

	tfe "github.com/hashicorp/go-tfe"
)

var placeholder = "sample"
var token = flag.String("token", placeholder, "authentication API token")
var address = flag.String("address", "https://app.terraform.io", "terraform address")

func getClient() *tfe.Client {
	config := &tfe.Config{
		Token:             *token,
		RetryServerErrors: true,
		Address:           *address,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func wsList() {
	fmt.Println("List workspaces for an Organization")
	wsCmd := flag.NewFlagSet("ws_list", flag.ExitOnError)
	org := wsCmd.String("org", placeholder, "organization name")
	asJson := wsCmd.Bool("json", false, "output as json")

	wsCmd.Parse(flag.Args()[1:])
	wsMod.List(getClient(), *org, *asJson)
}
func wsCreate() {
	fmt.Println("Create Workspace")
	cmd := flag.NewFlagSet("ws_create", flag.ExitOnError)
	org := cmd.String("org", placeholder, "organization name")
	name := cmd.String("name", placeholder, "workspace name")

	cmd.Parse(flag.Args()[1:])
	wsMod.Create(*getClient(), *org, *name)
}
func wsLock() {
	fmt.Println("Lock Workspace")
	cmd := flag.NewFlagSet("ws_lock", flag.ExitOnError)
	id := cmd.String("id", placeholder, "Workspace id")

	cmd.Parse(flag.Args()[1:])
	wsMod.Lock(*getClient(), *id)
}
func wsDelete() {
	fmt.Println("Delete Workspace")
	cmd := flag.NewFlagSet("ws_delete", flag.ExitOnError)
	id := cmd.String("id", placeholder, "Workspace id")

	cmd.Parse(flag.Args()[1:])
	wsMod.Delete(*getClient(), *id)
}
func wsPurge() {
	fmt.Println("Delete all Workspaces")
	cmd := flag.NewFlagSet("ws_purge", flag.ExitOnError)
	orgName := cmd.String("org_name", placeholder, "Organization name")

	cmd.Parse(flag.Args()[1:])
	wsMod.Purge(*getClient(), *orgName)
}

func cvList() {
	fmt.Println("List Configuration Versions")
	cmd := flag.NewFlagSet("ws_lock", flag.ExitOnError)
	id := cmd.String("ws_id", placeholder, "Workspace id")

	cmd.Parse(flag.Args()[1:])
	configVerMod.List(getClient(), *id)
}

func cvDownload() {
	fmt.Println("Download Configuration Version")
	cmd := flag.NewFlagSet("ws_lock", flag.ExitOnError)
	id := cmd.String("id", placeholder, "Workspace id")

	cmd.Parse(flag.Args()[1:])
	configVerMod.Download(getClient(), *id)
}

func userGetMe() {
	cmd := flag.NewFlagSet("whoami", flag.ExitOnError)

	cmd.Parse(flag.Args()[1:])
	userMod.ReadCurrent(getClient())
}

func orgList() {
	fmt.Println("List Organizations")
	orgMod.List(getClient())
}

func orgDelete() {
	fmt.Println("Delete Organization")
	cmd := flag.NewFlagSet("org_delete", flag.ExitOnError)
	name := cmd.String("name", placeholder, "organization name")

	cmd.Parse(flag.Args()[1:])

	orgMod.Delete(*getClient(), *name)
}
func orgCreate() {
	fmt.Println("Create Organization")
	cmd := flag.NewFlagSet("org_create", flag.ExitOnError)
	name := cmd.String("name", placeholder, "organization name")
	email := cmd.String("email", placeholder, "organization email")

	cmd.Parse(flag.Args()[1:])

	orgMod.Create(*getClient(), *name, *email)
}

func runList() {
	fmt.Println("List Runs")
	cmd := flag.NewFlagSet("run_list", flag.ExitOnError)
	wsId := cmd.String("ws_id", placeholder, "workspace ID")
	cmd.Parse(flag.Args()[1:])

	runMod.List(getClient(), *wsId)
}

func varList() {
	fmt.Println("List Vars")
	cmd := flag.NewFlagSet("var_list", flag.ExitOnError)
	wsId := cmd.String("ws_id", placeholder, "workspace ID")
	cmd.Parse(flag.Args()[1:])

	varMod.List(getClient(), *wsId)
}
func varDelete() {
	fmt.Println("Delete Var")
	cmd := flag.NewFlagSet("var_delete", flag.ExitOnError)
	wsId := cmd.String("ws_id", placeholder, "workspace ID")
	varId := cmd.String("var_id", placeholder, "variable ID")

	cmd.Parse(flag.Args()[1:])

	varMod.Delete(*getClient(), *wsId, *varId)
}

func runCreate() {
	fmt.Println("Create Run")
	cmd := flag.NewFlagSet("run_create", flag.ExitOnError)
	wsId := cmd.String("ws_id", placeholder, "workspace ID")

	cmd.Parse(flag.Args()[1:])

	dir := cmd.Args()

	runMod.Create(*getClient(), *wsId, dir)
}

func varCreate() {
	fmt.Println("Create Variable")
	cmd := flag.NewFlagSet("var_create", flag.ExitOnError)
	wsId := cmd.String("ws_id", placeholder, "workspace id")
	key := cmd.String("key", placeholder, "key")
	value := cmd.String("value", placeholder, "value")
	description := cmd.String("description", placeholder, "description")
	category := cmd.String("category", "env", "category (env | terraform)")
	sensitive := cmd.Bool("sensitive", false, "sensitive")
	hcl := cmd.Bool("hcl", true, "hcl")

	cmd.Parse(flag.Args()[1:])

	validCategories := []string{"terraform", "env"}
	if !helper.Contains(validCategories, *category) {
		fmt.Println("[ERROR] Category must be in:", strings.Join(validCategories, ","))
		os.Exit(0)
	}

	variable := helper.Variable{
		Key:         *key,
		Value:       *value,
		Description: *description,
		Category:    *category,
		Sensitive:   *sensitive,
		HCL:         *hcl,
	}
	varMod.Create(*getClient(), *wsId, variable)
}

func main() {
	flag.Usage = func() {
		fmt.Println("* Usage: <common tags> <command> <command tags>\n* Common tags:")
		flag.PrintDefaults()
		fmt.Println("* Commands (use <command> -h for the tags):")
		//ws
		fmt.Println("ws_list")
		fmt.Println("ws_create")
		fmt.Println("ws_lock")
		fmt.Println("ws_delete")
		fmt.Println("ws_purge")
		//var
		fmt.Println("var_create")
		fmt.Println("var_list")
		fmt.Println("var_delete")
		//org
		fmt.Println("org_list")
		fmt.Println("org_create")
		fmt.Println("org_delete")

		//run
		fmt.Println("run_list")
		fmt.Println("run_create")
		//cv
		fmt.Println("cv_list")
		fmt.Println("cv_download")
		//user
		fmt.Println("whoami")

	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	switch flag.Arg(0) {

	case "ws_list":
		wsList()
	case "ws_create":
		wsCreate()
	case "ws_lock":
		wsLock()

	case "ws_delete":
		wsDelete()
	case "ws_purge":
		wsPurge()
	case "var_create":
		varCreate()
	case "var_list":
		varList()
	case "var_delete":
		varDelete()
	case "org_list":
		orgList()
	case "org_create":
		orgCreate()
	case "org_delete":
		orgDelete()

	case "cv_list":
		cvList()
	case "cv_download":
		cvDownload()
	case "run_list":
		runList()
	case "run_create":
		runCreate()
	case "whoami":
		userGetMe()

	default:
		flag.PrintDefaults()
		os.Exit(0)
	}
}
