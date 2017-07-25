package main

import (
    "fmt"
    "bufio"
    "github.com/cloudbuy/go-pkg-optarg"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "encoding/json"
    "github.com/koron/go-dproxy"
)

type dbConf struct {
    name string
    user string
    pass string
    host string
}

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
    setupCmd := "/usr/local/bin/wp-setup"

    _, lookErr := exec.LookPath(setupCmd)
    if lookErr != nil {
        panic(lookErr)
    }

    cmd := exec.Command(setupCmd, site)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        panic(err)
    }

    cmd.Start()
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
        //fmt.Println()
    }
    cmd.Wait()

    db, err := readAmimotoDBConf(site)
    if err != nil {
        panic(err)
    }
    fmt.Println()
    fmt.Println("DB Host: ", db.host)
    fmt.Println("DB Name: ", db.name)
    fmt.Println("DB User: ", db.user)
    fmt.Println("DB Pass: ", db.pass)
}

func readAmimotoDBConf(site string) (dbConf, error) {
    var c interface{}

    filename := "/opt/local/" + site + ".json"
    jsonString, err := ioutil.ReadFile(filename)
    if err != nil {
        return dbConf{}, err
    }
    err = json.Unmarshal(jsonString, &c)
    if err != nil {
        return dbConf{}, err
    }
    db_name, _ := dproxy.New(c).M("wordpress").M("db").M("db_name").String()
    db_user, _ := dproxy.New(c).M("wordpress").M("db").M("user_name").String()
    db_pass, _ := dproxy.New(c).M("wordpress").M("db").M("password").String()
    db_host, _ := dproxy.New(c).M("wordpress").M("db").M("host").String()
    conf := dbConf{db_name, db_user, db_pass, db_host}

    return conf, nil
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
    paths := dirwalk(dir)
    fmt.Println("Cached URLs: ", len(paths))
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
        fmt.Println("Cache purge successfull!")
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
