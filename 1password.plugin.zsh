# make sure you execute this *after* asdf or other version managers are loaded
if (( $+commands[op] )); then
  eval "$(op completion zsh)"
  compdef _op op
  
  # load plugins configuration
  if [[ -f ~/.config/op/plugins.sh ]]; then
    source ~/.config/op/plugins.sh
  fi
fi