// Copyright © 2018 mritd <mritd1234@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"net/http"

	"bytes"

	"net"

	"strconv"
	"strings"

	"io"

	"github.com/gobuffalo/packr"
	"github.com/mritd/crxdl/downloader"
	"github.com/mritd/crxdl/utils"
	"github.com/spf13/cobra"
)

var ostype string
var arch string
var prod string
var prodchannel string
var prodversion string
var listen string
var port int

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "crxdl",
	Short: "Chrome crx download tool",
	Long: `
A simple chrome crx download tool`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		Start()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {

	RootCmd.PersistentFlags().StringVarP(&ostype, "ostype", "o", "mac", "OS type")
	RootCmd.PersistentFlags().StringVarP(&arch, "arch", "a", "x86-64", "OS arch")
	RootCmd.PersistentFlags().StringVarP(&prod, "prod", "d", "chromecrx", "Chrome crx type(chromecrx/chromiumcrx)")
	RootCmd.PersistentFlags().StringVarP(&prodchannel, "prodchannel", "c", "unknown", `Channel is "unknown" on Chromium on ArchLinux, so using "unknown" will probably be fine for everyone.`)
	RootCmd.PersistentFlags().StringVarP(&prodversion, "prodversion", "v", "66.0.3359.139", "Chrome version")
	RootCmd.PersistentFlags().StringVarP(&listen, "listen", "l", "0.0.0.0", "http listen address")
	RootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "http listen port")
}

func dlServer(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	crxid := req.Form.Get("crxid")
	if strings.TrimSpace(crxid) == "" {
		fmt.Fprint(w, "Error: crxid is blank!")
		return
	}
	if crxid == "fuckgfw" {
		crxid = "padekgcemlokbadohgkifijomclgjgif"
	}

	dl := downloader.Downloader{
		OsType:      ostype,
		Arch:        arch,
		Prod:        prod,
		ProdChannel: prodchannel,
		ProdVersion: prodversion,
		ExtensionID: crxid,
	}
	buf := dl.Download()
	if len(buf) == 0 {
		fmt.Fprint(w, "Extension not found!")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("content-disposition", "attachment; filename=\""+dl.ExtensionID+".crx\"")
		io.Copy(w, bytes.NewReader(buf))
	}

}

func Start() {

	addr, err := net.ResolveTCPAddr("tcp", listen+":"+strconv.Itoa(port))
	utils.CheckAndExit(err)

	fmt.Println("starting server at", addr.String())
	http.HandleFunc("/crxdl", dlServer)
	box := packr.NewBox("../resources")
	http.Handle("/", http.FileServer(box))

	fmt.Println("server starting...")
	utils.CheckAndExit(http.ListenAndServe(addr.String(), nil))
}
