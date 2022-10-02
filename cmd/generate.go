package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/petewall/eink-radiator-image-source-image/v2/internal"
)

var Config *internal.Config

func parseConfig(cmd *cobra.Command, args []string) error {
	var err error
	Config, err = internal.ParseConfig(viper.GetString("config"))
	return err
}

var GenerateCmd = &cobra.Command{
	Use:     "generate",
	Short:   "Generates a " + ImageTypeName + " image",
	PreRunE: parseConfig,
	RunE: func(cmd *cobra.Command, args []string) error {
		imageContext := Config.GenerateImage(viper.GetInt("width"), viper.GetInt("height"))

		if viper.GetBool("to-stdout") {
			return imageContext.EncodePNG(cmd.OutOrStdout())
		}

		return imageContext.SavePNG(viper.GetString("output"))
	},
}

const (
	DefaultImageHeight    = 480
	DefaultImageWidth     = 640
	DefaultOutputFilename = ImageTypeName + ".png"
)

func init() {
	rootCmd.AddCommand(GenerateCmd)
	GenerateCmd.Flags().StringP("config", "c", "", "the path to the image config file")
	GenerateCmd.Flags().Int("height", DefaultImageHeight, "the height of the image")
	GenerateCmd.Flags().Int("width", DefaultImageWidth, "the width of the image")

	GenerateCmd.Flags().StringP("output", "o", DefaultOutputFilename, "path to write the file")
	GenerateCmd.Flags().Bool("to-stdout", false, "print the image to stdout")
	_ = viper.BindPFlags(GenerateCmd.Flags())
}
