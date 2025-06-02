package whisker

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/vialtek/whisker/model"
)

func (s *NodeState) recipeByName(recipeName string) *model.Recipe {
	for _, recipe := range s.Recipes {
		if recipe.Name == recipeName {
			return recipe
		}
	}

	return nil
}

func execRecipe(recipe *model.Recipe, result *Result) {
	var cmd *exec.Cmd
	scriptPath := fmt.Sprintf("%s/%s", recipe.Pwd, recipe.Entrypoint)

	cmd = exec.Command(recipe.Runner, scriptPath)

	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		m := scanner.Text()
		log.Println(m)
		result.Output = append(result.Output, m)
	}

	scannerErr := bufio.NewScanner(stderr)
	for scannerErr.Scan() {
		m := scannerErr.Text()
		log.Println("Error:", m)
		result.Output = append(result.Output, m)

		// Job failed, abort.
		result.Error = m
		result.Success = false
		result.EndedAt = time.Now()
	}

	cmd.Wait()
}

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
