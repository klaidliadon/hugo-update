package cmd

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:    "run",
	Short:  "Runs the web server",
	Long:   `Runs the web server and keeps the $DSTPATH updated`,
	PreRun: runConf,
	Run:    runRun,
}

func init() {
	RootCmd.AddCommand(runCmd)
}

func runRun(cmd *cobra.Command, args []string) {
	http.HandleFunc(conf.Handler, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.NotFound(w, r)
			return
		}
		defer r.Body.Close()

		mac := hmac.New(sha1.New, []byte(conf.Secret))
		io.Copy(mac, r.Body)
		if s := fmt.Sprintf("sha1=%x", mac.Sum(nil)); s != r.Header.Get("X-Hub-Signature") {
			logger.Printf("Invalid signature %q, expected %q:", r.Header.Get("X-Hub-Signature"), s)
			http.Error(w, "Signature mismatch", http.StatusConflict)
			return
		}

		logger.Print("Start Update...")
		v, err := refreshContent()
		if err != nil {
			logger.Println("Update ERROR:", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		logger.Println("Update Complete!", v)

	})
	logger.Printf("Listening on :%v from %q to %q", conf.Port, conf.SrcPath, conf.DstPath)
	http.ListenAndServe(fmt.Sprintf(":%v", conf.Port), nil)
}

func refreshContent() (string, error) {
	if out, err := exec.Command("git", "pull").CombinedOutput(); err != nil {
		return "", fmt.Errorf("git: %s\n%s", err, string(out))
	}
	out, err := exec.Command("git", "rev-parse", "HEAD").CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git: %s", err)
	}
	if out, err := exec.Command("hugo").CombinedOutput(); err != nil {
		return "", fmt.Errorf("hugo: %s\n%s", err, string(out))
	}
	if out, err := exec.Command("rsync", "-ravz", "./public/", conf.DstPath).CombinedOutput(); err != nil {
		return "", fmt.Errorf("rsync: %s\n%s", err, string(out))
	}

	return strings.TrimSpace(string(out)), nil
}
