source <(/path/to/grit completion bash)

grit() {
    local file="$(mktemp)"
    trap "rm -f '$file'" EXIT
    /path/to/grit --shell-executor-output="$file" "$@" && source "$file"
}
