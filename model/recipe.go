package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type Recipe struct {
	ID              uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Slug            string `gorm:"column:slug" json:"slug,omitempty"`
	Name            string `gorm:"column:name" json:"name,omitempty"`
	Title           string `gorm:"column:title" json:"title,omitempty"`
	Subtitle        string `gorm:"column:subtitle" json:"subtitle,omitempty"`
	ImageURL        string `gorm:"column:image_url" json:"image_url,omitempty"`
	Tags            string `gorm:"column:tags;type:text" json:"tags,omitempty"`
	PrepTimeMinutes int    `gorm:"column:prep_time_minutes" json:"prepTimeMinutes,omitempty"`
	CookTimeMinutes int    `gorm:"column:cook_time_minutes" json:"cookTimeMinutes,omitempty"`
	Servings        int    `gorm:"column:servings" json:"servings,omitempty"`
	Ingredients     string `gorm:"column:ingredients;type:text" json:"ingredients,omitempty"`
	Steps           string `gorm:"column:steps;type:text" json:"steps,omitempty"`
	Calories        int    `gorm:"column:calories" json:"calories,omitempty"`
	ProteinG        int    `gorm:"column:protein_g" json:"proteinG,omitempty"`
	FatG            int    `gorm:"column:fat_g" json:"fatG,omitempty"`
	CarbsG          int    `gorm:"column:carbs_g" json:"carbsG,omitempty"`

	CreatedAt *time.Time `gorm:"column:created_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:datetime not null;default:CURRENT_TIMESTAMP;" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;type:datetime;default:NULL;omitempty;" json:"deleted_at,omitempty"`
}

// Nutrition mirrors types/recipe.ts Nutrition.
type Nutrition struct {
	Calories int `json:"calories"`
	ProteinG int `json:"proteinG"`
	FatG     int `json:"fatG"`
	CarbsG   int `json:"carbsG"`
}

// RecipeCard mirrors types/recipe.ts Recipe. Keep JSON tags in sync with the frontend.
type RecipeCard struct {
	ID              uint      `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Subtitle        string    `json:"subtitle,omitempty"`
	ImageURL        string    `json:"imageUrl,omitempty"`
	Tags            []string  `json:"tags"`
	PrepTimeMinutes int       `json:"prepTimeMinutes"`
	CookTimeMinutes int       `json:"cookTimeMinutes"`
	Servings        int       `json:"servings"`
	Ingredients     []string  `json:"ingredients"`
	Steps           []string  `json:"steps"`
	Nutrition       Nutrition `json:"nutrition"`
}

type GetListOfRecipesResponse struct {
	Name     string `json:"name"`
	RecipeID uint   `json:"recipeId"`
}

func ToRecipeCard(recipe Recipe) RecipeCard {
	return RecipeCard{
		ID:              recipe.ID,
		Slug:            recipeSlug(recipe),
		Title:           recipeTitle(recipe),
		Subtitle:        recipe.Subtitle,
		ImageURL:        recipe.ImageURL,
		Tags:            parseStringSlice(recipe.Tags),
		PrepTimeMinutes: recipe.PrepTimeMinutes,
		CookTimeMinutes: recipe.CookTimeMinutes,
		Servings:        recipe.Servings,
		Ingredients:     parseStringSlice(recipe.Ingredients),
		Steps:           parseStringSlice(recipe.Steps),
		Nutrition: Nutrition{
			Calories: recipe.Calories,
			ProteinG: recipe.ProteinG,
			FatG:     recipe.FatG,
			CarbsG:   recipe.CarbsG,
		},
	}
}

func RecipeFromPayload(input RecipeCard) Recipe {
	title := strings.TrimSpace(input.Title)
	slug := strings.TrimSpace(input.Slug)
	if slug == "" && title != "" {
		slug = slugify(title)
	}

	return Recipe{
		Slug:            slug,
		Title:           title,
		Name:            title,
		Subtitle:        input.Subtitle,
		ImageURL:        input.ImageURL,
		Tags:            marshalStringSlice(input.Tags),
		PrepTimeMinutes: input.PrepTimeMinutes,
		CookTimeMinutes: input.CookTimeMinutes,
		Servings:        input.Servings,
		Ingredients:     marshalStringSlice(input.Ingredients),
		Steps:           marshalStringSlice(input.Steps),
		Calories:        input.Nutrition.Calories,
		ProteinG:        input.Nutrition.ProteinG,
		FatG:            input.Nutrition.FatG,
		CarbsG:          input.Nutrition.CarbsG,
	}
}

func RecipeDisplayName(recipe Recipe) string {
	if name := recipeTitle(recipe); name != "" {
		return name
	}
	return recipe.Name
}

func parseStringSlice(raw string) []string {
	if raw == "" {
		return []string{}
	}

	var parsed []string
	if err := json.Unmarshal([]byte(raw), &parsed); err == nil {
		return parsed
	}

	parts := strings.Split(raw, ";")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func recipeSlug(recipe Recipe) string {
	if recipe.Slug != "" {
		return recipe.Slug
	}
	if recipe.Title != "" {
		return slugify(recipe.Title)
	}
	if recipe.Name != "" {
		return slugify(recipe.Name)
	}
	return fmt.Sprintf("recipe-%d", recipe.ID)
}

func recipeTitle(recipe Recipe) string {
	if recipe.Title != "" {
		return recipe.Title
	}
	return recipe.Name
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	value = strings.ReplaceAll(value, " ", "-")
	var b strings.Builder
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	slug := strings.Trim(b.String(), "-")
	if slug == "" {
		return "recipe"
	}
	return slug
}

func marshalStringSlice(values []string) string {
	if values == nil {
		values = []string{}
	}
	encoded, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(encoded)
}
