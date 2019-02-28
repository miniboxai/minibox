package apiserver

import "html/template"

const mainNavTpl = `
<ul class="main-nav">
	<li>
		<a href="/blog">博客</a>
	</li>
	<li>
		<a href="/docs">文档</a>
	</li>
	<li>
		<a href="/dev_docs">开发文档</a>
	</li>
</ul>
`

func buildNavTemplate() *template.Template {
	return template.Must(template.New("main-nav.tmpl").Parse(mainNavTpl))
}
