package utils

import (
	"encoding/json"
	"os"

	"github.com/renja-g/Barkeeper/constants"
)

func SaveProfiles(profiles []constants.Profile) error {
	// Open a new file for writing only
	file, err := os.OpenFile("data/profiles.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the profiles to the JSON format and write it to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(profiles)
	if err != nil {
		return err
	}

	return nil
}

func GetProfiles() ([]constants.Profile, error) {
	// Open the file for reading only
	file, err := os.Open("data/profiles.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data from the file
	var profiles []constants.Profile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func SaveMatches(matches []constants.Match) error {
	// Open a new file for writing only and truncate it
	file, err := os.OpenFile("data/matches.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the matches to the JSON format and write it to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(matches)
	if err != nil {
		return err
	}

	return nil
}

func GetMatches() ([]constants.Match, error) {
	// Open the file for reading only
	file, err := os.Open("data/matches.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data from the file
	var matches []constants.Match
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&matches)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
