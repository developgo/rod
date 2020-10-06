package main

import (
	"net/http"
	"strings"

	"github.com/go-rod/rod/lib/utils"
	"github.com/ysmood/gson"
)

func main() {
	list := getDeviceList()

	code := ``
	for _, d := range list.Arr() {
		d = d.Get("device")
		name := d.Get("title").String()
		code += utils.S(`

			// {{.name}} device
			{{.name}} = Device{gson.NewFrom({{.desc}})}`,
			"name", normalizeName(name),
			"title", name,
			"desc", utils.EscapeGoString(strings.TrimSpace(d.JSON("\t", "  "))),
		)
	}

	code = utils.S(`// generated by running "go generate" on project root

		package devices

		var (
			{{.code}}
		)
	`, "code", code)

	path := "./lib/devices/list.go"
	utils.E(utils.OutputFile(path, code))

	utils.Exec("gofmt", "-s", "-w", path)
	utils.Exec("goimports", "-w", path)
	utils.Exec("misspell", "-w", "-q", path)
}

func getDeviceList() gson.JSON {
	// we use the list from the web UI of devtools
	res, err := http.Get(
		"https://raw.githubusercontent.com/ChromeDevTools/devtools-frontend/master/front_end/emulated_devices/module.json",
	)
	utils.E(err)

	return gson.New(res.Body).Get("extensions")
}

func normalizeName(name string) string {
	name = strings.ReplaceAll(name, "/", "or")

	list := []string{}
	for _, s := range strings.Split(name, " ") {
		if len(s) > 1 {
			list = append(list, strings.ToUpper(s[0:1])+s[1:])
		} else {
			list = append(list, strings.ToUpper(s))
		}
	}

	return strings.Join(list, "")
}
