## Metin2 Deployment Configuration Tool Documentation
This documentation provides detailed instructions on using the Metin2 Deployment Configuration Tool, focusing on arguments, preprocessing YAML configurations, available template functions, and custom template variables.  test

## Usage

To use the Deployer, you have to first build it using `build.sh` script after installing [**golang**](https://go.dev/doc/install "golang").

You can then create all template files in `command`, `root_command` and `additional_folders` using **go-template** syntax.

When creating a template, you can define as many variables as you want, the program is not aware about the meaning of your variables, it will only take yaml files in input and procss your templates accordingly.

It is up to you to give yaml's properties a meaning.

Finally run `deploy.sh` script specifiyng yamls, output path and channel count.

### Command-Line Arguments
The program accepts several command-line arguments for configuration and deployment:

- --config (string slice): Paths to the configuration files.
- --out (string): Defines the output path for the deployment. Defaults to deploy.
- --channels (int): Defines the amount of channels to generate. This parameter is mandatory.

Example Usage
```
./deploy.sh --config=config1.yaml,config2.yaml --out=output_folder --channels=5
```

## Deployer configuration

The deployer uses **yaml**s and **go-template** to represent and generate all needed folders to run a Metin2 Server.

It allows for multiple yaml files to be ingested and it will take care of merging them in a single configuration file, allowing to create a `base_config` file, and multiple specific yamls (one for test, one for prod, ...).

The process is divided in 6 steps:

- **Preprocess of the yamls**, injecting the arguments passed to the program such as `channels`. In this phase, all the yaml files, possibly containining go-template syntax, are *exploded*, allowing for more cohincise yamls. One of the examples on how to exploit the preprocessing, is to generate `ports`, `p2p_ports` and `hostnames` following some simple rules and using `channels` argument to generate all ports dynamically.

- **Merging of the yamls**. All yamls are merged in a single configuration file

- **Process of command, root_command and any additional folder template**. in this step all scripts and files in the `command` and `root_command` are procesed using **go-template** to generate the files which will be placed in the server folders.

- **Directory creation and file touching**. In this step the deployer can create directories and touch files, if specified in the `metadata` section

- **Run of install script**. If specified in the `metadata` section, the deployer will run the install script to install all cores.

### Yaml structure

Below is an example for a configuration yaml

```yaml
metadata:
  create_dirs:
    - common/package
    - common/data
    - common/locale/italy
    - common/prototypes
  touch:
    - common/VERSION
  install_script:  install.sh

globals:
  compiled_bin_dir: /
  test_server: false
  root_dir: /
  connection.ip: 172.21.0.22
  connection.proxy_ip: 172.18.16.1
  cycles.passes_per_sec: 25
  cycles.save_seconds: 180
  cycles.ping_seconds: 180
  table_postfix: ""
  locale.service: italy
  view_range: 20000

cores:
  db:
    renames: 
      config_template: conf.txt
    values:
      hostname: db
      path: db
      table_postfix: ""
      connection.port: 15000
      sleep_ms: 10
      heart_fps: 25
      hash_player_life_sec: 600
      backup_limit_sec: 3600
      player_id_start: 100
      player_delete_level_max: 101
      player_delete_level_min: 0
      player_delete_check: 0
      locale.encoding: big5
      item_id_range:
        - 20000001
        - 400000000

  auth:
    renames: 
      config_template: CONFIG
    values:
      path: auth
      channel: 1
      hostname: auth
      connection.port: 25045
      connection.p2p_port: 35045
      auth.enabled: 1

  game99:
    renames: 
      config_template: CONFIG
    values:
      path: game99
      channel: 99
      hostname: game99
      connection.port: 31099
      connection.p2p_port: 40099
      football.master: true
      mark.master: true
      maps_allowed:
        - 81
        - 110
        - 111
        - 113
        - 114
        - 118
        - 119
        - 120
        - 121
        - 122
        - 123
        - 124
        - 125
        - 126
        - 127
        - 128
        - 181
        - 182
        - 183
        - 227
        - 242
        - 246

  
  {{ range $i := (seq 0 $.channels) }}
  {{ $ch := add $i 1 }}
  game1-ch{{ $ch }}:
    renames: 
      config_template: server.xml
    values:
      path: channel{{ $ch }}/game1
      channel: {{ $ch }}
      hostname: game1-ch{{ $ch }}
      connection:
        port:  {{ add 10201 (mul $i 200) }}
        p2p_port:  {{ add 17300 (mul $i 200) }}
      lrdb.port: {{ add 10100 (mul $i 1) }}
      mark.enabled: {{ if eq $ch 1 }} 1 {{ else }} 0 {{ end }}
      mark.min_level: {{ if eq $ch 1 }} 4 {{ else }} 0 {{ end }}
      intel.is_master: {{ if eq $ch 1 }} 1 {{ else }} 0 {{ end }}
      maps_allowed:
        - 61
        - 63
        - 64
        - 67
        - 68
        - 71
        - 73
        - 172
        - 100
        - 167
        - 220
      rejoinable_dungeon:
        - 100
        - 167
        - 220
{{ end }}

additional_folders:
  iptables:
    out_path: iptables
    values: {}
```

### Metadata
- **create_dirs**: Directories to be created as part of the deployment.
- **touch**: Files to be created (if they do not exist) or updated.
- **install_script**: Script to be run for installation.

### Globals
Global parameters that apply to all configurations and additional_folders.

For `command` folder, globals are injected after yaml merging and contain only the values specified in the yaml.

For `root_command` and `additional_folders` globals are injected after all cores have been processed and it contains 2 additional properties:

| Property | Description |
| - | - |
| __out__cores | The list of all process cores with all computed values | 
| paths | The list of all cores paths after being processed |

### Cores
Specific configurations for different core components, defined under the `cores` section. Each core configuration can have:
- **renames**: Key-value pairs for renaming configuration elements.
- **values**: Specific configuration values for the core.

### Additional Folders
Defines additional folders to be created and their output paths:
- **out_path**: Path where additional folders should be created.
- **values**: Configuration values related to the additional folder.

### YAML Preprocessing
Before merging, the YAML configurations are preprocessed with the provided variables. This step ensures that the configurations are correctly combined.

During preprocessing Custom Functions Available in Templates
The following custom functions are available for use in your templates:

| Template Function | Description | Example |
| - | - | - | 
| add | Adds two numbers. | `{{ add 3 5 }}` |
| sub | Substracts two numbers. | `{{ sub 3 5 }}` |
| mul | Multiplies two numbers. | `{{ mul 3 5 }}` |
| seq | Generates a sequence of integers from start to end. | `range $i := (seq 0 $.channels)` |
| channels |  Channel number passed as argument to the program | `channels: {{ .channels }}`

### YAML processing
After preprocessing, the deployer processes all files in `configurations`. files here can be written with the standard `go-template` function.
Foreach `core` generated by the **preprocessor**, the deployer process all files in the folder and generates a core at the **path** specified in `path` field.

The file can use fields inside the `core configuration` with root in `Values`

Example of file:

```xml
<server>
	<!-- MAIN CONFIGURATION -->
	<channel>{{ .channel }}</channel>
	<hostname>{{ .hostname }}</hostname>
	
	<!-- CONNECTION CONFIGURATION -->
	<bind_ip>{{ .connection.ip }}</bind_ip>
	{{ if .connection.proxy_ip }}<proxy_ip>{{ .connection.proxy_ip }}</proxy_ip>
	{{ end }}
	<port>{{ .connection.port }}</port>
	<p2p_port>{{ .connection.p2p_port }}</p2p_port>
	
    <!-- AUTH CONFIGURATION -->
	<auth>
		<enabled>{{ or .auth.enabled 0 }}</enabled>
	</auth>
```

When finally processing the templates in `configurations` folder, there are some other utility functions available

| Template Function | Description | Example |
| - | - | - |
| excludeStrings | Filters out strings from a slice that do not contain a specified word. | ```{{ $filtered := excludeStrings .Values "excludeWord" }}```
| filterStrings | Filters strings in a slice that contain a specified word. | ```{{ $filtered := filterStrings .Values "includeWord" }}``` |

The below template iterates over all `paths` excluding the one containing `db`
```sh
ROOT=$PWD
{{ range excludeStrings .paths "db" }}
echo "Clearing {{ . }}"
cd $ROOT/{{ . }} && sh clear.sh
{{end}}

sleep 2
```

The below template iterates over all `paths` which contain `db`


```sh
ROOT=$PWD

echo "Clearing {{ index (filterStrings .paths "db") 0 }}"
cd $ROOT/{{ index (filterStrings .paths "db") 0 }} && sh clear.sh

sleep 2
```