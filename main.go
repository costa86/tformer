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
	"gopkg.in/yaml.v2"
)

var placeholder = "sample"
var configFile = "config.yaml"

func getTfConfig() helper.TfConfig {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	var tfConfig helper.TfConfig
	err = yaml.Unmarshal(yamlFile, &tfConfig)
	if err != nil {
		panic(err)
	}
	return tfConfig
}

func getClient() *tfe.Client {
	tfConfig := getTfConfig()

	config := &tfe.Config{
		Token:             tfConfig.Token,
		RetryServerErrors: true,
		Address:           tfConfig.Address,
	}

	client, err := tfe.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func wsList() {
	org := getTfConfig().Organization
	fmt.Println("List workspaces for Organization")
	wsCmd := flag.NewFlagSet("ws_list", flag.ExitOnError)
	asJson := wsCmd.Bool("json", false, "output as json")

	wsCmd.Parse(flag.Args()[1:])
	wsMod.List(getClient(), org, *asJson)
}

func wsCreate() {
	org := getTfConfig().Organization
	fmt.Println("Create Workspace")
	cmd := flag.NewFlagSet("ws_create", flag.ExitOnError)
	name := cmd.String("name", placeholder, "workspace name")

	cmd.Parse(flag.Args()[1:])
	wsMod.Create(*getClient(), org, *name)
}

func wsLock() {
	fmt.Println("Lock Workspace")
	cmd := flag.NewFlagSet("ws_lock", flag.ExitOnError)
	name := cmd.String("name", placeholder, "workspace name")
	org := getTfConfig().Organization

	cmd.Parse(flag.Args()[1:])
	wsMod.LockOrUnlock(*getClient(), *name, org, true)
}
func wsLockAll() {
	fmt.Println("Lock all Workspaces")
	cmd := flag.NewFlagSet("ws_lock_all", flag.ExitOnError)
	org := getTfConfig().Organization

	cmd.Parse(flag.Args()[1:])
	wsMod.LockOrUnlockAll(*getClient(), org, true)
}

func wsUnlockAll() {
	fmt.Println("Unlock all Workspaces")
	cmd := flag.NewFlagSet("ws_unlock_all", flag.ExitOnError)
	org := getTfConfig().Organization

	cmd.Parse(flag.Args()[1:])
	wsMod.LockOrUnlockAll(*getClient(), org, false)
}

func wsUnlock() {
	fmt.Println("Unlock Workspace")
	cmd := flag.NewFlagSet("ws_lock", flag.ExitOnError)
	name := cmd.String("name", placeholder, "workspace name")
	org := getTfConfig().Organization

	cmd.Parse(flag.Args()[1:])
	wsMod.LockOrUnlock(*getClient(), *name, org, false)
}

func wsDelete() {
	fmt.Println("Delete Workspace")
	cmd := flag.NewFlagSet("ws_delete", flag.ExitOnError)
	name := cmd.String("name", placeholder, "workspace name")

	cmd.Parse(flag.Args()[1:])
	wsMod.Delete(*getClient(), *name, getTfConfig().Organization)
}

func wsPurge() {
	fmt.Println("Delete all Workspaces")
	cmd := flag.NewFlagSet("ws_purge", flag.ExitOnError)

	cmd.Parse(flag.Args()[1:])
	wsMod.Purge(*getClient(), getTfConfig().Organization)
}

func cvList() {
	fmt.Println("List Configuration Versions")
	cmd := flag.NewFlagSet("cv_list", flag.ExitOnError)
	ws := cmd.String("ws", placeholder, "workspace name")

	cmd.Parse(flag.Args()[1:])
	configVerMod.List(getClient(), *ws, getTfConfig().Organization)
}

func cvDownload() {
	fmt.Println("Download Configuration Version")
	cmd := flag.NewFlagSet("cv_download", flag.ExitOnError)
	id := cmd.String("id", placeholder, "configuration version id")

	cmd.Parse(flag.Args()[1:])
	configVerMod.Download(getClient(), *id)
}

func cvGet() {
	fmt.Println("Get Configuration Version")
	cmd := flag.NewFlagSet("cv_get", flag.ExitOnError)
	id := cmd.String("id", placeholder, "configuration version id")

	cmd.Parse(flag.Args()[1:])
	configVerMod.Get(getClient(), *id)
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
	ws := cmd.String("ws", placeholder, "workspace name")
	cmd.Parse(flag.Args()[1:])

	runMod.List(getClient(), *ws, getTfConfig().Organization)
}

func varList() {
	fmt.Println("List Variables")
	cmd := flag.NewFlagSet("var_list", flag.ExitOnError)
	ws := cmd.String("ws", placeholder, "workspace name")
	cmd.Parse(flag.Args()[1:])

	varMod.List(getClient(), *ws, getTfConfig().Organization)
}
func varDelete() {
	fmt.Println("Delete Variable")
	cmd := flag.NewFlagSet("var_delete", flag.ExitOnError)
	ws := cmd.String("ws", placeholder, "workspace name")
	key := cmd.String("key", placeholder, "variable key")

	cmd.Parse(flag.Args()[1:])

	varMod.Delete(*getClient(), *ws, getTfConfig().Organization, *key)
}

func getRunFlags(title string) (string, string, string, bool, string) {
	fmt.Println(title)
	cmd := flag.NewFlagSet("run_create", flag.ExitOnError)
	wsName := cmd.String("ws", placeholder, "workspace name")
	message := cmd.String("msg", placeholder, "message")
	dir := cmd.String("dir", placeholder, "dir")
	autoApply := cmd.Bool("aa", false, "auto-apply")
	cvId := cmd.String("cv_id", "latest", "configuration version id")

	cmd.Parse(flag.Args()[1:])

	return *wsName, *message, *dir, *autoApply, *cvId

}

func runCreate() {
	wsName, message, dir, autoApply, cvId := getRunFlags("Create Run")
	ws := wsMod.GetByName(*getClient(), getTfConfig().Organization, wsName)
	runMod.CreateOrDestroy(*getClient(), ws.ID, message, dir, autoApply, false, cvId)
}
func runDestoy() {
	wsName, message, dir, autoApply, cvId := getRunFlags("Destroy Run")
	ws := wsMod.GetByName(*getClient(), getTfConfig().Organization, wsName)
	runMod.CreateOrDestroy(*getClient(), ws.ID, message, dir, autoApply, true, cvId)

}
func runGet() {
	fmt.Println("Get Run")
	cmd := flag.NewFlagSet("run_get", flag.ExitOnError)
	id := cmd.String("id", placeholder, "run id")
	cmd.Parse(flag.Args()[1:])
	runMod.Get(*getClient(), *id)
}

func sample() {
	fmt.Println("sample")
	cmd := flag.NewFlagSet("sample", flag.ExitOnError)
	autoApply := cmd.Bool("apply", true, "auto-apply")
	other := cmd.Int("other", 1, "other")

	cmd.Parse(flag.Args()[1:])
	println(*autoApply)
	println(*other)

}

func varCreate() {
	fmt.Println("Create Variable")
	cmd := flag.NewFlagSet("var_create", flag.ExitOnError)
	ws := cmd.String("ws", placeholder, "workspace name")
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
	varMod.Create(*getClient(), *ws, getTfConfig().Organization, variable)
}

func main() {
	flag.Usage = func() {
		fmt.Printf("* Usage: <command> <command flags>\n* Make sure '%s' exists in the current dir\n", configFile)
		flag.PrintDefaults()
		fmt.Println("* Commands (use <command> -h for the flags):")
		//ws
		fmt.Println("ws_list")
		fmt.Println("ws_create")
		fmt.Println("ws_lock")
		fmt.Println("ws_lock_all")
		fmt.Println("ws_unlock")
		fmt.Println("ws_unlock_all")
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
		fmt.Println("run_destroy")
		fmt.Println("run_get")
		//cv
		fmt.Println("cv_list")
		fmt.Println("cv_download")
		fmt.Println("cv_get")
		//user
		fmt.Println("whoami")

	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	flag.Parse()

	switch flag.Arg(0) {
	// ws
	case "ws_list":
		wsList()
	case "ws_create":
		wsCreate()
	case "ws_lock":
		wsLock()
	case "ws_unlock":
		wsUnlock()
	case "ws_lock_all":
		wsLockAll()
	case "ws_unlock_all":
		wsUnlockAll()
	case "ws_delete":
		wsDelete()
	case "ws_purge":
		wsPurge()
	// var
	case "var_create":
		varCreate()
	case "var_list":
		varList()
	case "var_delete":
		varDelete()
	// org
	case "org_list":
		orgList()
	case "org_create":
		orgCreate()
	case "org_delete":
		orgDelete()
	// cv
	case "cv_list":
		cvList()
	case "cv_download":
		cvDownload()
	case "cv_get":
		cvGet()
	// run
	case "run_list":
		runList()
	case "run_create":
		runCreate()
	case "run_destroy":
		runDestoy()
	case "run_get":
		runGet()
	// user
	case "whoami":
		userGetMe()
	case "sample":
		sample()
	default:
		flag.PrintDefaults()
		os.Exit(0)
	}
}
