package cmd

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"animedrive-dl/config"
	"animedrive-dl/utils"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [link] [folder]",
	Short: "🎬 Download an anime link",
	Long:  `🎬 Download an anime link and perform the necessary actions.`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		baseFolder := args[1]
		var headers map[string]string
		var err error

		fmt.Println("🔗 Parsing URL...")

		if strings.Contains(url, "animedrive.hu/anime/?id=") {
			headers, err = config.LoadHeaders("headers.json", "baseHeaders")
			if err != nil {
				fmt.Printf("❌ Error loading baseHeaders: %v\n", err)
				return
			}

			fmt.Println("📡 Fetching number of episodes...")
			episodes, err := utils.FetchNumberOfEpisodes(url, headers)
			if err != nil {
				fmt.Printf("❌ Error fetching number of episodes: %v\n", err)
				return
			}
			fmt.Printf("📺 Available episodes: %d\n", episodes)

			prompt := promptui.Prompt{
				Label: fmt.Sprintf("How many episodes do you want to download? (max %d)", episodes),
				Validate: func(input string) error {
					n, err := strconv.Atoi(input)
					if err != nil || n < 1 || n > episodes {
						return fmt.Errorf("Invalid number of episodes")
					}
					return nil
				},
			}

			result, err := prompt.Run()
			if err != nil {
				fmt.Printf("❌ Prompt failed: %v\n", err)
				return
			}

			numEpisodes, _ := strconv.Atoi(result)
			fmt.Printf("Downloading %d episodes...\n", numEpisodes)

			// Extract the anime ID
			re := regexp.MustCompile(`animedrive.hu/anime/\?id=(\d+)`)
			match := re.FindStringSubmatch(url)
			if len(match) != 2 {
				fmt.Println("❌ Error: could not extract anime ID from URL")
				return
			}
			animeID := match[1]

			for i := 1; i <= numEpisodes; i++ {
				playerURL := fmt.Sprintf("https://player.animedrive.hu/player_v1.5.php?id=%s&ep=%d", animeID, i)
				fmt.Printf("🌐 Downloading episode %d: %s\n", i, playerURL)
				animeName, err := utils.FetchAnimeName(fmt.Sprintf("https://animedrive.hu/watch/?id=%s&ep=%d", animeID, i), headers)
				if err != nil {
					fmt.Printf("❌ Error fetching anime name: %v\n", err)
					return
				}
				playerHeaders, err := config.LoadHeaders("headers.json", "playerHeaders")
				if err != nil {
					fmt.Printf("❌ Error loading playerHeaders: %v\n", err)
					return
				}
				utils.FetchAndDownload(playerURL, playerHeaders, baseFolder, fmt.Sprintf("%d", i), animeName)
			}

		} else if strings.Contains(url, "animedrive.hu/watch/?id=") {
			headers, err = config.LoadHeaders("headers.json", "playerHeaders")
			if err != nil {
				fmt.Printf("❌ Error loading playerHeaders: %v\n", err)
				return
			}

			if strings.Contains(url, "&ep=") {
				// Extract id and ep from URL
				re := regexp.MustCompile(`animedrive.hu/watch/\?id=(\d+)&ep=(\d+)`)
				match := re.FindStringSubmatch(url)
				if len(match) == 3 {
					id := match[1]
					ep := match[2]
					fmt.Println("📡 Fetching anime name...")
					animeName, err := utils.FetchAnimeName(url, headers)
					if err != nil {
						fmt.Printf("❌ Error fetching anime name: %v\n", err)
						return
					}
					fmt.Printf("📺 Anime: %s\n", animeName)
					playerURL := fmt.Sprintf("https://player.animedrive.hu/player_v1.5.php?id=%s&ep=%s", id, ep)
					utils.FetchAndDownload(playerURL, headers, baseFolder, ep, animeName)
					return
				} else {
					fmt.Println("❌ Error: wrong URL format for watch URL")
					return
				}
			}
		} else {
			fmt.Println("❌ Error: wrong URL")
			return
		}

		fmt.Printf("🌐 Downloading from: %s\n", url)
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
