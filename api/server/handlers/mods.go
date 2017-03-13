package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/registry/actions"
	"github.com/danielkrainas/shex/registry/queries"
)

func Mods(actionPack actions.Pack) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateMod(actionPack, w, r)
		case http.MethodGet:
			SearchMods(actionPack, w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func CreateMod(c cqrs.CommandHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := acontext.GetLogger(ctx)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	m := &v1.ModInfo{}
	if err = json.Unmarshal(body, m); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	log.Infof("mod %s@%s created", m.Name, m.SemVersion)
	if err := v1.ServeJSON(w, m); err != nil {
		log.Errorf("error sending user json: %v", err)
	}
}

func SearchMods(q cqrs.QueryExecutor, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mods, err := q.Execute(ctx, &queries.SearchMods{})
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		acontext.TrackError(ctx, err)
		return
	}

	if err := v1.ServeJSON(w, mods); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending mods json: %v", err)
	}
}
