#compdef docflow
_docflow() {
  local -a commands
  commands=(
    "validate:validate docs"
    "compliance:governance compliance"
    "scan:scan repo"
    "plan:daily plan"
    "recommend:template recommendations"
    "template-sets:template co-usage"
    "templates:list/deprecate"
    "template-impact:impact of template changes"
    "stats:section coverage"
    "migrate-sections:alias migration"
    "queue:queue go/no-go"
    "history:run history"
    "health:healthcheck"
    "quickstart:demo run"
  )
  _describe 'command' commands
}
compdef _docflow docflow
