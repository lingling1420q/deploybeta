package version

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"

	etcdv3 "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"

	"github.com/jysperm/deploying/lib/builder"
	"github.com/jysperm/deploying/lib/etcd"
	appModel "github.com/jysperm/deploying/lib/models/app"
)

const DefaultRegistry = "localhost:5000"

type Version struct {
	Shasum   string `json:"shasum"`
	Tag      string `json:"tag"`
	Registry string `json:"registry"`
}

func CreateVersion(app *appModel.Application, registry string, gitTag string) (Version, error) {
	version := generateTag()

	var nameVersion string
	var newVersion Version
	if registry == "" {
		registry = DefaultRegistry
	}
	nameVersion = fmt.Sprintf("%s/%s:%s", registry, app.Name, version)
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, version)

	buildOpts := types.ImageBuildOptions{
		Tags: []string{nameVersion},
	}
	fmt.Println(app)
	shasum, err := builder.BuildImage(buildOpts, app.GitRepository, gitTag)
	if err != nil {
		return Version{}, err
	}

	if err := builder.PushImage(nameVersion); err != nil {
		return Version{}, err
	}

	newVersion.Shasum = shasum
	newVersion.Tag = version
	newVersion.Registry = registry

	jsonVersion, err := json.Marshal(newVersion)
	if err != nil {
		return Version{}, err
	}
	if _, err := etcd.Client.Put(context.Background(), versionKey, string(jsonVersion)); err != nil {
		return Version{}, err
	}

	return newVersion, nil
}

func DeleteVersion(app appModel.Application, version string) error {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, version)

	if _, err := etcd.Client.Delete(context.Background(), versionKey); err != nil {
		return err
	}

	return nil
}

func FindByTag(app appModel.Application, tag string) (*Version, error) {
	versionKey := fmt.Sprintf("/apps/%s/versions/%s", app.Name, tag)

	resp, err := etcd.Client.Get(context.Background(), versionKey)
	if err != nil {
		return nil, err
	}

	var v Version
	if err := json.Unmarshal(resp.Kvs[0].Value, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func ListAll(app appModel.Application) (*[]Version, error) {
	versionPrefix := fmt.Sprintf("/apps/%s/versions/", app.Name)
	resp, err := etcd.Client.Get(context.Background(), versionPrefix, etcdv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var versionArray []Version
	for _, ev := range resp.Kvs {
		temp := Version{}
		_ = json.Unmarshal(ev.Value, &temp)
		fmt.Println(temp)
		versionArray = append(versionArray, temp)
	}

	return &versionArray, nil
}

func generateTag() string {
	now := time.Now()
	return fmt.Sprintf("%d%d%d-%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
}
