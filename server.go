// © 2014 the Collectinator Authors under the WTFPL license. See AUTHORS for the list of authors.

package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/auth"
	"github.com/martini-contrib/render"
)

var fileDir string

type user struct {
	Name string
	Pass string
}

var users []user

func main() {
	flag.StringVar(&fileDir, "dir", "markdown", "directory to serve")
	userFile := flag.String("users", "users.json", "user definition file (JSON)")
	flag.Parse()
	m := martini.Classic()

	readUsers(*userFile)
	m.Use(auth.BasicFunc(authUser))
	m.Use(render.Renderer())

	m.Use(martini.Static("app", martini.StaticOptions{Exclude: "/files"}))
	m.Use(martini.Static(fileDir, martini.StaticOptions{Prefix: "files"}))
	m.Put("/files/**", writeFile)
	m.Post("/files/**", listFiles)
	m.Delete("/files/**", deleteFiles)

	m.Run()
}

func readUsers(file string) {
	if data, err := ioutil.ReadFile(file); err != nil {
		log.Fatal(err)
	} else {
		if err := json.Unmarshal(data, &users); err != nil {
			log.Fatal(err)
		}
	}
}

func authUser(user, pass string) bool {
	for _, u := range users {
		if auth.SecureCompare(user, u.Name) && auth.SecureCompare(pass, u.Pass) {
			return true
		}
	}
	return false
}

func deleteFiles(params martini.Params, rr render.Render) {
	err := os.RemoveAll(path.Join(fileDir, params["_1"]))
	if err != nil {
		rr.JSON(500, err)
		return
	}
	rr.JSON(200, params["_1"])
}

type info struct {
	Name  string
	IsDir bool
}

func listFiles(params martini.Params, rr render.Render) {
	files := getFiles(fileDir+"/", path.Join(fileDir, params["_1"]))

	rr.JSON(200, files)
}

func getFiles(base, pth string) (names []info) {
	files, _ := ioutil.ReadDir(pth)

	for _, file := range files {
		filePath := strings.Replace(path.Join(pth, file.Name()), base, "", -1)
		names = append(names, info{
			Name:  filePath,
			IsDir: file.IsDir(),
		})
		if file.IsDir() {
			names = append(names, getFiles(base, filePath)...)
		}
	}
	return
}

func writeFile(req *http.Request, params martini.Params, rr render.Render) {
	fileName := params["_1"]
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		rr.JSON(500, err)
		return
	}

	file := path.Join(fileDir, fileName)
	dir := path.Dir(file)

	err = os.MkdirAll(dir, 0755)
	if err != nil {
		rr.JSON(500, err)
		return
	}

	err = ioutil.WriteFile(
		file,
		data,
		0666,
	)
	if err != nil {
		rr.JSON(500, err)
		return
	}
	rr.Data(200, data)
}
