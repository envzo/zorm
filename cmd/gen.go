package cmd

import (
	"os"
	"path/filepath"

	"github.com/envzo/zorm/gen"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "zorm",
	Short: "Generate Go code for ORM",
	Run: func(cmd *cobra.Command, args []string) {
		if yaml == "" {
			println("yaml file must be provided")
			return
		}
		if folder == "" {
			println("destination folder must be provided")
			return
		}

		var err error

		if folder, err = filepath.Abs(folder); err != nil {
			panic(err)
		}

		if pkg == "" {
			pkg = filepath.Base(folder)
			// in this scenario, we don't create new folder
		} else {
			if err = os.MkdirAll(filepath.Join(folder, pkg), 0755); err != nil {
				panic(err)
			}
			folder = filepath.Join(folder, pkg)
		}

		if err = gen.Gen(yaml, folder, pkg); err != nil {
			panic(err)
		}
	},
}

var (
	yaml   string
	folder string
	pkg    string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&yaml, "file", "f", "", "ORM definition YAML file")
	rootCmd.PersistentFlags().StringVarP(&folder, "folder", "F", "", "folder where generated sql scripts & code will write to")
	rootCmd.PersistentFlags().StringVarP(&pkg, "package", "p", "", "package name of generated code")
}

func Exec() error {
	return rootCmd.Execute()
}
