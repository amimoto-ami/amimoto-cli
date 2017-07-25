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
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
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
    optarg.Header("Sub command site")
    optarg.Add("a", "add", "Add site (string)", "")
    optarg.Add("e", "enable", "Enable site (string)", "")
    optarg.Add("d", "disable", "Disable site (string)", "")
    optarg.Add("r", "remove", "Remove site (string)", "")
    ch := optarg.Parse()
    <-ch

    if len(optarg.Remainder) >= 1 {
        cmd := optarg.Remainder[0]
        switch cmd {
        case "cache":
            cache()
        case "site":
            site()
        default:
            optarg.Usage()
        }
    } else {
        optarg.Usage()
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

func site() {
    var site string
    for opt := range optarg.Parse() {
        site = opt.String()
        if site != "" {
            switch opt.ShortName {
            case "a":
                fmt.Println("Site add: ", site)
                add(site)
            case "d":
                fmt.Println("Site disable: ", site)
                disable(site)
            case "e":
                fmt.Println("Site enableL: ", site)
                enable(site)
            case "r":
                fmt.Println("Site remove: ", site)
                remove(site)
            default:
                optarg.Usage()
            }
        } else {
        optarg.Usage()
        }
    }
}

func Exists(name string) bool {
    _, err := os.Stat(name)
    return !os.IsNotExist(err)
}

func add(site string) {
    setupCmd := "/usr/local/bin/wp-setup"

    // wp-setup
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

    if Exists("/etc/nginx/conf.d/"+site+".conf.disable") {
        enable(site)
    }

    // read DB Configuration and show
    dbConf, readErr := readAmimotoDBConf(site)
    if readErr != nil {
        panic(readErr)
    }
    fmt.Println()
    fmt.Println("DB Host: ", dbConf.host)
    fmt.Println("DB Name: ", dbConf.name)
    fmt.Println("DB User: ", dbConf.user)
    fmt.Println("DB Pass: ", dbConf.pass)
}

func nginxConfRename(conf string, enable bool) error {
    var err error
    confDir := "/etc/nginx/conf.d/"
    if enable {
        if Exists(filepath.Join(confDir,conf+".disable")) {
            err = os.Rename(filepath.Join(confDir,conf+".disable"), filepath.Join(confDir,conf))
            return err
        }
    } else {
        if Exists(filepath.Join(confDir,conf)) {
            err = os.Rename(filepath.Join(confDir,conf), filepath.Join(confDir,conf+".disable"))
            return err
        }
    }
    return nil
}

func nginxReload() error {
    cmd := exec.Command("/sbin/service", "nginx", "reload")
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        return err
    }

    cmd.Start()
    scanner := bufio.NewScanner(stdout)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }
    cmd.Wait()
    return nil
}

func disable(site string) {
    var err error

    if Exists("/etc/nginx/conf.d/"+site+".conf") {
        // rename nginx conf files
        confs :=  map[string]string{
            "front": site+".conf",
            "ssl": site+"-ssl.conf",
            "backend": site+".backend.conf",
        }
        for _, filename := range confs {
            err = nginxConfRename(filename, false)
            if err != nil {
                panic(err)
            }
        }

        // nginx conf reload
        err = nginxReload()
        if err != nil {
            panic(err)
        }

        fmt.Println("Site (" + site + ") disabled!")

    } else {
        fmt.Println("Site (" + site + ") config not found.")
        os.Exit(1)
    }
}

func enable(site string) {
    var err error

    if Exists("/etc/nginx/conf.d/"+site+".conf.disable") {
        // rename nginx conf files
        confs :=  map[string]string{
            "front": site+".conf",
            "ssl": site+"-ssl.conf",
            "backend": site+".backend.conf",
        }
        for _, filename := range confs {
            err = nginxConfRename(filename, true)
            if err != nil {
                panic(err)
            }
        }

        // nginx conf reload
        err = nginxReload()
        if err != nil {
            panic(err)
        }

        fmt.Println("Site " + site + " enabled!")

    } else {
        fmt.Println("Site (" + site + ") config not found.")
        os.Exit(1)
    }
}

func remove(site string) {
    if Exists(filepath.Join("/opt/local/", site+".json")) {
        // get DB Connect info
        dbConf, readErr := readAmimotoDBConf(site)
        if readErr != nil {
            panic(readErr)
        }

        //drop database
        db, dbErr := sql.Open("mysql", dbConf.user+":"+dbConf.pass+"@/"+dbConf.name)
        if dbErr != nil {
            panic(dbErr)
        }
        defer db.Close()
        _, queryErr := db.Query("DROP DATABASE "+dbConf.name)
        if queryErr != nil {
            panic(queryErr.Error())
        }

        // remove site.conf files
        var err error
        dir := "/etc/nginx/conf.d/"
        confs :=  map[string]string{
            "front": filepath.Join(dir, site+".conf"),
            "ssl": filepath.Join(dir, site+"-ssl.conf"),
            "backend": filepath.Join(dir, site+".backend.conf"),
            "db_json": filepath.Join("/opt/local/", site+".json"),
            "db_sql": filepath.Join("/opt/local/", "createdb-"+site+".sql"),
        }
        for _, filename := range confs {
            if Exists(filename) {
                err = os.Remove(filename)
                if err != nil {
                    panic(err)
                }
            }
            if Exists(filename+".disable") {
                err = os.Remove(filename+".disable")
                if err != nil {
                    panic(err)
                }
            }
        }
        err = os.RemoveAll(filepath.Join("/var/www/vhosts/", site))

        // nginx conf reload
        err = nginxReload()
        if err != nil {
            panic(err)
        }

        fmt.Println("Site " + site + " removed!")
    }
}

func readAmimotoDBConf(site string) (dbConf, error) {
    filename := "/opt/local/" + site + ".json"
    if Exists(filename) {
        var c interface{}
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

    } else {
        return dbConf{}, nil
    }
}
