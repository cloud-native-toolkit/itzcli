package cmd

// // importCmd represents the import command
// var listConfigCmd = &cobra.Command{
// 	Use:   "list",
// 	Short: "Lists cluster configuration.",
// 	Long:  `Lists cluster configuration from ocpnow.`,
// 	RunE: func(cmd *cobra.Command, args []string) error {
// 		ocpCfg := viper.GetStringSlice("ocpnow.configFiles")
// 		if len(ocpCfg) == 0 {
// 			logger.Warnf("No ocpnow configuration files found")
// 			return nil
// 		}
// 		logger.Debugf("Listing configuration from <%s>...", ocpCfg[0])
// 		project, err := pkg.LoadProject(ocpCfg[0])
// 		if err != nil {
// 			return err
// 		}
// 		err = project.Write(configureCmd.OutOrStdout())
// 		return err
// 	},
// }

// func init() {
// 	clusterCmd.AddCommand(listConfigCmd)
// }
