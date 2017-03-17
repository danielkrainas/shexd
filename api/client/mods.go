package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/shexd/api/v1"
)

type ModsApi interface {
	//SearchMods() ([]*v1.ModInfo, error)
	CreateMod(m *v1.ModInfo) (*v1.ModInfo, error)
	GetModInfo(t *v1.NameVersionToken) (*v1.ModInfo, error)
}

type modsApi struct {
	*Client
}

func (c *Client) Mods() ModsApi {
	return &modsApi{c}
}

func (api *modsApi) GetModInfo(t *v1.NameVersionToken) (*v1.ModInfo, error) {
	url, err := api.urls().BuildModVersionMeta(t)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m := &v1.ModInfo{}
	if err = json.Unmarshal(body, m); err != nil {
		return nil, err
	}

	return m, nil
}

/*func (api *modsApi) SearchMods() ([]*v1.ModInfo, error) {
	url, err := api.urls().BuildModsRegistry()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := api.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	p := make([]*v1.ModInfo, 0)
	if err = json.Unmarshal(body, &p); err != nil {
		return nil, err
	}

	return p, nil
}*/

func (api *modsApi) CreateMod(m *v1.ModInfo) (*v1.ModInfo, error) {
	body, err := json.Marshal(&m)
	if err != nil {
		return nil, err
	}

	url, err := api.urls().BuildModsRegistry()
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := api.do(r)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	m = &v1.ModInfo{}
	if err = json.Unmarshal(body, m); err != nil {
		return nil, err
	}

	return m, nil
}

/*
import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/danielkrainas/shexd/api/v1"
)

const (
	ApiProfilesPath       = "profiles"
	ApiModsPath           = "mods"
	ApiModsMetaPathSuffix = "meta"
)

func downloadContents(url string) ([]byte, error) {
	resp, err := http.Get(url)
	contents, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return contents, err
	// smarter to io.Copy to local (temp?) file in case mods happen to be large
}

func DownloadAndCacheContent(url string, cachePath string) (int64, error) {
	resp, err := http.Get(url)
	var read int64
	if err == nil {
		cacheFile, err := os.Create(cachePath)
		if err == nil {
			defer resp.Body.Close()
			defer cacheFile.Close()
			read, err = io.Copy(cacheFile, resp.Body)
		}
	}

	return read, err
}

func PostContent(url string, bodyContent []byte) ([]byte, error) {
	body := bytes.NewBuffer(bodyContent)
	res, err := http.Post(url, "application/json", body)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func DownloadModInfo(source string, mod *v1.NameVersionToken) (*v1.RemoteModInfo, error) {
	// NOTE: expect remote name to be something like user/package

	// http://somesite.com/mods/admin/cool-package/v/0.1.2/meta
	url := source + path.Join(ApiModsPath, mod.Name, "v", mod.Version, ApiModsMetaPathSuffix)
	contents, err := downloadContents(url)
	info := &v1.RemoteModInfo{}
	info.Source = url
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(contents, info)
	return info, err
}

func DownloadModVersionList(source string, modName string) ([]string, error) {
	// http://somesite.com/mods/admin/cool-package/v
	url := source + path.Join(ApiModsPath, modName, "v")
	contents, err := downloadContents(url)
	if err != nil {
		return nil, err
	}

	versions := make([]string, 0)
	err = json.Unmarshal(contents, &versions)
	return versions, err
}

func DownloadMod(source string, destPath string, info *v1.RemoteModInfo) error {
	// http://somesite.com/mods/admin/cool-package/v/0.1.2
	url := source + path.Join(ApiModsPath, info.Name, "v", info.Version)

	_, err := DownloadAndCacheContent(url, destPath)
	if err != nil {
		return err
	}

	return nil
}

func DownloadProfileAsLocal(source *v1.ProfileSource, localName string) (*v1.Profile, error) {
	rp, err := DownloadProfile(source)
	if err != nil {
		return nil, err
	}

	return v1.MakeLocalProfile(localName, rp), nil
}

func DownloadProfile(source *v1.ProfileSource) (*v1.RemoteProfile, error) {
	url := path.Join(source.Location, ApiProfilesPath, source.Uid)
	jsonContent, err := downloadContents(url)
	if err != nil {
		return nil, err
	}

	remoteProfile := v1.NewRemoteProfile(source)
	err = json.Unmarshal(jsonContent, &remoteProfile)
	return remoteProfile, err
}
*/
