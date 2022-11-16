package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"runtime"
)

type SysInfo struct {
	OS   string
	Arch string
}

func main() {
	var execCmds = []*exec.Cmd{} // List to store every command to kill later
	sysInfo := SysInfo{runtime.GOOS, runtime.GOARCH}

	templates := make(map[string]*template.Template)
	templates["index"] = template.Must(template.ParseFiles("web/templates/index.html", "web/templates/_auth.html"))

	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Kill all the commands still runing
		killCmds(execCmds)

		if cmd := r.FormValue("cmd"); cmd != "" {
			log.Printf("Command input: %s", cmd)

			var app string
			var args []string
			var execCmd *exec.Cmd
			switch sysInfo.OS {
			case "windows":
				app = "cmd"
				args = append(args, "/C")
			case "linux":
				app = "/bin/sh"
				args = append(args, "-c")
			default:
				log.Println("Using default app and args")
			}

			if app != "" {
				args = append(args, cmd)
				execCmd = exec.Command(app, args...)
			} else {
				execCmd = exec.Command(cmd)
			}
			log.Printf("Command to run: %s", execCmd.Args)

			var stdout, stderr bytes.Buffer
			execCmd.Stdout = &stdout
			execCmd.Stderr = &stderr

			if err := execCmd.Run(); err != nil {
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				w.Write(stderr.Bytes())
				return
			} else {
				// Store the command in a list to kill
				//	if the page is refreshed
				// 	or if another command is sent
				execCmds = append(execCmds, execCmd)

				out := stdout.String()
				log.Println(out)
				w.Write(stdout.Bytes())
				return
			}
		}

		if err := templates["index"].ExecuteTemplate(w, "index.html", map[string]interface{}{"SysInfo": sysInfo}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Println("Application started")
	if err := http.ListenAndServe(":8189", nil); err != nil {
		log.Fatal(err)
	}
}

func killCmds(cmds []*exec.Cmd) {
	for _, v := range cmds {
		v.Process.Kill()
	}
}
