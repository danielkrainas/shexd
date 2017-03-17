package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/gorilla/mux"

	"github.com/danielkrainas/shexd/actions"
	"github.com/danielkrainas/shexd/api/v1"
	"github.com/danielkrainas/shexd/queries"
)

func tokenFromRoute(r *http.Request) *v1.NameVersionToken {
	vars := mux.Vars(r)
	return &v1.NameVersionToken{
		Name:      vars["mod"],
		Namespace: vars["namespace"],
		Version:   vars["version"],
	}
}

func ModMetadata(q cqrs.QueryExecutor) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}

		ctx := r.Context()
		query := &queries.FindMod{Token: tokenFromRoute(r)}
		info, err := q.Execute(ctx, query)
		if err != nil {
			acontext.GetLogger(ctx).Error(err)
			acontext.TrackError(ctx, err)
			return
		}

		if err := v1.ServeJSON(w, info); err != nil {
			acontext.GetLogger(ctx).Errorf("error sending mod info json: %v", err)
		}
	})
}

func Mods(actionPack actions.Pack) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			CreateMod(actionPack, w, r)
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
