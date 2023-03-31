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

## Usage Examples

<details>
<summary>Kubernetes</summary>

This example shows how to define the environment variables needed by `gritd`
on a [Kubernetes container] within a Kubenetes deployment manifest.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deployment
spec:
  template:
    spec:
      containers:
        - name: example-container
          env:
            - name: GRIT_CONFIG_DIR # the directory containing Grit's configuration files (defaults to '~/.config/grit')
              value: ~/.config/grit
```

Alternatively, the environment variables can be defined within a [config map][kubernetes config map]
then referenced from a deployment manifest using `configMapRef`.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: example-config-map
data:
  GRIT_CONFIG_DIR: ~/.config/grit # the directory containing Grit's configuration files (defaults to '~/.config/grit')
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: example-deployment
spec:
  template:
    spec:
      containers:
        - name: example-container
          envFrom:
            - configMapRef:
                name: example-config-map
```

</details>

<details>
<summary>Docker</summary>

This example shows how to define the environment variables needed by `gritd`
when running as a [Docker service] defined in a Docker compose file.

```yaml
service:
  example-service:
    environment:
      GRIT_CONFIG_DIR: ~/.config/grit # the directory containing Grit's configuration files (defaults to '~/.config/grit')
```

</details>

<!-- references -->

[docker service]: https://docs.docker.com/compose/environment-variables/#set-environment-variables-in-containers
[ferrite]: https://github.com/dogmatiq/ferrite
[`grit_config_dir`]: #GRIT_CONFIG_DIR
[kubernetes config map]: https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/#configure-all-key-value-pairs-in-a-configmap-as-container-environment-variables
[kubernetes container]: https://kubernetes.io/docs/tasks/inject-data-application/define-environment-variable-container/#define-an-environment-variable-for-a-container
[rfc 2119]: https://www.rfc-editor.org/rfc/rfc2119.html
