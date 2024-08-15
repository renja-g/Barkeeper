#!/bin/bash

# Copy the example config file to config.json
cp example.config.json config.json

# Function to prompt for a value and replace it in config.json
prompt_and_replace() {
    local key=$1
    local placeholder="FILL_ME_IN"
    local prompt_text=$2
    
    read -p "$prompt_text: " value
    # Escape special characters in value
    value=$(echo $value | sed 's/[&/\]/\\&/g')
    # Replace the placeholder in config.json
    sed -i "s/\"$key\": \"$placeholder\"/\"$key\": \"$value\"/" config.json
}

# Prompt the user for each value
prompt_and_replace "dev_guild_id" "Enter the Development Guild ID"
prompt_and_replace "token" "Enter the Token"
prompt_and_replace "blue_channel_id" "Enter the Blue Channel ID"
prompt_and_replace "red_channel_id" "Enter the Red Channel ID"
prompt_and_replace "lobby_channel_id" "Enter the Lobby Channel ID"
prompt_and_replace "riot_api_key" "Enter the Riot API Key"

echo "Configuration complete! Your config.json file has been created."
