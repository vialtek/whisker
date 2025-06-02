package whisker

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"

	"github.com/vialtek/whisker/model"
)

func loadRecipes() []*model.Recipe {
	var recipes []*model.Recipe

	filepath.Walk(GetConfig().RecipeDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			recipePath := filepath.Join(path, "recipe.yaml")
			if _, err := os.Stat(recipePath); err == nil {
				recipes = append(recipes, parseRecipeYaml(path))
			}
		}

		return nil
	})

	log.Println("Recipes loaded:", len(recipes))
	if len(recipes) == 0 {
		log.Println("Warning: no recipes loaded.")
	}

	return recipes
}

func parseRecipeYaml(path string) *model.Recipe {
	yamlFile, err := os.ReadFile(path + "/recipe.yaml")
	if err != nil {
		log.Printf("parseRecipeYaml could not open file:", err)
	}

	newRecipe := &model.Recipe{}
	err = yaml.Unmarshal(yamlFile, newRecipe)
	if err != nil {
		log.Fatalf("parseRecipeYaml could not unmarshal:", err)
	} else {
		newRecipe.Pwd = path
	}

	return newRecipe
}
