# Environment Variables

This document describes the environment variables used by `gritd`.

| Name                | Optionality                    | Description                                         |
| ------------------- | ------------------------------ | --------------------------------------------------- |
| [`GRIT_CONFIG_DIR`] | defaults to `'~/.config/grit'` | the directory containing Grit's configuration files |

⚠️ `gritd` may consume other undocumented environment variables. This document
only shows variables declared using [Ferrite].

## Specification

All environment variables described below must meet the stated requirements.
Otherwise, `gritd` prints usage information to `STDERR` then exits.
**Undefined** variables and **empty** values are equivalent.

The key words **MUST**, **MUST NOT**, **REQUIRED**, **SHALL**, **SHALL NOT**,
**SHOULD**, **SHOULD NOT**, **RECOMMENDED**, **MAY**, and **OPTIONAL** in this
document are to be interpreted as described in [RFC 2119].

### `GRIT_CONFIG_DIR`

> the directory containing Grit's configuration files

The `GRIT_CONFIG_DIR` variable **MAY** be left undefined, in which case the
default value of `~/.config/grit` is used.

```bash
export GRIT_CONFIG_DIR='~/.config/grit' # (default)
```

<!-- references -->

[ferrite]: https://github.com/dogmatiq/ferrite
[`grit_config_dir`]: #GRIT_CONFIG_DIR
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
