package main

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var hostAddress, username, password string
var vlanID int
var all bool = false
var debug bool = false // Global debug flag

func main() {
	// Create the root command
	var rootCmd = &cobra.Command{
		Use:   "redfish-cli",
		Short: "A CLI tool for managing Redfish-enabled devices",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set log level to debug if --debug or -v is passed
			if debug {
				logrus.SetLevel(logrus.DebugLevel)
				logrus.Debug("Debug logging enabled")
			} else {
				logrus.SetLevel(logrus.InfoLevel)
			}
		},
	}

	// Add the global --debug flag
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "v", false, "Enable debug logging")

	// Add the "interfaces" command
	interfacesCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	interfacesCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	interfacesCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	interfacesCmd.Flags().BoolVarP(&all, "all", "A", false, "Show all interfaces (default: only active)")
	interfacesCmd.MarkFlagRequired("hostaddress")
	interfacesCmd.MarkFlagRequired("username")
	interfacesCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(interfacesCmd)

	// Add the "setvlanid" command
	setVlanIDCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	setVlanIDCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	setVlanIDCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	setVlanIDCmd.MarkFlagRequired("hostaddress")
	setVlanIDCmd.MarkFlagRequired("username")
	setVlanIDCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(setVlanIDCmd)

	// Add the "setOneTimeHTTPBoot" command
	setOneTimeHTTPBootCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	setOneTimeHTTPBootCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	setOneTimeHTTPBootCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	setOneTimeHTTPBootCmd.MarkFlagRequired("hostaddress")
	setOneTimeHTTPBootCmd.MarkFlagRequired("username")
	setOneTimeHTTPBootCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(setOneTimeHTTPBootCmd)

	// Add the "RebootServer" command
	rebootServerCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	rebootServerCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	rebootServerCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	rebootServerCmd.MarkFlagRequired("hostaddress")
	rebootServerCmd.MarkFlagRequired("username")
	rebootServerCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(rebootServerCmd)

	// Add the "StopServer" command
	stopServerCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	stopServerCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	stopServerCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	stopServerCmd.MarkFlagRequired("hostaddress")
	stopServerCmd.MarkFlagRequired("username")
	stopServerCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(stopServerCmd)

	// Add the "StartServer" command
	startServerCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	startServerCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	startServerCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	startServerCmd.MarkFlagRequired("hostaddress")
	startServerCmd.MarkFlagRequired("username")
	startServerCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(startServerCmd)

	// Add the "setVlanReboot" command
	setVlanRebootCmd.Flags().StringVarP(&hostAddress, "hostaddress", "H", "", "Host address (e.g., https://19.89.89.8:443)")
	setVlanRebootCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	setVlanRebootCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	setVlanRebootCmd.MarkFlagRequired("hostaddress")
	setVlanRebootCmd.MarkFlagRequired("username")
	setVlanRebootCmd.MarkFlagRequired("password")
	rootCmd.AddCommand(setVlanRebootCmd)

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing command: %v", err)
	}
}
