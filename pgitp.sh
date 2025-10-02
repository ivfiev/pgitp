#!/bin/bash

cmd="$1"
local="$2"
remote="$3"

ask_password() {
  read -s -p "Password: " PASSWORD
  echo
}

pull() {
  git -C $remote pull   
  ask_password
  for file in $remote/*; do
    fname="${file##*/}"
    cat $file | gpg --decrypt --passphrase $PASSWORD --batch | tar -xzf - > $local$fname
  done
  unset PASSWORD
}

push() {
  ask_password
  for file in $local/*; do
    fname="${file##*/}"
    hash1=$(cat $remote$fname | gpg --decrypt --passphrase $PASSWORD --batch | tar -xzOf -)
    hash2=$(cat $file)
    if [[ "$hash1" != "$hash2" ]]; then
      tar -czf - $file | gpg --symmetric --cipher-algo AES256 --armor --batch --yes --passphrase $PASSWORD > $remote$fname
    fi
  done
  unset PASSWORD
  git -C $remote add -A
  git -C $remote commit -m "WIP" 
  git -C $remote push
}

case "$1" in
  "pull") pull;;
  "push") push;;
  *) echo "fuck you!";;
esac
