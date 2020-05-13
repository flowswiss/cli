# Usage Documentation

## Global Flags
```
  --endpoint-url string     base url to the api endpoints for the platform. the
                            default is the Flow Swiss platform, but this flag
                            can be used to change to any other environment. can
                            also be specified through the `FLOW_ENDPOINT_URL`
                            environment variable or the `endpoint_url` item in
                            your configuration file.

                            example: 'https://api.cloudbit.ch/'


  --username string         the username to authenticate against the Flow Swiss
                            platform. can also be specified through the
                            `FLOW_USERNAME` environment variable or the
                            `username` item in your configuration file.

                            example: 'example@flow.swiss'


  --password string         the password for the specified user on the Flow
                            Swiss platform. can also be specified through the
                            `FLOW_PASSWORD` environment variable or the
                            `password` item in your configuration file.

                            example: 'MySuperSecurePassword'


  --two-factor-code string  the current two factor code for your account. by
                            default it will ask for user input during runtime
                            when you have two factor authentication enabled on
                            your account. 

                            example: '123456'


  --format string           the output format to use for list operations. can be
                            one of table, csv or json. by default all list
                            operations will generate a table with the most
                            important information. if you want to integrate the
                            response into an other system, we suggest using the
                            json option. can also be specified through the
                            `FLOW_FORMAT` environment variable or the `format`
                            item in your configuration file.

                            example: 'json'


  --verbosity level         specifies the verbose output level for the current
                            exectuion. the higher the level, the more output
                            will be generated. this can idealy be used for
                            debugging purposes. the verbose output will be
                            posted to `stderr` to still allow pipes.

                            level 0: default output
                            level 1: show request paths
                            level 2: dump responses from the api
                            level 3: dump requests to the api

                            example: '1'
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

  --name string (required)  the name for the server. this name will also be used
                            as the hostname of the virtual machine, if you want
                            this to match the name on the portal choose only
                            hostname allowed characters (a-z, 0-9 and hyphen).

                            example: 'My First Virtual Machine'


  --location string (required)  filter for the desired location of the server.
                            can be an interger if you want to search by id or a
                            string if you want to search for any other
                            identification (e.g. name or city). to get a list of
                            available locations see `flow locations --module
                            compute`

                            example: 'alp1'


  --image string (required)  filter for the desired operating system image. can
                            be an integer if you want to search by id or a
                            string if you want to search for any other
                            identification (e.g. name)

                             example: 'ubuntu-20.04'


  --product string (required)  filter for the desired product to use. can be an
                            interger if you want to search by id or a string if
                            you want to search for any other identification
                            (e.g. name). to get a list of available products see
                            `flow products --module compute`

                            example: 'b1.1x1'


  --network string          filter for the desired network for the first network
                            interface of the virtual machine. by default it will
                            select your first network automatically. can be an
                            integer if you want to search by id or a string if
                            you want to search for any other identification
                            (e.g. name or cidr)

                            example: 'Default Network'


  --private-ip string       an optional private ip address to assign to the
                            first network interface. must be inside the cidr of
                            the network and must not be assigned to any other
                            item in the network. by default dhcp will choose an
                            address automatically.

                            example: '172.31.0.11'


  --key-pair string (required if image is linux)
                            an initial ssh keypair for the virtual machine to
                            connect to via ssh. can be an integer if you want to
                            search by id or a string if you want to search for
                            any other identification (e.g. name or md5
                            fingerprint). if you do not have a key, you can
                            generate one by typing `ssh-keygen`.

                            example: '12:ab:9c:de:93:01:6d:7a:b7:2d:27:06:65:90:c9:cf'


  --windows-password string (required if image is windows)
                            an initial windows password for the windows machine.
                            can be anything, but should match the windows
                            password requirements (at least: one uppercase
                            character, one lowercase character, one digit, one
                            special character)

                            exmaple: 'z2.xHfQFi8vM27t,o9Ft' (please choose a
                                     different one)


  --cloud-init file         a cloud init script for the setup of your virtual
                            machine. to see how cloud init scripts work, please
                            visit https://cloudinit.readthedocs.io/en/latest/

                            example: './cloud-init.yaml'


  --attach-external-ip bool  whether to attach an elastic ip (public ip) to the
                            virtual machine or not. by default an elastic ip
                            will be attached. to disable use `--attach-external-
                            ip=false`  
```


```
flow compute server [start|stop|reboot] <server>
  Perform an action on the virtual machine. The command will block until the
  server status updates to the expected value.

  <server> string           filter for the server to execute the action on. can
                            be an integer if you want to search by id or a
                            string if you want to search for any other
                            identification (e.g. name, public ip, private ip,
                            etc)

                            example: 'My First Virtual Machine'
                   
```

```
flow compute server delete <server>
  Detaches all elastic ips attached to the selected server and deletes the
  server product itself.

  <server> string           filter for the server to delete. can be an integer
                            if you want to search by id or a string if you want
                            to serarch for any other identification (e.g. name,
                            public ip, private ip, etc)

                            example: 'My First Virtual Machine'


  --force                   if this flag is present, the interaction with the
                            user whether he really wants to delete the server
                            will be skipped and the server will be deleted
                            immediately


  --detach-only             by default all attached elastic ips will be deleted
                            once they were detached. this can be prevented
                            through this flag and all elastic ips will only be
                            detached. 
```