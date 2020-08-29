package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/dsnezhkov/cflsh/common"
	"github.com/manifoldco/promptui"

)

func init() {
	rootCmd.AddCommand(bucketsCmd)
}

var bucketsCmd = &cobra.Command{
	Use:   "buckets",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		bucketsDriver(cmd, args)
	},
}

func bucketsDriver(cmd *cobra.Command, args []string) {

	nameSpace := common.NameSpace{
		Name: "test",
		UUID: "9b5e2ff9353f4f618a29193d2facc01b",
	}

	prompt := promptui.Prompt{
		Label: "commands",
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	switch result {
	case "list":
		resp, err := common.ListWorkersKVs(nameSpace.UUID)
		if err !=nil {
			fmt.Printf("%s\n", err)
		}
		if len(resp.Result) != 0 {
			for _,x := range resp.Result {
				fmt.Printf("%s\n", x.Name)
			}
		}
	case "dump":
		resp, err := common.ListWorkersKVs(nameSpace.UUID)
		if err !=nil {
			fmt.Printf("%s\n", err)
		}
		if len(resp.Result) != 0 {
			cmdRRequest := common.CmdKVRRequest{
				NameSpace:  nameSpace,
			}
			for _,x := range resp.Result {
				cmdRRequest.CmdKName = x.Name
				payload, err := common.ReadCmdKV(cmdRRequest)
				if err!= nil {
					fmt.Printf("Cannot get %s: %s\n", x.Name, err)
				}

				fmt.Printf("%s: %s\n", x.Name, payload)
			}
		}

	default:
		fmt.Printf("Invalid %s.\n", result )
	}
}
