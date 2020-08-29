package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dsnezhkov/cflsh/common"


	"github.com/manifoldco/promptui"

	"github.com/lithammer/shortuuid/v3"
)

func init() {
	rootCmd.AddCommand(redCmd)
}

var redCmd = &cobra.Command{
	Use:   "red",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		redDriver(cmd, args)
	},
}

func redDriver(cmd *cobra.Command, args []string) {

	nameSpace := common.NameSpace{
		Name: "test",
		UUID: "9b5e2ff9353f4f618a29193d2facc01b",
	}
	CmdWRequest := common.CmdKVWRequest{
		NameSpace:  nameSpace,
	}

	// Init command bucket with known state
	CmdWRequest.CmdKName = common.CmdBucketKeyPrefix
	CmdWRequest.CmdPayload = []byte("")
	resp, err := common.WriteCmdKV(CmdWRequest)
	if err != nil {
		fmt.Printf("Error in WriteCmdKV worker: %s\n", err)
		fmt.Printf("%+v", resp)
	}

	// Loop
	for {

		prompt := promptui.Prompt{
			Label: "command",
		}

		pcmd, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		oneTimeKey := shortuuid.New()
		requestBucketName := common.CmdBucketKeyPrefix + oneTimeKey
		fmt.Printf("Setting up bucket: %s\n", requestBucketName)

		// Write command request payload

		CmdWRequest.CmdKName = requestBucketName
		CmdWRequest.CmdPayload =[]byte(pcmd)

		resp, err := common.WriteCmdKV(CmdWRequest)
		if err != nil {
			fmt.Printf("Error in WriteCmdKV worker: %s\n", err)
			fmt.Printf("%+v", resp)
		}

		// Wire up bucket in KV index
		CmdWRequest.CmdKName = common.CmdBucketKeyPrefix
		CmdWRequest.CmdPayload = []byte(requestBucketName)
		resp, err = common.WriteCmdKV(CmdWRequest)
		if err != nil {
			fmt.Printf("Error in WriteCmdKV worker: %s\n", err)
			fmt.Printf("%+v", resp)
		}

		// Pre-setup output bucket
		outputPayloadBucket := strings.ReplaceAll(
			string(requestBucketName), "command", "output")

		fmt.Printf("Setting up output bucket: %s\n", outputPayloadBucket)

		// Write command request payload
		OutWRequest := common.OutKVWRequest{
			NameSpace:  nameSpace,
		}
		OutWRequest.OutKName = outputPayloadBucket
		OutWRequest.OutPayload =[]byte("")

		resp, err = common.WriteOutKV(OutWRequest)
		if err != nil {
			fmt.Printf("Error in WriteOutKV worker: %s\n", err)
			fmt.Printf("%+v", resp)
		}

		// Wire up bucket in KV index
		OutWRequest.OutKName = common.OutBucketKeyPrefix
		OutWRequest.OutPayload = []byte(outputPayloadBucket)
		resp, err = common.WriteOutKV(OutWRequest)
		if err != nil {
			fmt.Printf("Error in WriteOutKV worker: %s\n", err)
			fmt.Printf("%+v", resp)
		}



		OutRRequest := common.OutKVRRequest{
			NameSpace:  nameSpace,
			OutKName: outputPayloadBucket,
		}

		for {

			time.Sleep(common.WorkerInterval * time.Second)

			// Get bucket name from worker KV index
			outPayload, err := common.ReadOutKV(OutRRequest)
			if err != nil {
				fmt.Printf("Error getting output bucket name: %s\n", err)
				fmt.Printf("%+v", resp)
			}

			// Empty / waiting?
			if len(outPayload) == 0 {
				continue
			}else {
				fmt.Printf("Got output: `%s`\n", outPayload)

				// Signal complete command by removing output payload
				OutWRequest.OutPayload = []byte("")
				resp, err = common.WriteOutKV(OutWRequest)
				if err != nil {
					fmt.Printf("Error writing output bucket: %s\n", err)
					fmt.Printf("%+v\n", resp)
				}

				break // get out of the wait loop into prompt
			}
		}

	}

}
