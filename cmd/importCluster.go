package cmd

// var ocpnowCfg string

// // importCmd represents the import command
// var importCmd = &cobra.Command{
// 	Use:    "import",
// 	Short:  "Imports cluster configuration from ocpnow.",
// 	Long:   `Imports cluster configuration from ocpnow.`,
// 	PreRun: SetLoggingLevel,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		logger.Debugf("Importing configuration from <%s>...", ocpnowCfg)
// 		return importOcpnowConfig(ocpnowCfg)
// 	},
// }

// func importOcpnowConfig(cfg string) error {
// 	// First thing we're going to do is to copy the file into the .itz home directory..
// 	// HACK: Here we should peek inside the file and use the project name as
// 	// the config file name.
// 	configDir := filepath.Dir(viper.ConfigFileUsed())
// 	importedCfg := filepath.Join(configDir, "project.yaml")
// 	err := pkg.AppendToFile(cfg, importedCfg)
// 	if err != nil {
// 		return err
// 	}
// 	logger.Tracef("Storing project config <%s> in configuration directory <%s>...", cfg, configDir)
// 	currentFiles := viper.GetStringSlice("ocpnow.configFiles")
// 	currentFiles = append(currentFiles, importedCfg)
// 	viper.Set("ocpnow.configFiles", currentFiles)
// 	return viper.WriteConfig()
// }

// func init() {
// 	clusterCmd.AddCommand(importCmd)
// 	importCmd.Flags().StringVarP(&ocpnowCfg, "from-ocpnow-project", "p", "", "Specifies the project.yaml file created by ocpnow.")
// }
