# shellcheck disable=SC2148

function auto-check() {
  owner=$(git config --get remote.origin.url | cut -d '/' -f 4)
  repo=$(git rev-parse --show-toplevel | xargs basename)

  echo "---> 😎 OWNER: $owner"
  echo "---> 🧰 REPO: $repo"
  echo ""

  if ! update-workflow-names "$repo" "$owner"; then
    return 1
  fi

  if ! update-mod-file "$repo" "$owner"; then
    return 1
  fi

  if ! update-source-id-variable-in-translate-defs "$repo" "$owner"; then
    return 1
  fi

  
  if ! update-astrolib-in-taskfile "$repo" "$owner"; then
    return 1
  fi

  
  if ! update-astrolib-in-goreleaser "$repo" "$owner"; then
    return 1
  fi

  
  if ! rename-templ-data-id "$repo" "$owner"; then
    return 1
  fi

  
  if ! update-import-statements "$repo" "$owner"; then
    return 1
  fi

  
  if ! update-readme "$repo" "$owner"; then
    return 1
  fi

  
  if ! rename-language-files "$repo" "$owner"; then
    return 1
  fi

  
  if ! reset-version; then
    return 1
  fi

  touch ./.env
  echo "✔️ done"
  return 0
}

function update-all-generic() {
  local title=$1
  local repo=$2
  local owner=$3
  local folder=$4
  local name=$5
  local target=$6
  local replacement=$7

  echo "  🎯 --->        title: $title"
  echo "  ✅ ---> file pattern: $name"
  echo "  ✅ --->       folder: $folder"
  echo "  ✅ --->       target: $target"
  echo "  ✅ --->  replacement: $replacement"

  # !!!WARNING: sed
  # does not work the same way between mac and linux
  #
  # find "$folder" -name "$name" -type f -print0: This part of the command uses the find utility to search for files within the specified folder ($folder) matching the given pattern ($name). Here's what each option does:
  #     -name "$name": Specifies the filename pattern to match. In this case, it matches the filenames with the pattern stored in the variable $name.
  #     -type f: Specifies that only regular files should be matched, excluding directories, symbolic links, etc.
  #     -print0: This option tells find to print the matched filenames separated by null characters (\0). This is essential for correctly handling filenames with spaces or special characters.

  # while IFS= read -r -d '' file; do: This initiates a while loop that reads each null-delimited filename (-d '') produced by the find command. Here's what each part does:
  #     IFS=: This sets the Internal Field Separator to nothing. This is to ensure that leading and trailing whitespace characters are not trimmed from each filename.
  #     read -r: This command reads input from the pipe (|) without interpreting backslashes (-r option).
  #     -d '': This option specifies the delimiter as null characters (\0), ensuring that filenames containing spaces or special characters are read correctly.
  #     file: This is the variable where each filename from the find command is stored during each iteration of the loop.
  #
  # trying to modify:
  # - ci-workflow.yml
  # - release-workflow.yml
  find "$folder" -name "$name" -type f -print0 | while IFS= read -r -d '' file; do
    echo "Processing file: $file"
    uname_output=$(uname)
    # sed help:
    # '': no file backup needed (we modify the file in place without backing up original)
    # but note that this is only required for mac. for linux, you dont need the ''.
    #
    # -i: The in-place edit flag, which tells sed to modify the original file inline.
    # 's/search_pattern/replacement_text/g':
    # 
    # s: Indicates that this is a substitution command.
    # /search_pattern/: The pattern to search for.
    # /replacement_text/: The text to replace the search pattern with.
    # g: The global flag, which ensures that all occurrences of the
    #    search pattern are replaced, not just the first one.
    if [[ "$uname_output" == *"Darwin"* ]]; then
      if ! sed -i '' "s/${target}/${replacement}/g" "$file"; then
        echo "!!! ⛔ Sed on mac failed for $file"
        return 1
      fi
    else
      if ! sed -i "s/${target}/${replacement}/g" "$file"; then
        echo "!!! ⛔ Sed on linux failed for $file"
        return 1
      fi
    fi
  done

  echo "  ✔️ --->  DONE"
  echo ""
  return 0
}

function update-mod-file() {
  local repo=$1
  local owner=$2
  local folder=./
  local file_pattern=go.mod
  local target="module github.com\/snivilised\/astrolib"
  local replacement="module github.com\/$owner\/$repo"
  update-all-generic "update-mod-file" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function update-source-id-variable-in-translate-defs() {
  local repo=$1
  local owner=$2
  local folder=./i18n/
  local file_pattern=translate-defs.go
  local target="AstrolibSourceID"
  tc_repo=$(echo "${repo:0:1}" | tr '[:lower:]' '[:upper:]')${repo:1}
  local replacement="${tc_repo}SourceID"
  update-all-generic "update-source-id-variable-in-translate-defs" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function update-astrolib-in-taskfile() {
  local repo=$1
  local owner=$2
  local folder=./
  local file_pattern=Taskfile.yml
  local target=astrolib
  local replacement=$repo
  update-all-generic "update-astrolib-in-taskfile" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function update-workflow-names() {
  local repo=$1
  local owner=$2
  local folder=.github/workflows
  local file_pattern="*.yml"
  local target="name: Astrolib"
  tc_repo=$(echo "${repo:0:1}" | tr '[:lower:]' '[:upper:]')${repo:1}
  local replacement="name: $tc_repo"
  update-all-generic "💥 update-workflow-names" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function update-astrolib-in-goreleaser() {
  local repo=$1
  local owner=$2
  local folder=./
  local file_pattern=.goreleaser.yaml
  local target=astrolib
  local replacement=$repo
  update-all-generic "update-astrolib-in-goreleaser" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function rename-templ-data-id() {
  local repo=$1
  local owner=$2
  local folder=./
  local file_pattern="*.go"
  local target="astrolibTemplData"
  local replacement="${repo}TemplData"
  update-all-generic "rename-templ-data-id" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function update-readme() {
  local repo=$1
  local owner=$2
  local folder=./
  local file_pattern=README.md
  local target="astrolib: "
  local replacement="${repo}: "

  
  if ! update-all-generic "update-readme(astrolib:)" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"; then
    return 1
  fi

  target="snivilised\/astrolib"
  replacement="$owner\/$repo"
  
  if ! update-all-generic "update-readme(snivilised/astrolib)" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"; then
    return 1
  fi

  target="astrolib Continuous Integration"
  tc_repo=$(echo "${repo:0:1}" | tr '[:lower:]' '[:upper:]')${repo:1}
  replacement="$tc_repo Continuous Integration"
  
  if ! update-all-generic "update-readme(astrolib Continuous Integration)" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"; then
    return 1
  fi

  return 0
}

function update-import-statements() {
  local repo=$1
  local owner=$2
  local folder=./
  local file_pattern="*.go"
  local target="snivilised\/astrolib"
  local replacement="$owner\/$repo"
  update-all-generic "update-import-statements" "$repo" "$owner" "$folder" "$file_pattern" "$target" "$replacement"
}

function rename-language-files() {
  local repo=$1
  find . -name 'astrolib*.json' -type f -print0 |
  while IFS= read -r -d '' file; do
    mv "$file" "$(dirname "$file")/$(basename "$file" | sed "s/^astrolib/$repo/")"
  done
  return $?
}

function reset-version() {
  echo "v0.1.0" > ./VERSION
  return 0
}
