package main

/**
Para usar o blog, crie um novo arquivo markdown. Por convenção, o arquivo e o slug devem ter o mesmo nome. O arquivo do tipo .md deve estar em /markdown

Para instruções mais detalhadas, ler o blog do criador original em https://fluxsec.red,
ou entre em contato pelo twitter https://twitter.com/0xfluxsec.

O arquivo .md deve conter o seguinte cabeçalho, incluindo os 3 traços.
Eles são utilizados para separar o cabeçalho do conteudo.

Title: Titulo da página
Slug: slug-da-url
Parent: O nome que você quer dar à série de posts
Order: numero em termo da ordem no parente
Description: Curta descrição que aparece abaixo do título
MetaPropertyTitle: Titulo para compartilhamento
MetaDescription: Descrição de +/- 150 a 200 palavras para o SEO da página.
MetaPropertyDescription: descrição curta para redes sociais.
MetaOgURL: https://www.[[SEU-DOMINIO-AQUI]]/slug-da-url
---
Conteudo

Downloads adicionais:
* FontAwesome free, em /static/
Em caso de dúvidas, verificar o arquivo /markdown/index.md
*/

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type BlogPost struct {
	Title                   string
	Slug                    string
	Parent                  string
	Content                 template.HTML
	Description             string
	Order                   int
	Headers                 []string // São os <h2> ou ## de cada post
	MetaDescription         string
	MetaPropertyTitle       string
	MetaPropertyDescription string
	MetaOgURL               string
}

type SidebarData struct {
	Categories []Category
}

type Category struct {
	Name  string
	Pages []BlogPost
	Order int
}

var BaseURL = "http://mahzze.dev"

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// dados da sidebar (barra lateral)
	sidebarData, err := loadSidebarData("./markdown")
	if err != nil {
		log.Fatal(err) // Se não for possivel carregar os posts, para o site
	}

	// Registra o template da sidebar como um partial
	r.SetFuncMap(template.FuncMap{
		"loadSidebar": func() SidebarData {
			return sidebarData
		},
		"dict": dict,
	})

	// carrega os templates
	r.LoadHTMLGlob("templates/*")

	// carrega arquivos estáticos - como imagens, arquivos css, fontes, etc.
	r.Static("/static", "./static")

	// carrega os arquivos markdown (posts) e "traduz" para HTML.
	posts, err := loadMarkdownPosts("./markdown")
	if err != nil {
		log.Fatal(err) // Se não for possivel carregar nada, causa um erro
	}

	// Roteamento da página principal
	r.GET("/", func(c *gin.Context) {
		indexPath := "./markdown/index.md"
		indexContent, err := os.ReadFile(indexPath)
		if err != nil {
			log.Printf("Erro durante operacao: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
			return
		}

		post, err := parseMarkdownFile(indexContent)
		if err != nil {
			log.Printf("Erro durante operacao: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro interno do servidor"})
			return
		}

		sidebarLinks := createSidebarLinks(post.Headers)

		c.HTML(http.StatusOK, "index.html", gin.H{
			"Title":                   post.Title,
			"Content":                 post.Content,
			"SidebarData":             sidebarData,
			"Headers":                 post.Headers,
			"SidebarLinks":            sidebarLinks,
			"CurrentSlug":             post.Slug,
			"MetaDescription":         post.MetaDescription,
			"MetaPropertyTitle":       post.MetaPropertyTitle,
			"MetaPropertyDescription": post.MetaPropertyDescription,
			"MetaOgURL":               post.MetaOgURL,
		})
	})

	// roteamento de cada post, com base no slug (resultado: posts em /{{slug}})
	for _, post := range posts {
		localPost := post
		if localPost.Slug != "" {
			sidebarLinks := createSidebarLinks(localPost.Headers)
			r.GET("/"+localPost.Slug, func(c *gin.Context) {
				c.HTML(http.StatusOK, "layout.html", gin.H{
					"Title":                   localPost.Title,
					"Content":                 localPost.Content,
					"SidebarData":             sidebarData,
					"Headers":                 localPost.Headers,
					"Description":             localPost.Description,
					"SidebarLinks":            sidebarLinks,
					"CurrentSlug":             localPost.Slug,
					"MetaDescription":         localPost.MetaDescription,
					"MetaPropertyTitle":       localPost.MetaPropertyTitle,
					"MetaPropertyDescription": localPost.MetaPropertyDescription,
					"MetaOgURL":               localPost.MetaOgURL,
				})
			})
		} else {
			log.Printf("AVISO: Post '%s' está com o slug vazio, e não vai ser acessável via uma única URL.\n", localPost.Title)
		}
	}

	r.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", gin.H{
			"Title": "Não existe",
		})
	})

	r.Run()
}

func loadMarkdownPosts(dir string) ([]BlogPost, error) {
	var posts []BlogPost
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".md") {
			path := dir + "/" + file.Name()
			content, err := os.ReadFile(path)
			if err != nil {
				return nil, err
			}

			post, err := parseMarkdownFile(content)
			if err != nil {
				return nil, err
			}

			posts = append(posts, post)
		}
	}

	return posts, nil
}

func parseMarkdownFile(content []byte) (BlogPost, error) {
	sections := strings.SplitN(string(content), "---", 2)
	if len(sections) < 2 {
		return BlogPost{}, errors.New("o texto markdown não pôde ser formatado")
	}

	metadata := sections[0]
	mdContent := sections[1]

	// lida com \r (carriage returns)
	metadata = strings.ReplaceAll(metadata, "\r", "")
	mdContent = strings.ReplaceAll(mdContent, "\r", "")

	title, slug, parent, description, order, metaDescriptionStr,
		metaPropertyTitleStr, metaPropertyDescriptionStr,
		metaOgURLStr := parseMetadata(metadata)

	htmlContent := mdToHTML([]byte(mdContent))
	headers := extractHeaders([]byte(mdContent))

	return BlogPost{
		Title:                   title,
		Slug:                    slug,
		Parent:                  parent,
		Description:             description,
		Content:                 template.HTML(htmlContent),
		Headers:                 headers,
		Order:                   order,
		MetaDescription:         metaDescriptionStr,
		MetaPropertyTitle:       metaPropertyTitleStr,
		MetaPropertyDescription: metaPropertyDescriptionStr,
		MetaOgURL:               metaOgURLStr,
	}, nil
}

func extractHeaders(content []byte) []string {
	var headers []string
	//match only level 2 markdown headers
	re := regexp.MustCompile(`(?m)^##\s+(.*)`)
	matches := re.FindAllSubmatch(content, -1)

	for _, match := range matches {
		// match[1] contains header text without the '##'
		headers = append(headers, string(match[1]))
	}

	return headers
}

func mdToHTML(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	parser := parser.NewWithExtensions(extensions)

	opts := html.RendererOptions{
		Flags: html.CommonFlags | html.HrefTargetBlank,
	}
	renderer := html.NewRenderer(opts)

	doc := parser.Parse(md)

	output := markdown.Render(doc, renderer)

	return output
}

func parseMetadata(metadata string) (
	title string,
	slug string,
	parent string,
	description string,
	order int,
	metaDescription string,
	metaPropertyTitle string,
	metaPropertyDescription string,
	metaOgURL string,
) {
	re := regexp.MustCompile(`(?m)^(\w+):\s*(.+)`)
	matches := re.FindAllStringSubmatch(metadata, -1)

	metaDataMap := make(map[string]string)
	for _, match := range matches {
		if len(match) == 3 {
			metaDataMap[match[1]] = match[2]
		}
	}

	title = metaDataMap["Title"]
	slug = metaDataMap["Slug"]
	parent = metaDataMap["Parent"]
	description = metaDataMap["Description"]
	orderStr := metaDataMap["Order"]
	metaDescriptionStr := metaDataMap["MetaDescription"]
	metaPropertyTitleStr := metaDataMap["MetaPropertyTitle"]
	metaPropertyDescriptionStr := metaDataMap["MetaPropertyDescription"]
	metaOgURLStr := metaDataMap["MetaOgURL"]

	orderStr = strings.TrimSpace(orderStr)
	order, err := strconv.Atoi(orderStr)
	if err != nil {
		log.Printf("Erro convertendo ordem a partir de string: %v", err)
		order = 9999 // por segurança, deixar este valor bem alto
	}

	return title, slug, parent, description, order, metaDescriptionStr,
		metaPropertyTitleStr, metaPropertyDescriptionStr, metaOgURLStr
}

func loadSidebarData(dir string) (SidebarData, error) {
	var sidebar SidebarData
	categoriesMap := make(map[string]*Category)

	posts, err := loadMarkdownPosts(dir)
	if err != nil {
		return sidebar, err
	}

	for _, post := range posts {
		if post.Parent != "" {
			if _, exists := categoriesMap[post.Parent]; !exists {
				categoriesMap[post.Parent] = &Category{
					Name:  post.Parent,
					Pages: []BlogPost{post},
					Order: post.Order,
				}
			} else {
				categoriesMap[post.Parent].Pages = append(categoriesMap[post.Parent].Pages, post)
			}
		}
	}

	// Transforma mapa em vetor/arranjo
	for _, cat := range categoriesMap {
		sidebar.Categories = append(sidebar.Categories, *cat)
	}

	// organiza categorias por ordem
	sort.Slice(sidebar.Categories, func(i, j int) bool {
		return sidebar.Categories[i].Order < sidebar.Categories[j].Order
	})

	return sidebar, nil
}

func createSidebarLinks(headers []string) template.HTML {
	var linksHTML string
	for _, header := range headers {
		sanitizedHeader := sanitizeHeaderForID(header)
		link := fmt.Sprintf(`<li><a href="#%s">%s</a></li>`, sanitizedHeader, header)
		linksHTML += link
	}
	return template.HTML(linksHTML)
}

func sanitizeHeaderForID(header string) string {
	// passa para letras minusculas
	header = strings.ToLower(header)

	// troca espaços por hífens
	header = strings.ReplaceAll(header, " ", "-")

	// remover todos os caracteres que não são letras, números ou hífens
	header = regexp.MustCompile(`[^a-z0-9\-]`).ReplaceAllString(header, "")

	return header
}

func dict(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("chamada inválida de dicionario")
	}
	dict := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("as chaves do dicionario devem ser strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}
