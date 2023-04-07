package cmd

import (
	"fmt"
	"strings"

	"github.com/DanielPickens/hercules/internal/report"
	"github.com/DanielPickens/hercules/pkg"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints version/build info",
		Long:  "Prints version/build information",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
}

func printVersion() {
	printLogo(report.ColorAqua, report.ColorLighSlate)
	printTuple("Version", version)
	printTuple("Commit", commit)
	printTuple("Date", date)
	printTuple("Logs", pkg.LogFile)
}

func printTuple(section, value string) {
	const secFmt = "%-10s "
	fmt.Printf("%s", report.Colorize(fmt.Sprintf(secFmt, section+":"), report.ColorAqua))
	fmt.Println(report.Colorize(value, report.ColorWhite))
}

func printLogo(title, logo report.Color) {
	for i, s := range report.Logo {
		if i < len(report.Hercules) {
			fmt.Printf("%s", report.Colorize(report.Hercules[i], title))
			fmt.Printf("%s", strings.Repeat(" ", 22))
		} else {
			if i == 4 {
				fmt.Printf("%s", report.Colorize("  Hercules!", logo))
				fmt.Printf("%s", strings.Repeat(" ", 26))
			} else {
				fmt.Printf("%s", strings.Repeat(" ", 50))
			}
		}
		fmt.Println(report.Colorize(s, logo))
	}
}
