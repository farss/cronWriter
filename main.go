package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var config map[string]interface{}

func myHandel(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		buf, _ := ioutil.ReadAll(r.Body)
		jsonData := make(map[string][][]string)
		err := json.Unmarshal(buf, &jsonData)
		if err != nil {
			fmt.Println(err)
			//log.Fatal(err)
		}
		fmt.Println(jsonData)
		fmt.Println(len(jsonData))

		//加在配置
		f, err := os.OpenFile("config.json", os.O_RDONLY, 0777)
		if err != nil {
			fmt.Println(err)
		}
		defer f.Close()
		buf, _ = ioutil.ReadAll(f)
		fmt.Println(string(buf))
		//config := make(map[string][]string)
		err = json.Unmarshal(buf, &config)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(config["list"])
		if _, ok := config["php"]; !ok {
			w.WriteHeader(500)
			w.Write([]byte("error:config:php not found"))
			return
		}
		php := config["php"].(string)
		if _, ok := config["list"]; !ok {
			w.WriteHeader(500)
			w.Write([]byte("error:config:list not found"))
			return
		}
		list := (config["list"]).(map[string]interface{})

		if len(jsonData) != 1 {
			w.WriteHeader(500)
			w.Write([]byte("error:jsut one id can be handel"))
			return
		}
		for file, content := range jsonData {
			fmt.Println(content)
			_, ok := list[file]
			if !ok {
				w.WriteHeader(500)
				w.Write([]byte("error:config:" + file + " does not exist"))
				return
			}
			c := (list[file]).(string)
			f, err := os.OpenFile(c, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				w.Write([]byte(err.Error()))
				return
			}
			defer f.Close()

			fileContent := ""

			if _, ok := config["pre"]; ok {
				pre := (config["pre"]).([]interface{})
				for _, v := range pre {
					fileContent += v.(string) + "\n"
				}
			}

			for _, v := range content {
				fileContent += v[1] + " " + php + " " + v[0] + "\n"
			}
			fileContent = strings.TrimSuffix(fileContent, "\n")
			fmt.Println(c)
			fmt.Println(fileContent)
			len, err := f.WriteString(fileContent)
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(500)
				w.Write([]byte("error:fail to save config"))
				return
			}
			fmt.Println(len)
		}
		w.WriteHeader(200)
	} else {
		w.WriteHeader(500)
	}
}

func main() {
	p := flag.String("p", "5202", "port")
	flag.Parse()

	http.HandleFunc("/", myHandel)

	err := http.ListenAndServe(":"+*p, nil)
	if err != nil {
		log.Fatal(err)
	}
}
