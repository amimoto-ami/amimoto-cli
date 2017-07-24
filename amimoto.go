package main

import (
	"fmt"
	"syscall"
	//"flag"
	"github.com/cloudbuy/go-pkg-optarg"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	optarg.Header("Sub command cache")
	optarg.Add("p", "purge", "Clear NGINX proxy cache", false)
	optarg.Header("Sub command add {site name}")
	//optarg.Add("s", "site", "Adds site with WordPress files and DB creation", "")

	ch := optarg.Parse()
	<-ch

	if len(optarg.Remainder) >= 1 {
		cmd := optarg.Remainder[0]
		switch cmd {
		case "cache":
			cache()
		case "add":
			if len(optarg.Remainder) > 1 {
				add(optarg.Remainder[1])
			}
		}
	} else {
		optarg.Usage()
	}
}

func add(site string) {
	fmt.Println("Site URL: ", site)

	binary, lookErr := exec.LookPath("/usr/local/bin/wp-setup")
	if lookErr != nil {
		panic(lookErr)
	}
	args := []string{"/usr/local/bin/wp-setup", site}
	env := os.Environ()

	execErr := syscall.Exec(binary, args, env)
	if execErr != nil {
		panic(execErr)
	}
}

func cache() {
	var purge bool
	var dir string
	dir = "/var/cache/nginx/proxy_cache/"
	for opt := range optarg.Parse() {
		switch opt.ShortName {
		case "p":
			purge = opt.Bool()
		}
	}
	if purge {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			err := os.RemoveAll(filepath.Join(dir, file.Name()))
			if err != nil {
				panic(err)
			}
		}
	} else {
		paths := dirwalk(dir)
		fmt.Println("Cached URLs: ", len(paths))
	}
}

func dirwalk(dir string) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if file.IsDir() {
			paths = append(paths, dirwalk(filepath.Join(dir, file.Name()))...)
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
