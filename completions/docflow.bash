_docflow_completions() {
  local cur prev words
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  commands="validate compliance scan plan recommend template-sets templates template-impact stats migrate-sections queue history health quickstart"
  case "${COMP_CWORD}" in
    1) COMPREPLY=( $(compgen -W "${commands}" -- "$cur") );;
    *) COMPREPLY=();;
  esac
}
complete -F _docflow_completions docflow
