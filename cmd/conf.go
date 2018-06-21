package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// confCmd represents the conf command
var confCmd = &cobra.Command{
	Use:   "conf",
	Short: "Verifies the configuration",
	Long: `Verifies that git, hugo and rsync are installed in the machine.
Checks that there's a valid hugo git project in $SRCPATH.`,
	Run: runConf,
}

func init() {
	RootCmd.AddCommand(confCmd)
}

func runConf(cmd *cobra.Command, args []string) {
	for _, cmd := range [][2]string{{"git", "version"}, {"hugo", "version"}, {"rsync", "--version"}} {
		out, err := exec.Command(cmd[0], cmd[1]).CombinedOutput()
		if err != nil {
			logger.Fatalf("%s: %s", cmd[0], err)
		}
		logger.Println(strings.SplitN(string(out), "\n", 2)[0])
	}
	if err := os.Chdir(conf.SrcPath); err != nil {
		logger.Fatalln("cd:", conf.SrcPath, err.Error())
	}
	if err := exec.Command("git", "status").Run(); err != nil {
		logger.Fatalln("git:", err.Error())
	}
}
