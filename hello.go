package main

// "strconv" // string converter
// "bufio" // ler arquivo linha a linha
// "os/exec" // executar comando cmd
// "runtime" // pegar sistema operacional

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const monitoramentos = 2
const intervaloDeMonitoramento = 5 * time.Second

func main() {
	exibeIntroducao()

	for {
		exibeMenu()

		comando := lerComando()

		switch comando {
		case 1:
			iniciarMonitoramento()
		case 2:
			lerLogsTXT()
		case 3:
			limparConsole()
		case 4:
			fmt.Println("Saindo do programa, Bye Bye...")
			os.Exit(0)
		default:
			fmt.Println("Não conheço este comando.")
			os.Exit(-1)
		}
	}
}

func devolveNomeEIdade() (string, int) {
	nome := "Brendo"
	idade := 27

	return nome, idade
}

func exibeIntroducao() {
	nome, idade := devolveNomeEIdade()
	versao := 1.1
	sistemaOperacional := runtime.GOOS
	fmt.Println("Olá", nome, "Sua idade é", idade)
	fmt.Println("Este programa está na versão", versao)
	fmt.Println("Sistema operacional ", sistemaOperacional)
}

func exibeMenu() {
	fmt.Println("\n \n1- Iniciar Monitoramento")
	fmt.Println("2- Exibir Logs")
	fmt.Println("3- Limpar terminal")
	fmt.Println("4- Sair do programa")
}

func lerComando() int {
	var comandoLido int

	_, err := fmt.Scan(&comandoLido)
	if err != nil {
		return 0
	}

	fmt.Println("O Comando escolhido foi", comandoLido, "\n ")

	return comandoLido
}

func iniciarMonitoramento() {
	fmt.Println("---------------------")
	fmt.Println("Monitorando...")
	fmt.Println("---------------------")

	sites := lerSitesTXT()

	for i := 0; i < monitoramentos; i++ {
		for index, value := range sites {
			fmt.Println("Testando site", index+1)
			testaSite(value)
		}
		time.Sleep(intervaloDeMonitoramento)
	}
	fmt.Println("---------------------\n ")

}

func testaSite(site string) {
	resp, err := http.Get(site)

	if err != nil {
		registraLog(site, false)
		fmt.Println("Error no site", site, "\nError:", err.Error())
		return
	}

	if resp.StatusCode == 200 {
		registraLog(site, true)
		fmt.Println("O site", site, "Foi carregado com sucesso")
	} else {
		registraLog(site, false)
		fmt.Println("Site offline", site)
	}
}

func lerSitesTXT() []string {
	var sites []string
	arquivo, err := os.Open("sites.txt")
	if err != nil {
		fmt.Println(err)
	}

	leitor := bufio.NewReader(arquivo)
	for {
		linha, err := leitor.ReadString('\n')
		linha = strings.TrimSpace(linha)
		sites = append(sites, linha)
		if err == io.EOF { // EOF erro de final do arquivo
			break // quebra o loop
		}
	}
	arquivo.Close()

	return sites
}

func registraLog(site string, status bool) {
	arquivo, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(err)
	}

	horario := time.Now().Local().Format("02/01/2006 15:04:05")
	arquivo.WriteString(horario + " - " + site + " - ONLINE: " + strconv.FormatBool(status) + "\n")
	arquivo.Close()
}

func lerLogsTXT() {
	arquivo, err := ioutil.ReadFile("log.txt")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("-----------------")
	fmt.Println("IMPRIMINDO LOGS")
	fmt.Println("-----------------")
	fmt.Println(string(arquivo))
	fmt.Println("-----------------")
}

func limparConsole() {
	limpar := make(map[string]func())

	limpar["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	limpar["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	sitemaOperacional := runtime.GOOS
	funcLimpar, ok := limpar[sitemaOperacional]
	if ok {
		funcLimpar()
	} else {
		panic("Plataforma não suportada!")
	}
}
