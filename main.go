package main

import (
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-version"
)

var DEBUG bool
var COMMIT_LOG bool
var NEXTCLOUD_VERSIONS *[]version.Version

const API_BASE = "https://apps.nextcloud.com/"

type ApiRelease struct {
	Version             string        `json:"version"`
	Download            string        `json:"download"`
	PhpVersionSpec      string        `json:"phpVersionSpec"`
	PlatformVersionSpec string        `json:"platformVersionSpec"`
	Licenses            []string      `json:"licenses"`
	Databases           []interface{} `json:"databases"`
	PhpExtension        []interface{} `json:"phpExtensions"`
}

type ApiAppTranslation struct {
	Name        string `json:"name"`
	Summary     string `json:"summary"`
	Description string `json:"description"`
}

type ApiApp struct {
	Id           string                       `json:"id"`
	Website      string                       `json:"website"`
	Releases     []ApiRelease                 `json:"releases"`
	Translations map[string]ApiAppTranslation `json:"translations"`
}
type ApiJson []ApiApp

type App struct {
	Sha256      string   `json:"sha256"`
	Url         string   `json:"url"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Homepage    string   `json:"homepage"`
	Licenses    []string `json:"licenses"`
}
type AppJson map[string]App

func major(t string) string {
	return strings.Split(t, ".")[0]
}
func loadFile(t string) (AppJson, error) {
	fname := major(t) + ".json"
	log.Print("Loading " + fname)
	m := make(AppJson)
	var err error
	if _, err := os.Stat(fname); !os.IsNotExist(err) {
		file, _ := os.OpenFile(fname, os.O_CREATE|os.O_RDONLY, 0644)
		defer file.Close()
		dat, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Failed to read file %s: %e", fname, err)
		}
		err = json.Unmarshal(dat, &m)
		if err != nil {
			log.Printf("Failed to parse %s: %e", fname, err)
		}
		log.Printf("Loaded %s", fname)
	}
	return m, err
}

func writeLog(t string, ao, an AppJson) {
	m := major(t)
	file, _ := os.OpenFile(m+"-new.log", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	log.Printf("Writing %s-new.log", m)
	for k, na := range an {
		oa, isOld := ao[k]
		if !isOld {
			file.WriteString(fmt.Sprintf("ADD %s %s\n", k, na.Version))
		} else if isOld && na.Version != oa.Version {
			file.WriteString(fmt.Sprintf("UPD %s %s -> %s\n", k, oa.Version, na.Version))
		}
	}
	log.Printf("Replacing %s.log with %s-new.log", m, m)
	os.Rename(m+"-new.log", m+".log")
}

func writeFile(t string, c AppJson) {
	m := major(t)
	file, _ := os.OpenFile(m+"-new.json", os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	log.Printf("Writing %s-new.json", m)
	enc.Encode(c)
	log.Printf("Replacing %s.json with %s-new.json", m, m)
	os.Rename(m+"-new.json", m+".json")
}

func prefetch(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print("Prefetch failed for: ", url, err)
		return "", err
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("Prefetch failed reading body for: ", url, err)
		return "", err
	}
	sha256 := fmt.Sprintf("%x", sha256.Sum256(contents))

	return sha256, err
}

// copy every element from every map into the resulting map
// meaning, merge all maps with the later maps having precedence over the previous one(s)
func mergeAs(as ...AppJson) AppJson {
	res := make(AppJson)
	for _, m := range as {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}

func queryApi(v string) (ApiJson, error) {
	url := API_BASE + "/api/v1/platform/" + v + "/apps.json"
	log.Printf("Querying API at %s", url)
	resp, err := http.Get(url)
	var apiResponse ApiJson
	if err != nil || resp.StatusCode != http.StatusOK {
		log.Print("API query failed (", resp.Status, ") for: ", url, err)
		return apiResponse, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("API query failed to read body: ", url, err)
		return apiResponse, err
	}
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		log.Print("API query failed to parse JSON: ", string(body), err)
		return apiResponse, err
	}
	return apiResponse, nil
}

func update(v string, apps []string) {
	log.Printf("Starting to process apps for version %s", v)
	j, err := queryApi(v)
	if err != nil {
		panic(err)
	}
	ao, err := loadFile(v)
	if err != nil {
		panic(err)
	}
	an := make(AppJson)
	for _, a := range j {
		var na App
		na.Description = a.Translations["en"].Description
		na.Homepage = a.Website
		latestVer, _ := version.NewVersion("0")
		for _, rel := range a.Releases {
			ver, _ := version.NewVersion(rel.Version)
			if ver.GreaterThan(latestVer) && len(ver.Prerelease()) == 0 {
				latestVer = ver
				na.Version = ver.String()
				na.Url = rel.Download
				na.Licenses = rel.Licenses
			}
		}
		if len(apps) > 0 {
			for _, app := range apps {
				if a.Id == app {
					log.Printf("Found app %s (%s) at %s", a.Id, na.Version, na.Url)
					an[a.Id] = na
					break
				}
			}
		} else {
			log.Printf("Found app %s (%s) at %s", a.Id, na.Version, na.Url)
			an[a.Id] = na
		}
	}

	for k, _ := range an {
		needsPrefetch := false
		oa, isOld := ao[k]
		na, _ := an[k]

		if !isOld {
			log.Printf("New app found -> prefetching %s (%s) from %s", k, na.Version, na.Url)
			needsPrefetch = true
		} else {
			if na.Version != oa.Version || na.Url != oa.Url {
				needsPrefetch = true
				log.Printf("App was updated -> prefetching %s (%s) from %s", k, na.Version, na.Url)
			}
		}

		if needsPrefetch {
			sha256, err := prefetch(na.Url)
			if err != nil {
				continue
			}
			na.Sha256 = sha256
		} else {
			na.Sha256 = oa.Sha256
		}
		an[k] = na
	}

	writeFile(v, mergeAs(ao, an))
	writeLog(v, ao, an)
	log.Printf("Finished processing version %s", v)
}

func main() {
	_, DEBUG = os.LookupEnv("DEBUG")
	_, COMMIT_LOG = os.LookupEnv("COMMIT_LOG")
	var isSet bool
	NEXTCLOUD_VERSIONS_ENV, isSet := os.LookupEnv("NEXTCLOUD_VERSIONS")
	if !isSet {
		log.Fatal("NEXTCLOUD_VERSIONS needs to be set to nextcloud release(s)")
		os.Exit(1)
	}
	var NEXTCLOUD_VERSIONS []*version.Version
	for _, s := range strings.Split(strings.Trim(NEXTCLOUD_VERSIONS_ENV, "\""), ",") {
		v, err := version.NewVersion(s)
		if err != nil {
			log.Fatal("Error parsing NEXTCLOUD_VERSIONS")
			panic(err)
		}
		NEXTCLOUD_VERSIONS = append(NEXTCLOUD_VERSIONS, v)
	}
	apps := flag.String("apps", "", "Apps to fetch. Defaults to all")
	flag.Parse()

	for _, v := range NEXTCLOUD_VERSIONS {
		if *apps == "" { // https://github.com/golang/go/issues/35130
			update(v.String(), nil)
		} else {
			update(v.String(), strings.Split(*apps, ","))
		}
	}
}
