package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/etiedem/pokedexcli/internal/pokecache"
	"github.com/etiedem/pokedexcli/internal/pokeconfig"
)

type Results struct {
	Count     int    `json:"count"`
	Next      string `json:"next"`
	Previous  string `json:"previous"`
	Locations []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}
type ExploreResults struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func get_remote(url string) ([]byte, bool) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, true
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", res.StatusCode, body)
		return nil, true
	}
	if err != nil {
		log.Fatal(err)
		return nil, true
	}
	return body, false
}

func SearchLocation(location string, cache *pokecache.Cache) ExploreResults {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location)
	body, error := cache.Get(url)

	if error == false {

		body, error = get_remote(url)
		cache.Add(url, body)

	}

	results := ExploreResults{}
	err := json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func GetPokemon(name string, cache *pokecache.Cache) pokeconfig.Pokemon {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", name)
	body, error := cache.Get(url)

	if error == false {

		body, error = get_remote(url)
		cache.Add(url, body)

	}

	results := pokeconfig.Pokemon{}
	err := json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results

}

func GetLocations(c *pokeconfig.Config, np string, cache *pokecache.Cache) Results {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=20&limit=20")
	if np == "next" && c.Next != nil {
		url = *c.Next
	}
	if np == "previous" && c.Previous != nil {
		url = *c.Previous
	}
	body, error := cache.Get(url)

	if error == false {

		body, error = get_remote(url)
		cache.Add(url, body)

	}

	results := Results{}
	err := json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	c.Next = &results.Next
	c.Previous = &results.Previous

	return results
}
