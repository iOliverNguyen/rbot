package genapi

const tplText = `
func init() {
	httprpc.Register(NewServer)
}

func NewServer(builder interface{}, hooks ...httprpc.HooksBuilder) (httprpc.Server, bool) {
	switch builder := builder.(type) {
	{{range $s := .Services -}}
	case func() {{.Name}}Service:
		return New{{.Name}}ServiceServer(builder, hooks...), true
	case {{.Name}}Service:
		fn := func() {{.Name}}Service { return builder }
		return New{{.Name}}ServiceServer(fn, hooks...), true
	{{end -}}
	default:
	{{if .Services|len|eq 0 -}}
		_ = builder
	{{end -}}
		return nil, false
	}
}

{{range $s := .Services}}
type {{.Name}}ServiceServer struct {
	hooks   httprpc.HooksBuilder
	builder func() {{.Name}}Service
}

func New{{.Name}}ServiceServer(builder func() {{.Name}}Service, hooks ... httprpc.HooksBuilder) httprpc.Server{
	return &{{.Name}}ServiceServer {
		hooks: httprpc.ChainHooks(hooks...),
		builder: builder,
	}
}

const {{.Name}}ServicePathPrefix = "/{{.APIPath}}/"

{{range $m := .Methods -}}
const Path_{{$s.Name}}_{{$m.Name}} = "/{{$s.APIPath}}/{{.APIPath}}"
{{end}}

func (s *{{$s.Name}}ServiceServer) PathPrefix() string {
	return {{.Name}}ServicePathPrefix
}

func (s *{{$s.Name}}ServiceServer) WithHooks(hooks httprpc.HooksBuilder) httprpc.Server {
	result := *s
	result.hooks = httprpc.ChainHooks(s.hooks, hooks)
	return &result
}

func (s *{{$s.Name}}ServiceServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	hooks := httprpc.WrapHooks(s.hooks)
	ctx, info := req.Context(), &httprpc.HookInfo{Route: req.URL.Path, HTTPRequest: req}
	ctx, err := hooks.RequestReceived(ctx, *info)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	serve, err := httprpc.ParseRequestHeader(req)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	reqMsg, exec, err := s.parseRoute(req.URL.Path, hooks, info)
	if err != nil {
		httprpc.WriteError(ctx, resp, hooks, *info, err)
		return
	}
	serve(ctx, resp, req, hooks, info, reqMsg, exec)
}

func (s *{{$s.Name}}ServiceServer) parseRoute(path string, hooks httprpc.Hooks, info *httprpc.HookInfo) (reqMsg httprpc.Message, _ httprpc.ExecFunc, _ error) {
	switch path {
{{range $m := .Methods -}}
	case "/{{$s.APIPath}}/{{.APIPath}}":
	msg := {{(index .Request.Items 0).Type|new}}
	fn := func(ctx context.Context) (newCtx context.Context, resp httprpc.Message, err error) {
		inner := s.builder()
		info.Request, info.Inner = msg, inner
		newCtx, err = hooks.RequestRouted(ctx, *info)
		if err != nil {
			return
		}
		resp, err = inner.{{.Name}}(newCtx, msg)
		return
	}
	return msg, fn, nil
{{end -}}
	default:
		msg := fmt.Sprintf("no handler for path %q", path)
		return nil, nil, httprpc.BadRouteError(msg, "POST", path)
	}
}
{{end}}
`
