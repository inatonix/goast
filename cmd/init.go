package cmd

// import (
// 	"fmt"
// 	"os"

// 	"github.com/spf13/cobra"
// )

// // func Initialize() {
// // 	err :=
// // }

// var apiChatEngine string
// var apiModel string
// var apiDBSchema string
// var apiReadFiles string
// var apiVerbose bool

// var apiCmd = &cobra.Command{
// 	Use:   "api [Project directory] [OpenAPI document path]",
// 	Short: "Generate the new API code",
// 	Long:  `This command is used for chat AI to generate the new API code from the project codes and OpenAPI document.`,
// 	Args:  cobra.ExactArgs(2),
// 	Run: func(cmd *cobra.Command, args []string) {
// 		projectDir := args[0]
// 		openAPIDoc := args[1]

// 		input := app.ApiInput{
// 			ProjectDir:    projectDir,
// 			OpeanAPIDoc:   openAPIDoc,
// 			DBSchema:      apiDBSchema,
// 			ReadFilesPath: apiReadFiles,
// 		}

// 		config := app.ApiConfig{
// 			Input:      input,
// 			ChatEngine: apiChatEngine,
// 			Model:      apiModel,
// 			Verbose:    apiVerbose,
// 		}

// 		err := app.Api(config)
// 		if err != nil {
// 			fmt.Println(err)
// 			_ = cmd.Usage()
// 			os.Exit(1)
// 		}
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(apiCmd)

// 	apiCmd.PersistentFlags().StringVar(&apiChatEngine, "chat_engine", "chatgpt", "Chat engine (chatgpt or claude)")
// 	apiCmd.PersistentFlags().StringVar(&apiModel, "model", "gpt-4-1106-preview", "Model name (gpt-3.5-turbo, gpt-4, claude2 etc)")
// 	apiCmd.PersistentFlags().StringVar(&apiDBSchema, "db_schema", "", "Database schema file path")
// 	apiCmd.PersistentFlags().StringVar(&apiReadFiles, "read_file", ".read", "List of the read file path")
// 	apiCmd.PersistentFlags().BoolVarP(&apiVerbose, "verbose", "v", false, "flag for verbose output")
// }
