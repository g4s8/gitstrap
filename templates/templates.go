package templates

import (
	"github.com/g4s8/gitstrap/context"
)

// TemplatesV1 - first version of templates
type TemplatesV1 []TemplateV1

// Upgrade v1 templates to latest version
func (t TemplatesV1) Upgrade(params map[string]string) Templates {
	res := make([]Template, len(t), len(t))
	for i, old := range t {
		old.upgrade(params, &res[i])
	}
	return Templates(res)
}

// Templates - gitstrap templates
type Templates []Template

// TemplateV1 - gitstrap template
type TemplateV1 struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	URL      string `yaml:"url"`
}

func (t *TemplateV1) upgrade(params map[string]string, out *Template) {
	out.Name = t.Name
	out.Location = t.Location
	out.URL = t.URL
	out.Args = params
}

// Template - gitstrap template
type Template struct {
	Name     string            `yaml:"name"`
	Location string            `yaml:"location"`
	URL      string            `yaml:"url"`
	Args     map[string]string `yaml:"args"`
}

func (tps Templates) Apply(ctx *context.Context) error {
	return nil
}

// func (strap *strapCtx) applyTemplates(repo *github.Repository) error {
// 	// apply templates
// 	tctx := &templateContext{repo, &strap.cfg.Gitstrap}
// 	for _, t := range strap.cfg.Gitstrap.Templates {
// 		tpl := template.New(t.Name)
// 		var data []byte
// 		var err error
// 		if t.Location != "" {
// 			data, err = readTemplate(t.Location)
// 			if err != nil {
// 				return err
// 			}
// 		} else if t.URL != "" {
// 			data, err = downloadTemplate(t.URL)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if _, err = tpl.Parse(string(data)); err != nil {
// 			return fmt.Errorf("failed to parse template %s: %s", tpl.Name(), err)
// 		}
// 		fout, err := os.Create(t.Name)
// 		if err != nil {
// 			return fmt.Errorf("failed to open output file for template %s: %s", tpl.Name(), err)
// 		}
// 		if err = tpl.Execute(fout, tctx); err != nil {
// 			return fmt.Errorf("failed to execute template %s: %s", tpl.Name(), err)
// 		}
// 		fmt.Printf("Template %s applied\n", tpl.Name())
// 	}
// 	return nil
// }

// func readTemplate(name string) ([]byte, error) {
// 	tf, err := os.Open(name)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to open template file %s: %s", name, err)
// 	}
// 	data, err := ioutil.ReadAll(bufio.NewReader(tf))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read template file %s: %s", name, err)
// 	}
// 	if err = tf.Close(); err != nil {
// 		return nil, fmt.Errorf("failed to close template file %s: %s", name, err)
// 	}
// 	return data, nil
// }

// func downloadTemplate(url string) ([]byte, error) {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to download template from %s: %s", url, err)
// 	}
// 	data, err := ioutil.ReadAll(bufio.NewReader(resp.Body))
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read template body from %s: %s", url, err)
// 	}
// 	if err := resp.Body.Close(); err != nil {
// 		return nil, fmt.Errorf("failed to close connection from %s: %s", url, err)
// 	}
// 	return data, nil
// }
