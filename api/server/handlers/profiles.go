package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/danielkrainas/gobag/api/errcode"
	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"

	"github.com/danielkrainas/shexd/actions"
	"github.com/danielkrainas/shexd/api/v1"
	"github.com/danielkrainas/shexd/queries"
)

func Profiles(actionPack actions.Pack) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			PublishProfile(actionPack, w, r)
		case http.MethodGet:
			SearchProfiles(actionPack, w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

func PublishProfile(c cqrs.CommandHandler, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := acontext.GetLogger(ctx)
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	rp := &v1.RemoteProfile{}
	if err = json.Unmarshal(body, rp); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	rp.Source = nil
	if err := c.Handle(ctx, &commands.StoreProfile{New: true, Profile: rp}); err != nil {
		log.Error(err)
		acontext.TrackError(ctx, errcode.ErrorCodeUnknown.WithDetail(err))
		return
	}

	log.Infof("profile published %q", rp.Name)
	if err := v1.ServeJSON(w, m); err != nil {
		log.Errorf("error sending user json: %v", err)
	}
}

func SearchProfiles(q cqrs.QueryExecutor, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	mods, err := q.Execute(ctx, &queries.SearchProfiles{})
	if err != nil {
		acontext.GetLogger(ctx).Error(err)
		acontext.TrackError(ctx, err)
		return
	}

	if err := v1.ServeJSON(w, mods); err != nil {
		acontext.GetLogger(ctx).Errorf("error sending mods json: %v", err)
	}
}
