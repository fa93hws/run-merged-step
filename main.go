package main

import (
	"github.com/fa93hws/run-merged-step/cmd"
)

func main() {
	cmd.Execute()
	// var (
	// 	verbose bool
	// 	rootCmd = &cobra.Command{
	// 		Use:   "hugo",
	// 		Short: "Hugo is a very fast static site generator",
	// 		Long: `A Fast and Flexible Static Site Generator built with
	// 									love by spf13 and friends in Go.
	// 									Complete documentation is available at https://gohugo.io/documentation/`,
	// 		Run: func(cmd *cobra.Command, args []string) {
	// 			fmt.Println("verbose:", verbose)
	// 			fmt.Println("rest of args:", args)
	// 			fmt.Println("hello world")
	// 		},
	// 	}
	// )

	// rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "whether to log verbose")
	// // rootCmd.MarkPersistentFlagRequired("verbose")

	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
}
