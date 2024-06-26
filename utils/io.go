package utils

import (
	"encoding/json"
	"os"

	"github.com/renja-g/Barkeeper/constants"
)

func SaveRatings(ratings []constants.Rating) error {
	// Open a new file for writing only
	file, err := os.OpenFile("data/ratings.json", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode the ratings to the JSON format and write it to the file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	err = encoder.Encode(ratings)
	if err != nil {
		return err
	}

	return nil
}

func GetRatings() ([]constants.Rating, error) {
	// Open the file for reading only
	file, err := os.Open("data/ratings.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decode the JSON data from the file
	var ratings []constants.Rating
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ratings)
	if err != nil {
		return nil, err
	}

	return ratings, nil
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
