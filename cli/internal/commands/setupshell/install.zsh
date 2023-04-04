source <(/path/to/grit completion zsh)
compdef _grit grit

grit() {
    local file="$(mktemp)"
    trap "rm -f '$file'" EXIT
    /path/to/grit --shell-executor-output="$file" "$@" && source "$file"
}
