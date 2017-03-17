package v1

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

type URLBuilder struct {
	root     *url.URL
	router   *mux.Router
	relative bool
}

func NewURLBuilder(root *url.URL, relative bool) *URLBuilder {
	return &URLBuilder{
		root:     root,
		router:   Router(),
		relative: relative,
	}
}

func NewURLBuilderFromString(root string, relative bool) (*URLBuilder, error) {
	u, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	return NewURLBuilder(u, relative), nil
}

func NewURLBuilderFromRequest(r *http.Request, relative bool) *URLBuilder {
	var scheme string
	forwarded := r.Header.Get("X-Forwarded-Proto")

	switch {
	case len(forwarded) > 0:
		scheme = forwarded
	case r.TLS != nil:
		scheme = "https"
	case len(r.URL.Scheme) > 0:
		scheme = r.URL.Scheme
	default:
		scheme = "http"
	}

	host := r.Host
	forwardedHost := r.Header.Get("X-Forwarded-Host")
	if len(forwardedHost) > 0 {
		hosts := strings.SplitN(forwardedHost, ",", 2)
		host = strings.TrimSpace(hosts[0])
	}

	basePath := routeDescriptorsMap[RouteNameBase].Path
	requestPath := r.URL.Path
	index := strings.Index(requestPath, basePath)

	u := &url.URL{
		Scheme: scheme,
		Host:   host,
	}

	if index > 0 {
		u.Path = requestPath[0 : index+1]
	}

	return NewURLBuilder(u, relative)
}

func (ub *URLBuilder) BuildBaseURL() (string, error) {
	route := ub.cloneRoute(RouteNameBase)

	baseURL, err := route.URL()
	if err != nil {
		return "", err
	}

	return baseURL.String(), nil
}

func (ub *URLBuilder) BuildModVersions(partialToken *NameVersionToken) (string, error) {
	route := ub.cloneRoute(RouteNameModVersions)
	routeUrl, err := route.URL("namespace", partialToken.Namespace, "mod", partialToken.Name)
	if err != nil {
		return "", err
	}

	return routeUrl.String(), nil
}

func (ub *URLBuilder) BuildModVersionMeta(t *NameVersionToken) (string, error) {
	route := ub.cloneRoute(RouteNameMods)
	routeUrl, err := route.URL("namespace", TokenNamespace(t), "mod", t.Name, "version", t.Version)
	if err != nil {
		return "", err
	}

	return routeUrl.String(), nil
}

func (ub *URLBuilder) BuildModsRegistry() (string, error) {
	route := ub.cloneRoute(RouteNameMods)

	routeUrl, err := route.URL()
	if err != nil {
		return "", err
	}

	return routeUrl.String(), nil
}

func (ub *URLBuilder) BuildProfilesRegistry() (string, error) {
	route := ub.cloneRoute(RouteNameProfiles)

	routeUrl, err := route.URL()
	if err != nil {
		return "", err
	}

	return routeUrl.String(), nil
}

type clonedRoute struct {
	*mux.Route

	root     *url.URL
	relative bool
}

func (ub *URLBuilder) cloneRoute(name string) clonedRoute {
	route := new(mux.Route)
	root := new(url.URL)

	*route = *ub.router.GetRoute(name)
	*root = *ub.root
	return clonedRoute{Route: route, root: root, relative: ub.relative}
}

func (cr clonedRoute) URL(pairs ...string) (*url.URL, error) {
	routeURL, err := cr.Route.URL(pairs...)
	if err != nil {
		return nil, err
	}

	if cr.relative {
		return routeURL, nil
	}

	if routeURL.Scheme == "" && routeURL.User == nil && routeURL.Host == "" {
		routeURL.Path = routeURL.Path[1:]
	}

	url := cr.root.ResolveReference(routeURL)
	url.Scheme = cr.root.Scheme
	return url, nil
}
