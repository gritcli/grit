/path/to/grit completion fish | source

function grit
    set file (mktemp)
    trap "rm -f $file" EXIT
    /path/to/grit --shell-executor-output=$file $argv; and source $file
end
