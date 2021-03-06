package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestFindFileIdentifier_Success(t *testing.T) {
	var testCases = []struct {
		name            string
		memberType      string
		patternToSearch string
		fakeHTMLFile    string
		desiredOutput   string
	}{
		{"Should get right code for membros ativos for years differente of 2017", "membrosAtivos", ":membros-ativos-02-2019", ".../4554/resource-fevereiro:download=5051:membros-ativos-02-2019", "5051"},
		{"Should get right code for membros ativos for year 2017", "membrosAtivos", ":quadro-de-membros-ativos-fevereiro-2017", ".../4554/resource-online:download=5051:quadro-de-membros-ativos-fevereiro-2017", "5051"},
		{"Should get right code for membros inativos for year differente of 2014", "membrosInativos", ":membros-inativos-03-2017", ".../4554/resource:download=4312:membros-inativos-03-2017", "4312"},
		{"Should get right code for membros inativos for year 2014 and month different of janeiro", "membrosInativos", ":membros-inativos-03-2014", ".../4554/resource:download=1234:membros-inativos-03-2014", "1234"},
		{"Should get right code for membros inativos for year 2014 and month janeiro", "membrosInativos", ":membros-inativos-01-2015", ".../4554/resource:download=4567:membros-inativos-01-2015", "4567"},
		{"Should get right code for servidores ativos", "servidoresAtivos", ":servidores-ativos-01-2015", ".../31342sas2/endpoint:download=9999:servidores-ativos-01-2015", "9999"},
		{"Should get right code for servidores inativos", "servidoresInativos", ":servidores-inativos-01-2015", ".../ghytr6/resource:download=1098:servidores-inativos-01-2015", "1098"},
		{"Should get right code for pensionistas", "pensionistas", ":pensionistas-01-2015", ".../5tghjuw2/Controller:random=5453:pensionistas-01-2015", "5453"},
		{"Should get right code for colaboradores", "colaboradores", ":contracheque-valores-percebidos-colaboradores-fevereiro", ".../random/servlet:code=3490:contracheque-valores-percebidos-colaboradores-fevereiro", "3490"},
		{"Should get right code for exercicios anteriores", "exerciciosAnteriores", ":dea-022019", ".../controller_servlet:download=5378:dea-022019", "5378"},
		{"Should get right code for indenizacoes e outros pagamentos", "indenizacoesEOutrosPagamentos", ":virt-abril-2019", ".../members_controller:code=8712:virt-abril-2019", "8712"},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			code, _ := findFileIdentifier(tt.memberType, tt.fakeHTMLFile, tt.patternToSearch)
			if code != tt.desiredOutput {
				t.Errorf("got %s, want %s", code, tt.desiredOutput)
			}
		})
	}
}

func TestPathResolver_Sucess(t *testing.T) {
	var testCases = []struct {
		name   string
		month  int
		year   int
		member string
		out    string
	}{
		{"should get path for membros ativos for year differente of 2017", 2, 2019, "remuneracao-de-todos-os-membros-ativos", ":membros-ativos-02-2019"},
		{"should get path for membros ativos for year 2017", 2, 2017, "remuneracao-de-todos-os-membros-ativos", ":quadro-de-membros-ativos-fevereiro-2017"},
		{"should get path for membros inativos for different year of 2014 and month different of january", 2, 2017, "proventos-de-todos-os-membros-inativos", ":membros-inativos-02-2017"},
		{"should get path for membros inativos for year 2014 and month different of january", 2, 2014, "proventos-de-todos-os-membros-inativos", ":membros-inativos-02-2015"},
		{"should get path for membros inativos for year 2014 and month january", 1, 2014, "proventos-de-todos-os-membros-inativos", ":membros-inativos-01-2014"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			path := pathResolver(tt.month, tt.year, tt.member)
			if path != tt.out {
				t.Errorf("got %s, want %s", path, tt.out)
			}
		})
	}
}

func TestPathResolver_Error(t *testing.T) {
	var testCases = []struct {
		name   string
		month  int
		year   int
		member string
		out    string
	}{
		{"should fail for membros ativos for year differente of 2017", 2, 2019, "remuneracao-de-todos-os-membros-ativos", ":quadro-de-membros-ativos-fevereiro-2019"},
		{"should fail for membros ativos for year 2017", 2, 2017, "remuneracao-de-todos-os-membros-ativos", ":membros-ativos-2-2017"},
		{"should fail for membros inativos for year different of 2014 and month different of january", 2, 2017, "proventos-de-todos-os-membros-inativos", ":membros-inativos-02-2015"},
		{"should fail for membros inativos for year 2014 and month different of january", 2, 2014, "proventos-de-todos-os-membros-inativos", ":membros-inativos-01-2015"},
		{"should fail for membros inativos for year 2014 and month january", 1, 2014, "proventos-de-todos-os-membros-inativos", ":membros-inativos-02-2014"},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			path := pathResolver(tt.month, tt.year, tt.member)
			if path == tt.out {
				t.Errorf("got %s, want %s", path, tt.out)
			}
		})
	}
}

func TestCrawl(t *testing.T) {
	var testCases = []struct {
		name           string
		outputPath     string
		month          int
		year           int
		crawlingServer *httptest.Server
		out            map[string]struct{}
	}{
		{"should get 8 files for february of 2019", "files", 2, 2019,
			httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w,
					"<div><link src=\".../4554/resource-fevereiro:download=5051:membros-ativos-02-2019\"/></div>"+
						"<div><link src=\".../4554/resource:download=4312:membros-inativos-02-2019\"/></div>"+
						"<div><link src=\".../31342sas2/endpoint:download=9999:servidores-ativos-02-2019\"/></div>"+
						"<div><link src=\".../ghytr6/resource:download=1098:servidores-inativos-02-2019\"/></div>"+
						"<div><link src=\".../5tghjuw2/Controller:random=5453:pensionistas-02-2019\"/></div>"+
						"<div><link src=\".../random/servlet:code=3490:contracheque-valores-percebidos-colaboradores-fevereiro\"/></div>"+
						"<div><link src=\".../controller_servlet:download=5378:dea-022019\"/></div>"+
						"<div><link src=\".../members_controller:code=8712:virt-fevereiro-2019\"/></div>")
			})),
			map[string]struct{}{
				"files/proventos-de-todos-os-membros-inativos-02-2019.xlsx":                  {},
				"files/proventos-de-todos-os-servidores-inativos-02-2019.xlsx":               {},
				"files/remuneracao-de-todos-os-membros-ativos-02-2019.xlsx":                  {},
				"files/remuneracao-de-todos-os-servidores-atuvos-02-2019.xlsx":               {},
				"files/valores-percebidos-por-todos-os-colaboradores-02-2019.xlsx":           {},
				"files/valores-percebidos-por-todos-os-pensionistas-02-2019.xlsx":            {},
				"files/verbas-indenizatorias-e-outras-remuneracoes-temporarias-02-2019.xlsx": {},
				"files/verbas-referentes-a-exercicios-anteriores-02-2019.xlsx":               {},
			}},
	}
	var pathsReference []string
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			outs, _ := Crawl(tt.outputPath, tt.month, tt.year, tt.crawlingServer.URL+"/")
			pathsReference = outs
			for _, out := range outs {
				_, ok := tt.out[out]
				if !ok {
					t.Errorf("got %s and it is not present on list", out)
				}
			}
		})
	}
	for _, path := range pathsReference {
		err := os.Remove(path)
		if err != nil {
			panic("failed to remove file at path: " + path)
		}
	}
}
