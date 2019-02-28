package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/RangelReale/osincli"
	"github.com/fatih/color"
	cobra "github.com/spf13/cobra"

	"minibox.ai/minibox/pkg/utils"
)

var WaitLoginDuration = 5 * time.Minute

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Sign in minibox.ai website, use apis",
	Long: `Used minibox.ai dev account to sign in, can use apis to managament project, download/upload data, 
		running a training jobs..
        Complete documentation is available at http://docs.minibox.ai`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("login start")

		// client := oclient.NewClient("12345", "asdfasdf", "http://localhost:14001/auth")
		// u := client.AuthCodeURL("")

		cliconfig := &osincli.ClientConfig{
			ClientId:     "12345",
			ClientSecret: "asdfasdf",
			AuthorizeUrl: "http://localhost:14000/oauth/authorize",
			TokenUrl:     "http://localhost:14000/oauth/token",
			RedirectUrl:  "http://localhost:14001/auth",
			Scope:        "cli",
		}
		client, err := osincli.NewClient(cliconfig)
		if err != nil {
			panic(err)
		}
		// create a new request to generate the url
		areq := client.NewAuthorizeRequest(osincli.CODE)

		u := areq.GetAuthorizeUrl()
		utils.Browser(u.String())
		fmt.Printf("Copy this URL: [%s] in your Browser, if you no saw any new Browser window are opend for login.\n", color.CyanString(u.String()))

		server := NewAuthService(areq, client)
		server.duration = time.Minute
		server.Run(":14001")
		// Do Stuff Here
	},
}

type AuthService struct {
	*http.ServeMux
	duration time.Duration
	ch       chan bool
	areq     *osincli.AuthorizeRequest
	client   *osincli.Client
}

func NewAuthService(areq *osincli.AuthorizeRequest, client *osincli.Client) *AuthService {
	svr := &AuthService{
		ServeMux: http.NewServeMux(),
		duration: WaitLoginDuration,
		ch:       make(chan bool),
		areq:     areq,
		client:   client,
	}
	svr.init()
	return svr
}

func (svr *AuthService) init() {
	svr.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received authorize info\n")
		areqdata, err := svr.areq.HandleRequest(r)
		if err != nil {
			// w.Write([]byte(fmt.Sprintf("ERROR: %s\n", err)))
			fmt.Printf("ERROR: %s\n", err)
			return
		}

		treq := svr.client.NewAccessRequest(osincli.AUTHORIZATION_CODE, areqdata)

		// show access request url (for debugging only)
		// u2 := treq.GetTokenUrl()
		// // w.Write([]byte(fmt.Sprintf("Access token URL: %s\n", u2.String())))
		// fmt.Printf("Access token URL: %s\n", u2.String())
		// exchange the authorize token for the access token
		ad, err := treq.GetToken()
		if err != nil {
			// w.Write([]byte(fmt.Sprintf("ERROR: %s\n", err)))
			fmt.Printf("ERROR: %s\n", err)
			return
		}
		// w.Write([]byte(fmt.Sprintf("Access token: %+v\n", ad)))
		// closeScript(w)
		jwtToken, err := svr.downloadJWT("http://localhost:14000/oauth/jwt", ad.AccessToken)
		if err != nil {
			log.Printf("download jwt error %s", err)
		}

		saveAccesToken(jwtToken)

		http.Redirect(w, r, "http://localhost:14000/oauth/app", http.StatusFound)
		go func() {
			time.Sleep(1 * time.Second)
			svr.ch <- true
		}()
	})
}

func (svr *AuthService) downloadJWT(url, code string) ([]byte, error) {
	var cli = &http.Client{}
	// dir := filepath.Dir(config.Endpoint.AuthURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+code)

	var jwtResp = struct {
		JWT string `json:"jwt"`
	}{}

	if resp, err := cli.Do(req); err != nil {
		return nil, err
	} else {
		dec := json.NewDecoder(resp.Body)
		if err := dec.Decode(&jwtResp); err != nil {
			return nil, err
		} else {
			return []byte(jwtResp.JWT), nil
		}
	}
}

func (svr *AuthService) Run(addr string) (err error) {
	var (
		server = &http.Server{Addr: addr, Handler: svr}
	)

	go func() {
		err = server.ListenAndServe()
	}()

	select {
	case <-time.After(svr.duration):
		server.Close()
		log.Printf("sign in failed, timeout.")
	case <-svr.ch:
		server.Close()
		fmt.Println("Login Success！")
	}

	return
}
