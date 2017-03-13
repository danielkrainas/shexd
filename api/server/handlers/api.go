package handlers

import (
	"net/http"

	"github.com/danielkrainas/gobag/context"
	"github.com/gorilla/mux"

	"github.com/danielkrainas/shex/api/v1"
	"github.com/danielkrainas/shex/registry/actions"
)

const ServerVersionHeader = "Shex-Registry-Version"
const ApiVersionHeader = "Shex-Api-Version"
const ApiVersion = "1"

type Api struct {
	router *mux.Router
}

func NewApi(actionPack actions.Pack) (*Api, error) {
	api := &Api{
		router: v1.RouterWithPrefix(""),
	}

	api.register(v1.RouteNameBase, Base)
	api.register(v1.RouteNameMods, Mods(actionPack))

	return api, nil
}

func (api *Api) register(routeName string, dispatch http.HandlerFunc) {
	api.router.GetRoute(routeName).Handler(api.dispatcher(dispatch))
}

func (api *Api) dispatcher(dispatch http.HandlerFunc) http.Handler {
	return http.Handler(dispatch)
}

func (api *Api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Add(ServerVersionHeader, acontext.GetVersion(r.Context()))
	w.Header().Add(ApiVersionHeader, ApiVersion)
	api.router.ServeHTTP(w, r)
}

/* API endpoint
app.post('/profiles/:user/:profileId', function (req, res) {
	console.log('profile posted:');
	console.log(JSON.stringify(req.body, null, 4));
	res.send('1');// return current version number of profile
});

app.head('/profiles/:user/:profileId', function (req, res) {
	//If-Modified-Since: version
});

app.get('/profiles/:user/:profileId', function (req, res) {
	var name = req.params.user + '/' + req.params.profileId;
	console.log('profile request:', name);
	res.json({
		rev: 1,
		name: name,
		mods: {},
	});
});

app.get('/mods/:user/:mod/v', function (req, res) {
	res.json(['0.1.8', '1.0.0']);
});

app.get('/mods/:user/:mod/v/:version/meta', function (req, res) {
	if (req.params.version === 'latest') {
		req.params.version = '1.0.0';
	}

	console.info('params ', req.params);
	var result = {
		name: req.params.user + '/' + req.params.mod,
		version: req.params.version
	};

	console.info('meta ', result);
	res.json(result);
});

app.get('/mods/:user/v/:mod/:version', function (req, res) {
	if (req.params.version === 'latest') {
		req.params.version = '1.0.0';
	}

	console.info('params ', req.params);

	res.sendFile(path.join(__dirname, '/candyland2.smod'), function (err) {
		if (err) {
			console.error(err);
			return;
		}
	});
});*/
