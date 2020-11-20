package helm

// type helmAgent interface {
// 	getReleases(*action.Configuration, string) ([]*release.Release, error)
// 	upgradeRelease(*action.Configuration, string, *poseidon.File, string, map[string]interface{}) (*release.Release, error)
// }

// type defaultHelmAgent struct{}

// func (h defaultHelmAgent) getReleases(cfg *action.Configuration, cluster string) ([]*release.Release, error) {
// 	list := action.NewList(cfg)
// 	list.All = true
// 	list.AllNamespaces = true
// 	list.Limit = 0
// 	list.SetStateMask()

// 	releases, err := list.Run()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return releases, nil
// }

// func (h defaultHelmAgent) upgradeRelease(cfg *action.Configuration, n string, f *poseidon.File, ns string, vals map[string]interface{}) (*release.Release, error) {
// 	chart, err := loader.LoadArchive(bytes.NewReader(f.Contents))
// 	if err != nil {
// 		return nil, err
// 	}

// 	upgrade := action.NewUpgrade(cfg)
// 	upgrade.Namespace = ns
// 	return upgrade.Run(n, chart, vals)
// }

// func toHelmClient(c *Cluster, namespace string) (*action.Configuration, error) {
// 	ep := c.Endpoint()
// 	tok := c.Token()
// 	flags := &genericclioptions.ConfigFlags{
// 		APIServer:   &ep,
// 		BearerToken: &tok,
// 		// TODO AdamP - flags only supports cert files, how do we want to handle?
// 		// CertFile:    &cluster.Cert,
// 		Insecure: &[]bool{true}[0],
// 	}
// 	actionConfig := new(action.Configuration)
// 	if err := actionConfig.Init(flags, namespace, "secrets", helmLogger); err != nil {
// 		return nil, err
// 	}

// 	return actionConfig, nil
// }

// func helmLogger(template string, args ...interface{}) {
// 	logger.Info(fmt.Sprintf(template, args...))
// }

// // DownloadChart downloads a chart to the local filesystem
// func (a DefaultAgent) DownloadChart(name, version, out, endpoint string, opts *action.ChartPathOptions) (string, error) {
// 	pull := action.NewPull()
// 	pull.ChartPathOptions = *opts
// 	pull.Settings = cli.New()
// 	pull.RepoURL = endpoint
// 	pull.Version = version

// 	// TODO AdamP - we may want to implement some kind of directory management to prevent overflow here..
// 	pull.DestDir = out

// 	if _, err := pull.Run(name); err != nil {
// 		return "", err
// 	}

// 	return path.Join(out, getFullName(name, version)), nil
// }

// // ToChartPathOptions returns the chart path options for helm
// func (r *Repo) ToChartPathOptions() *action.ChartPathOptions {
// 	return &action.ChartPathOptions{
// 		InsecureSkipTLSverify: true,
// 		RepoURL:               r.Endpoint,
// 	}
// }

// // ToChartPathOptions returns the chart path options for helm
// func (r *Repo) ToChartPathOptions() *action.ChartPathOptions {
// 	return &action.ChartPathOptions{
// 		InsecureSkipTLSverify: true,
// 		RepoURL:               r.Endpoint,
// 	}
// }
