#!/usr/bin/env bash

# use https://coolors.co to help create colur schemes
# glyphs:
# - https://github.com/ryanoasis/nerd-fonts/wiki/Glyph-Sets-and-Code-Points
# - https://b.agaric.net/page/agave
# - https://github.com/blobject/agave
# - https://www.nerdfonts.com/cheat-sheet

# <palette>
#   <color name="Jordy Blue" hex="85c0ff" r="133" g="192" b="255" />
#   <color name="Powder blue" hex="94b0da" r="148" g="176" b="218" />
#   <color name="Cool gray" hex="8f91a2" r="143" g="145" b="162" />
#   <color name="Davy's gray" hex="505a5b" r="80" g="90" b="91" />
#   <color name="Outer space" hex="343f3e" r="52" g="63" b="62" />
# </palette>
# use a red/pink colour for error: f72585

# some snippets:
#
# gum confirm "Commit changes?" && git commit -m "$SUMMARY" -m "$DESCRIPTION"
#
# This series of colours does not POP enough
#
# gum style "Hello, there! Welcome to $(gum style --foreground "#85c0ff" 'Jordy Blue')."
# gum style "Hello, there! Welcome to $(gum style --foreground "#94b0da" 'Powder Blue')."
# gum style "Hello, there! Welcome to $(gum style --foreground "#8f91a2" 'Cool gray')."
# gum style "Hello, there! Welcome to $(gum style --foreground "#505a5b" 'Davys gray')."
# gum style "Hello, there! Welcome to $(gum style --foreground "#343f3e" 'Outer space')."

_col_promt="#fbf8cc"
_col_remedy="#fde4cf"
_col_error="#ffcfd2"
_col_branch="#f1c0e8"
_col_action="#cfbaf0"
_col_item="#a3c4f3"
_col_query="#98f5e1"
_col_git="#b9fbc0"
_col_msg="#90dbf4"
_col_affirm="#8eecf5"

# common emojis:
# 😕 😎
# 🍭 🎀 🎁
# ✅
# 🚀 🔆
# 🥥
# 🔥 ⛔ ❌
#

# pastello pallette
#
# <palette>
#   <color name="prompt(Lemon chiffon)" hex="fbf8cc" r="251" g="248" b="204" />
#   <color name="remedy(Champagne pink)" hex="fde4cf" r="253" g="228" b="207" />
#   <color name="error(Tea rose (red))" hex="ffcfd2" r="255" g="207" b="210" />
#   <color name="branch(Pink lavender)" hex="f1c0e8" r="241" g="192" b="232" />
#   <color name="action(Mauve)" hex="cfbaf0" r="207" g="186" b="240" />
#   <color name="item(Jordy Blue)" hex="a3c4f3" r="163" g="196" b="243" />
#   <color name="msg(Non Photo blue)" hex="90dbf4" r="144" g="219" b="244" />
#   <color name="affirm(Electric blue)" hex="8eecf5" r="142" g="236" b="245" />
#   <color name="query(Aquamarine)" hex="98f5e1" r="152" g="245" b="225" />
#   <color name="git-command(Celadon)" hex="b9fbc0" r="185" g="251" b="192" />
# </palette>


# 🍭 gum utility
#

function _text() {
  # when you call this function, the colour must be inside quotes, eg:
  # _text greetings "#98f5e1"
  #
  text=$1
  colour=$2
  gum style --foreground "$colour" "$text"
}

function _a() { # action
  text=$1
  colour=#cfbaf0
  _text "$text" "$colour"
}

function _b() { # branch
  text=$1
  colour=#f1c0e8
  _text "$text" "$colour"
}

function _e() { # error
  text=$1
  colour=#f63e02
  _text "⛔ $text" "$colour"
}

function _g() { # git-command
  text=$1
  colour=#b9fbc0
  _text " $text" "$colour"
}

function _i() { # item
  text=$1
  colour=#a3c4f3
  _text "$text" "$colour"
}

function _m() { # msg
  text=$1
  colour=#90dbf4
  _text "$text" "$colour"
}

function _o() { # ok
  text=$1
  colour=#94fbab
  _text "✅ $text" "$colour"
}

function _p() { # prompt
  text=$1
  colour=#fbf8cc
  _text "😕 $text" "$colour"
}

function _q() { # query
  text=$1
  colour=#98f5e1
  _text "$text" "$colour"
}

function _r() { # remedy
  text=$1
  colour=#fde4cf
  _text "$text" "$colour"
}

function _w() { # warning
  text=$1
  colour=#ffcfd2
  _text "$text" "$colour"
}


# 🍯 git dev workflow commands; This script makes use of nerdfonts.com glyphs, eg 
#

# === 🥥 git-operations ========================================================

function get-def-branch() {
  echo main
}

function gad() {
  if [ -z "$(git status -s -uno | grep -v '^ ' | awk '{print $2}')" ]; then
    gum confirm "$(_text 'Stage all?' $_col_query)" && git add .
  fi

  return 0
}

function get-tracking-branch() {
  git config "branch.$(git_current_branch).remote"
}

function check-upstream-branch() {
  upstream_branch=$(get-tracking-branch)
  feature_branch=$(git_current_branch)

  if [ -z "$upstream_branch" ]; then
    # echo "===> 🐛 No upstream branch detected for : '🎀 $feature_branch'"
    gum style "===> 🐛 No upstream branch detected for : 🎀 $(_text "$feature_branch" $_col_item)"

    if ! _prompt-are-you-sure; then
      return 1
    fi
  fi

  return 0
}

#
# 🍭 end gum utility

# === 🥥 interactive-rebase ====================================================

function _do_start-interactive-rebase() {
  num_commits=$(($1))

  git rebase -i HEAD~"$num_commits"
}

function start-rebase() {
  feature_branch=$(git branch --show-current)
  num_commits=$(git log main.."$feature_branch" --pretty=oneline | wc -l)
  minimum=2
  display_commits="$num_commits"

  if [[ $num_commits -lt $minimum ]]; then
    # whenever time permits, this could be expanded to list all
    # the commits that are to be squashed and put them into
    # a bubbles table.
    # 
    gum style "$(_e "Not enought commits ($(_i "$display_commits")) " \
      "for rebase on branch '$(_b "$feature_branch")', Aborted!")"
  
    return 1
  fi

  _prompt-are-you-sure-with-context \
    "$(gum style "found $(_i "$display_commits") commits on branch '$(_b "$feature_branch")'")" && \
      (_do_start-interactive-rebase "$num_commits" || gum style "$(_e "Aborted!")")

  return 0
}

# === 🥥 gcan(git-commit-amend-no-edit) ========================================

function gcan() {
  feature_branch=$(git branch --show-current)
  num_commits=$(git log main.."$feature_branch" --pretty=oneline | wc -l)
  minimum=1
  display_commits="$num_commits"

  if [[ $num_commits -lt $minimum ]]; then
    gum style "$(_e "Not enought commits, for commit amend on branch '$(_b "$feature_branch")', Aborted!")"

    return 1
  fi

  last=$(git log -n 1 --pretty=format:"%s (HASH: '%h')" "$feature_branch")

  _prompt-are-you-sure-with-context \
    gum style "last commit: $(_i "$last")" && \
      git commit --amend --no-edit -v
}

# === 🥥 prompt ================================================================

function _prompt() {
  message="$1"
  gum confirm "$(_m "$message")"

  return $?
}

function _prompt-are-you-sure {
  _prompt "are you sure? 👾"
  result=$?

  if [ ! "$result" -eq 0 ]; then
    gum style "$(_e "Aborted!")"
  fi

  return $result
}

function _prompt-are-you-sure-with-context {
  context="$1"

  _prompt "$context, are you sure? 👾"
  result=$?

  if [ ! "$result" -eq 0 ]; then
    gum style "$(_e "Aborted!")"
  fi

  return $result
}

# === 🥥 start-feat ============================================================

function start-feat() {
  feature_branch=$1

  if [[ -n $1 ]]; then
    gum style "===> 🚀 $(_a 'START FEATURE'): '🎀 $(_b "$feature_branch")'"

    git checkout -b "$1"
  else
    # echo "!!! 😕 Please specify a feature branch"
    gum style "!!! $(_r 'Please specify a feature branch')"

    return 1
  fi

  return 0
}

# === 🥥 end-feat ==============================================================

function _do-end-feat() {
  feature_branch=$(git_current_branch)
  default_branch=$(get-def-branch)

  if _prompt "About to end feature 🎁 '$feature_branch' ... have you squashed commits"; then
    gum style "<=== ✨ $(_a 'END FEATURE'): $(_b "$feature_branch")"

    if [ "$feature_branch" != master ] && [ "$feature_branch" != main ]; then
      git branch --unset-upstream
      # can't reliably use prune here, because we have in effect
      # a race condition, depending on how quickly github deletes
      # the upstream branch after Pull Request "Rebase and Merge"
      #
      # gcm && git fetch --prune
      # have to wait a while and perform the prune or delete
      # local branch manually.
      #
      git checkout "$default_branch"
      git pull origin "$default_branch"

      gum style "$(_o 'Done!')"
    else
      # echo "!!! 😕 Not on a feature branch ($feature_branch)"
      gum style "!!! $(_e 'Not on a feature branch') ($(_b "$feature_branch"))"

      return 1
    fi
  else
    gum style "$(_e "Aborted!")"

    return 1
  fi

  return 0
}

function end-feat() {
  _prompt-are-you-sure && _do-end-feat
}

# === 🥥 push-feat =============================================================

function _do-push-feat() {
  current_branch=$(git_current_branch)
  default_branch=$(get-def-branch)

  if [ "$current_branch" = "$default_branch" ]; then
    # echo "!!! ⛔ Aborted! still on default branch($default_branch) branch ($current_branch)"
    gum style "!!! $(_e 'Aborted! still on default branch') " \
      "($(_b "$default_branch")) branch ($(_b "$current_branch"))"

    return 1
  fi

  if ! git push origin --set-upstream "$current_branch"; then
    # echo "!!! ⛔ Aborted! Failed to push feature for branch: $current_branch"
    gum style "!!! $(_e 'Aborted! Failed to push feature for branch:') " \
      "$(_b "$current_branch")"

    return 1
  fi

  gum style "$(_o 'pushed feature to') $(_b "$current_branch") ok!"

  return 0
}

function push-feat() {
  _prompt-are-you-sure && _do-push-feat
}

# === 🥥 check-tag =============================================================

function check-tag() {
  rel_tag=$1
  if ! [[ $rel_tag =~ ^[0-9] ]]; then
    # echo "!!! ⛔ Aborted! invalid tag"
    gum style "!!! $(_e 'Aborted! invalid tag')"

    return 1
  fi

  return 0
}

# === 🥥 do-release ============================================================

function _do-release() {
  if [[ -n $1 ]]; then
    if ! check-tag "$1"; then
      return 1
    fi

    raw_version=$1
    version_number=v$1
    current_branch=$(git_current_branch)
    default_branch=$(get-def-branch)

    if [[ $raw_version == v* ]]; then
      # the # in ${raw_version#v} is a parameter expansion operator
      # that removes the shortest match of the pattern "v" from the beginning
      # of the string variable.
      #
      version_number=$raw_version
      raw_version=${raw_version#v}
    fi

    if [ "$current_branch" != "$default_branch" ]; then
      # echo "!!! ⛔ Aborted! not on default($default_branch) branch; current($current_branch)" #error/branch

      gum style "!!! $(_e 'Aborted! not on default')($(_b "$default_branch")) " \
        "branch; current($(_b "$current_branch"))"

      return 1
    fi

    # echo "===> 🚀 START RELEASE: '🎁 $version_number'"
    gum style "===> 🚀 $(_a 'START RELEASE'): '🎁 $(_i "$version_number")'"
    release_branch="release/$version_number"

    if ! git checkout -b "$release_branch"; then
      # echo "!!! ⛔ Aborted! Failed to create the release branch: $release_branch" #error/branch
      gum style "!!! $(_e 'Aborted! Failed to create the release branch:') $(_b "$release_branch")"

      return 1
    fi

    if [ -e ./VERSION ]; then      
      if ! echo "$version_number" > ./VERSION; then
        # echo "!!! ⛔ Aborted! Failed to update VERSION file" #error/item
        gum style "!!! $(_e 'Aborted! Failed to update VERSION file')"

        return 1
      fi

      
      if ! git add ./VERSION; then
        # echo "!!! Aborted! Failed to git add VERSION file" #error/item
        gum style "$(_e "!!! Aborted! Failed to git add $(_i 'VERSION') file")"

        return 1
      fi

      if [ -e ./VERSION-PATH ]; then
        version_path=$(more ./VERSION-PATH)
        echo "$raw_version" > "$version_path"
        
        if ! git add "$version_path"; then
          # echo "!!! ⛔ Aborted! Failed to git add VERSION-PATH file ($version_path)" #error/item
          gum style "!!! $(_e 'Aborted! Failed to git add VERSION-PATH file') ($("_i $version_path"))"

          return 1
        fi
      fi

      if ! git commit -m "Bump version to $version_number"; then
        # echo "!!! ⛔ Aborted! Failed to commit VERSION file" #error/item
        gum style "!!! $(_e 'Aborted! Failed to commit VERSION file')"

        return 1
      fi
      
      if ! git push origin --set-upstream "$release_branch"; then
        # echo "!!! ⛔ Aborted! Failed to push release $version_number upstream" #error/item
        gum style "$(e "!!! Aborted! Failed to push release $(_i "$version_number") upstream")"

        return 1
      fi

      # echo "Done! ✅"
      gum style "$(_o 'Done!')"
    else
      # echo "!!! ⛔ Aborted! VERSION file is missing. (In root dir?)" #error
      gum style "!!! $(_e "Aborted! $(_i 'VERSION') file is missing. (In root dir?)")"

      return 1
    fi
  else
    # echo "!!! 😕 Please specify a semantic version to release" #remedy
    gum style "!!! $(_r 'Please specify a semantic version to release')"

    return 1
  fi

  return 0
}

# release <semantic-version>, !!! do not specify the v prefix, added automatically
# should be run from the root directory otherwise relative paths won't work properly.
function release() {
  _prompt-are-you-sure && _do-release "$1"
}

# === 🥥 tag-rel ===============================================================

function _do-tag-rel() {
  if [[ -n "$1" ]]; then
    version_number="v$1"
    current_branch=$(git_current_branch)
    default_branch=$(get-def-branch)

    if [ "$current_branch" != "$default_branch" ]; then
      # echo "!!! ⛔ Aborted! not on default($default_branch) branch; current($current_branch)" #error/branch
      gum style "$(_e "!!! Aborted! not on default($(_b "$default_branch")) " \
        "branch; current($(_b "$current_branch"))")"

      return 1
    fi

    gum style "$(_a "===> 🏷️  PUSH TAG: '$(_i "$version_number")'")"
    
    if ! git tag -a "$version_number" -m "Release $version_number"; then
      # echo "!!! ⛔ Aborted! Failed to create annotated tag: $version_number" #error
      gum style "$(_e "!!! Aborted! Failed to create annotated tag: $(_i "$version_number")")"

      return 1
    fi

    
    if ! git push origin "$version_number"; then
      # echo "!!! ⛔ Aborted! Failed to push tag: $version_number" #error
      gum style "$(_e "!!! Aborted! Failed to push tag: $(_i "$version_number")")"

      return 1
    fi

    # echo "Done! ✅"
    gum style "$(_o 'Done!')"
  else
    # echo "!!! 😕 Please specify a release semantic version to tag" # remedy
    gum style "$(_r '!!! Please specify a release semantic version to tag')"

    reurn 1
  fi

  return 0
}

# tag-rel <semantic-version>, !!! do not specify the v prefix, added automatically
function tag-rel() {
  _prompt-are-you-sure && _do-tag-rel "$1"
}

#
# 🍯 end git dev workflow commands
