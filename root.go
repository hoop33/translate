package translate

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

var (
	sl, tl   string
	insecure bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&sl, "sl", "en", "source language (os env property $TRANSLATE_SL has priority)")
	rootCmd.PersistentFlags().StringVar(&tl, "tl", "it", "target language (os env property $TRANSLATE_TL has priority)")
	rootCmd.PersistentFlags().BoolVar(&insecure, "insecure", false, "ignore certificate errors")
}

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "a simple cli for translation from google.",
	RunE:  transl}

func transl(_ *cobra.Command, args []string) error {

	if len(args) == 0 {
		return fmt.Errorf("please provide at least 1 arg to translate\n")
	}
	var query string
	for _, q := range args {
		escaped := q + " "
		query += escaped
	}
	slProp := os.Getenv("TRANSLATE_SL")
	if slProp != "" {
		sl = slProp
	}
	tlProp := os.Getenv("TRANSLATE_TL")
	if tlProp != "" {
		tl = tlProp
	}
	if insecure {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	res, err := Translate(sl, tl, query)
	if err != nil {
		return fmt.Errorf("could not translate [%v]: %v", query, err)
	}
	fmt.Fprintln(os.Stdout, res)
	return nil
}

//Execute translate command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "could not execute command: %v", err)
		os.Exit(1)
	}
}
