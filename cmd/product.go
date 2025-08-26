/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/opt-nc/geol/cmd/cache/utilities"
	"github.com/spf13/cobra"
)

// productCmd represents the product command
var productCmd = &cobra.Command{
	Use:   "product",
	Short: "Display the latest version of one or more products and the end of life date.",
	Long:  "Display the latest version of one or more products and the end of life date.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Veuillez spécifier au moins un produit.")
			return
		}

		// Charger le cache local
		productsPath, err := utilities.GetProductsPath()
		if err != nil {
			fmt.Println("Erreur lors de la récupération du chemin du cache:", err)
			return
		}
		cacheFile, err := os.Open(productsPath)
		if err != nil {
			fmt.Println("Erreur lors de l'ouverture du cache local:", err)
			return
		}
		defer func() {
			if err := cacheFile.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Erreur lors de la fermeture du cache local: %v\n", err)
			}
		}()

		var cache struct {
			Products map[string][]string `json:"products"`
		}
		if err := json.NewDecoder(cacheFile).Decode(&cache); err != nil {
			fmt.Println("Erreur lors du décodage du cache:", err)
			return
		}

		// Structure pour stocker les résultats
		type ProductResult struct {
			Name        string
			EolLabel    string
			ReleaseName string
			ReleaseDate string
			EolFrom     string
			Link        string
		}
		var results []ProductResult

		for _, prod := range args {
			found := false
			for name, aliases := range cache.Products {
				if strings.EqualFold(prod, name) {
					found = true
					prod = name
					break
				}
				for _, alias := range aliases {
					if strings.EqualFold(prod, alias) {
						found = true
						prod = name
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue // produit non trouvé dans le cache
			}

			// Requête API pour ce produit
			url := utilities.ApiUrl + "products/" + prod
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("Erreur lors de la requête pour %s: %v\n", prod, err)
				continue
			}
			body, err := io.ReadAll(resp.Body)
			if cerr := resp.Body.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "Erreur lors de la fermeture du body HTTP pour %s: %v\n", prod, cerr)
			}
			if err != nil {
				fmt.Printf("Erreur lecture réponse pour %s: %v\n", prod, err)
				continue
			}
			if resp.StatusCode != 200 {
				fmt.Printf("Produit %s introuvable sur l'API.\n", prod)
				continue
			}

			// Décodage JSON
			var apiResp struct {
				Result struct {
					Name   string `json:"name"`
					Labels struct {
						Eol string `json:"eol"`
					} `json:"labels"`
					Releases []struct {
						Name        string `json:"name"`
						ReleaseDate string `json:"releaseDate"`
						EolFrom     string `json:"eolFrom"`
					} `json:"releases"`
					Links struct {
						HTML string `json:"html"`
					} `json:"links"`
				} `json:"result"`
			}
			if err := json.Unmarshal(body, &apiResp); err != nil {
				fmt.Printf("Erreur décodage JSON pour %s: %v\n", prod, err)
				continue
			}
			var relName, relDate, relEol string
			if len(apiResp.Result.Releases) > 0 {
				relName = apiResp.Result.Releases[0].Name
				relDate = apiResp.Result.Releases[0].ReleaseDate
				relEol = apiResp.Result.Releases[0].EolFrom
			}
			results = append(results, ProductResult{
				Name: apiResp.Result.Name,
				//EolLabel:    apiResp.Result.Labels.Eol,
				ReleaseName: relName,
				ReleaseDate: relDate,
				EolFrom:     relEol,
				Link:        apiResp.Result.Links.HTML,
			})
		}

		// Affichage tableau markdown avec glamour
		if len(results) == 0 {
			fmt.Println("Aucun produit trouvé dans le cache et l'API.")
			return
		}
		var buf bytes.Buffer
		// En-tête
		buf.WriteString("| **Name** | **Version** | **Release Date** | **EOL From** | **Link** |\n")
		buf.WriteString("|------|--------------|--------------|----------|------|\n")
		for _, r := range results {
			name := strings.ReplaceAll(r.Name, "|", "\\|")
			//eolLabel := strings.ReplaceAll(r.EolLabel, "|", "\\|")
			relName := strings.ReplaceAll(r.ReleaseName, "|", "\\|")
			relDate := strings.ReplaceAll(r.ReleaseDate, "|", "\\|")
			relEol := strings.ReplaceAll(r.EolFrom, "|", "\\|")
			link := strings.ReplaceAll(r.Link, "|", "\\|")
			buf.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", name, relName, relDate, relEol, link))
		}
		// r, _ := glamour.NewTermRenderer(
		// 	glamour.WithAutoStyle(),
		// 	//glamour.WithWordWrap(120),
		// )

		//out, err := glamour.Render(buf.String(), "dark")
		//out, err := r.Render(buf.String())
		out, err := glamour.RenderWithEnvironmentConfig(buf.String())
		if err != nil {
			fmt.Print(buf.String()) // fallback brut
		} else {
			fmt.Print(out)
		}
	},
}

func init() {
	rootCmd.AddCommand(productCmd)
}
