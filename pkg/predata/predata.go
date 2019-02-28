package predata

import (
	"fmt"
	"log"
	"net/url"
	"path"
	"path/filepath"

	ut "minibox.ai/pkg/utils"
)

var (
	hubScheme = "hub"
	hosthub   = "minibox.ai"
)

type ErrInvalidDataSource struct {
	Src string
}

func (err *ErrInvalidDataSource) Error() string {
	return fmt.Sprintf("Invalid data soruce %s", err.Src)
}

//go:generate stringer  -type=DataSource
type DataSource int

const (
	HubData DataSource = iota
	UserData
	PrivateData
	WideData
	LocalData
	InvalidData
)

func ParseURL(uri string) (DataSource, *url.URL, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return InvalidData, nil, &ErrInvalidDataSource{uri}
	}

	source, _, resu := parsePath(*u)
	if source == InvalidData {
		return source, nil, &ErrInvalidDataSource{uri}
	}
	uu, _ := url.Parse(resu)
	return source, uu, nil
}

func ImportFromHub(u url.URL, target string) error {
	return Download(target, u.String())
}

func LocalPath(pat string) bool {
	if len(pat) == 0 {
		return false
	}

	if []rune(pat)[0] == '.' {
		return true
	} else if len(filepath.VolumeName(pat)) > 0 {
		return true
	}

	return false
}

func SetHubSite(site string) {
	hosthub = site
}

func parsePath(u url.URL) (source DataSource, pat, resURI string) {
	if u.Scheme == hubScheme {
		u.Path = path.Join(u.Host, u.Path)
		u.Host = hosthub
		u.Scheme = "https"
	}

	dir, file := filepath.Split(u.Path)
	log.Printf("scheme: %s, dir: %s, file %s", u.Scheme, dir, file)

	if u.Host == hosthub && ut.Empty(dir) && len(file) > 0 {
		source = HubData
		pat = file
		resURI = fmt.Sprintf("https://%s/datasets/%s", hosthub, file)
	} else if u.Host == hosthub && len(dir) > 0 && len(file) > 0 {
		source = UserData
		pat = u.Path
		resURI = fmt.Sprintf("https://%s/%sdatasets/%s", hosthub, dir, file)
	} else if ut.Empty(u.Host) && ut.Empty(dir) && len(file) > 0 {
		source = HubData
		pat = file
		resURI = fmt.Sprintf("https://%s/datasets/%s", hosthub, file)
	} else if ut.Empty(u.Host) && len(dir) > 0 && len(file) > 0 {
		source = UserData
		pat = u.Path
		resURI = fmt.Sprintf("https://%s/%sdatasets/%s", hosthub, dir, file)
	} else if u.Scheme == "file" || LocalPath(u.Path) {
		source = LocalData
		pat = u.Path
		resURI = pat
	} else if len(u.Host) > 0 && u.Host != hosthub {
		source = WideData
		pat = u.String()
		resURI = pat
	} else {
		source = InvalidData
	}

	return
}
