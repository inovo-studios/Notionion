#!/bin/bash

echo -n "Enter integration token: " 
read -s token
echo
echo -n "Enter notion page id: "
read pageid

export NOTION_TOKEN="$token"
export NOTION_PAGEID="$pageid"