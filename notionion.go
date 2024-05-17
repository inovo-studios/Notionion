package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/elazarl/goproxy"
	"github.com/jomei/notionapi"
)

func main() {
	port := "8080"
	flag.Parse()
	if len(flag.Args()) > 0 {
		port = flag.Arg(0)
	}

	// Integration token
	token := os.Getenv("NOTION_TOKEN")
	if token == "" {
		log.Fatal("❌ Please set NOTION_TOKEN envvar with your integration token before launching notionion")
	}

	// Page URL
	pageurl := os.Getenv("NOTION_PAGE_URL")
	if pageurl == "" {
		log.Fatal("❌ Please set NOTION_PAGE_URL envvar with your page id before launching notionion (CTRL+L on desktop app)")
	}

	// Extract Page ID
	pageid := pageurl[strings.LastIndex(pageurl, "-")+1:]
	if pageid == pageurl {
		log.Fatal("❌ PAGEID was not found in NOTION_PAGEURL. Ensure the url is in the form of https://notion.so/[pagename]-[pageid]")
	}

	// Initialize Notion client
	client := notionapi.NewClient(notionapi.Token(token))
	if client == nil {
		log.Fatal("❌ Failed to initialize Notion client")
	}

	// Check Page Content
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		log.Fatalf("Failed retrieving page children blocks: %v", err)
	}

	// Check Proxy Status
	if active, err := notionion.GetProxyStatus(children); err != nil {
		log.Println(err)
	} else if active {
		log.Println("📶 Proxy is active")
	} else {
		log.Println("📴 Proxy is inactive. Activate it by checking the \"OFF\" box")
	}

	// Request section checks
	if _, err := notionion.GetRequestBlock(children); err != nil {
		log.Fatalf("❌ Request block not found in the proxy page: %v", err)
	} else {
		log.Println("➡️ Request block found")
	}

	// Disable Request Buttons
	if err := notionion.DisableRequestButtons(client, pageid); err != nil {
		log.Println(err)
	}

	// Get Request Code Block
	codeReq, err := notionion.GetRequestCodeBlock(children)
	if err != nil {
		log.Fatalf("❌ Request code block not found in the proxy page: %v", err)
	}
	notionion.ClearRequestCode(client, codeReq.ID)

	// Response section checks
	if _, err := notionion.GetResponseBlock(children); err != nil {
		log.Fatalf("❌ Response block not found in the proxy page: %v", err)
	} else {
		log.Println("⬅️ Response block found")
	}

	// Get Response Code Block
	codeResp, err := notionion.GetResponseCodeBlock(children)
	if err != nil {
		log.Fatalf("❌ Response code block not found in the proxy page: %v", err)
	}
	notionion.ClearResponseCode(client, codeResp.ID)

	// Proxy Section
	proxy := goproxy.NewProxyHttpServer()
	//proxy.Verbose = true

	// Request HTTP Handler
	proxy.OnRequest().Do(notionion.ProxyRequestHTTPHandler(client, pageid, codeReq, codeResp))

	// Response Handler
	proxy.OnResponse().Do(notionion.ProxyResponseHTTPHandler(client, pageid, codeResp))

	log.Printf("🧅 Launch notionion proxy on port %s!\n", port)
	log.Fatal(http.ListenAndServe(":"+port, proxy))
}
