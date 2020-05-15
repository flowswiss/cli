# Usage Documentation

## Global Flags
```
  --endpoint-url string     Base url to the api endpoints for the platform. The
                            default is the Flow Swiss platform, but this flag
                            can be used to change to any other environment. Can
                            also be specified through the `FLOW_ENDPOINT_URL`
                            environment variable or the `endpoint_url` item in
                            your configuration file.

                            Example: 'https://api.cloudbit.ch/'


  --username string         The username to authenticate against the Flow Swiss
                            platform. Can also be specified through the
                            `FLOW_USERNAME` environment variable or the
                            `username` item in your configuration file.

                            Example: 'example@flow.swiss'


  --password string         The password for the specified user on the Flow
                            Swiss platform. Can also be specified through the
                            `FLOW_PASSWORD` environment variable or the
                            `password` item in your configuration file.

                            Example: 'MySuperSecurePassword'


  --two-factor-code string  The current two factor code for your account. By
                            default it will ask for user input during runtime
                            when you have two factor authentication enabled on
                            your account. 

                            Example: '123456'

  --organization string     Filter for the organization to run all API Requests
                            in. Can be an interger if you want to search by id
                            or a string if you want to search for any other
                            identification (e.g. name). Can also be specified
                            through the `FLOW_ORGANIZATION` environment variable
                            or the `organization` item in your configuration
                            file.

                            Example: 'Flow Swiss AG'


  --format string           The output format to use for list operations. Can be
                            one of table, csv or json. By default all list
                            operations will generate a table with the most
                            important information. If you want to integrate the
                            response into an other system, we suggest using the
                            json option. Can also be specified through the
                            `FLOW_FORMAT` environment variable or the `format`
                            item in your configuration file.

                            Example: 'json'


  --verbosity level         Specifies the verbose output level for the current
                            exectuion. Hhe higher the level, the more output
                            will be generated. This can idealy be used for
                            debugging purposes. The verbose output will be
                            posted to `stderr` to still allow pipes.

                            Level 0: default output
                            Level 1: show request paths
                            Level 2: dump responses from the api
                            Level 3: dump requests to the api

                            Example: '1'
```

## Commands

### Authentication

```
flow auth login
  Checks whether your login credentials are valid and stores the username and
  password in your users home directory under `$HOME/.flow/credentials.json`
```

### Compute

```
flow compute server list
  Lists all server of the selected organization.
```

```
flow compute server create
  Creates a new virtual machine for the selected organization

  --name string (required)  The name for the server. This name will also be used
                            as the hostname of the virtual machine, if you want
                            this to match the name on the portal choose only
                            hostname allowed characters (a-z, 0-9 and hyphen).

                            Example: 'My First Virtual Machine'


  --location string (required)  Filter for the desired location of the server.
                            Can be an interger if you want to search by id or a
                            string if you want to search for any other
                            identification (e.g. name or city). To get a list of
                            available locations see `flow locations --module
                            compute`

                            Example: 'alp1'


  --image string (required)  Filter for the desired operating system image. Can
                            be an integer if you want to search by id or a
                            string if you want to search for any other
                            identification (e.g. name)

                             Example: 'ubuntu-20.04'


  --product string (required)  Filter for the desired product to use. Can be an
                            interger if you want to search by id or a string if
                            you want to search for any other identification
                            (e.g. name). To get a list of available products see
                            `flow products compute`

                            Example: 'b1.1x1'


  --network string          Filter for the desired network for the first network
                            interface of the virtual machine. By default it will
                            select your first network automatically. Can be an
                            integer if you want to search by id or a string if
                            you want to search for any other identification
                            (e.g. name or cidr)

                            Example: 'Default Network'


  --private-ip string       An optional private ip address to assign to the
                            first network interface. Must be inside the cidr of
                            the network and must not be assigned to any other
                            item in the network. By default dhcp will choose an
                            address automatically.

                            Example: '172.31.0.11'


  --key-pair string (required if image is linux)
                            An initial ssh keypair for the virtual machine to
                            connect to via ssh. Can be an integer if you want to
                            search by id or a string if you want to search for
                            any other identification (e.g. name or md5
                            fingerprint). If you do not have a key, you can
                            generate one by typing `ssh-keygen`.

                            Example: '12:ab:9c:de:93:01:6d:7a:b7:2d:27:06:65:90:c9:cf'


  --windows-password string (required if image is windows)
                            An initial windows password for the windows machine.
                            Can be anything, but should match the windows
                            password requirements (at least: one uppercase
                            character, one lowercase character, one digit, one
                            special character)

                            Exmaple: 'z2.xHfQFi8vM27t,o9Ft' (please choose a
                                     different one)


  --cloud-init file         A cloud init script for the setup of your virtual
                            machine. To see how cloud init scripts work, please
                            visit https://cloudinit.readthedocs.io/en/latest/

                            Example: './cloud-init.yaml'


  --attach-external-ip bool  Whether to attach an elastic ip (public ip) to the
                            virtual machine or not. By default an elastic ip
                            will be attached. To disable use `--attach-external-
                            ip=false`  
```


```
flow compute server [start|stop|reboot] <server>
  Perform an action on the virtual machine. The command will block until the
  server status updates to the expected value.

  <server> string           Filter for the server to execute the action on. Can
                            be an integer if you want to search by id or a
                            string if you want to search for any other
                            identification (e.g. name, public ip, private ip,
                            etc)

                            Example: 'My First Virtual Machine'
                   
```

```
flow compute server delete <server>
  Detaches all elastic ips attached to the selected server and deletes the
  server product itself.

  <server> string           Filter for the server to delete. Can be an integer
                            if you want to search by id or a string if you want
                            to serarch for any other identification (e.g. name,
                            public ip, private ip, etc)

                            Example: 'My First Virtual Machine'


  --force                   If this flag is present, the interaction with the
                            user whether he really wants to delete the server
                            will be skipped and the server will be deleted
                            immediately


  --detach-only             By default all attached elastic ips will be deleted
                            once they were detached. This can be prevented
                            through this flag and all elastic ips will only be
                            detached.
```

```
flow compute key-pair list
  Lists all key pairs of the selected organization.
```

```
flow compute key-pair upload <file>
  Uploads the selected public key if it does not already exist.

  <file> string             Path to the public key file which should be
                            uploaded. The selected file must be in the openssh
                            authorized keys format as generated through
                            `ssh-keygen`

  --name string             The name for the new key pair. By default the
                            optional public key comment is chosen if it exists.
```

```
flow compute network list
  Lists all networks of the selected organization.
```

### Products

```
flow products compute
  Lists all compute virtual machine products.

  --location string         Filter for the location desired where the product
                            should be available. Can be an integer if you want
                            to search by id or a string if you want to serarch
                            for any other identification (e.g. name or city)

                            Example: 'ALP1'
```