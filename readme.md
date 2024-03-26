# EdgeFig

EdgeFig is a tool for managing EdgeMax devices (EdgeRouter, EdgeSwitch) in a declarative fashion in a simplified yaml configuration, which can easily be applied to all managed devices from a central location. The tool was born out of a desire to manage these routers and switches with a gitops flow, where changes flow through a pull request flow with required approvals and are then automatically applied on merge.

This is currently a work in progress and rapidly evolving.

## Configuration

Configuration for the devices is a simplified yaml file, which gets expanded to the full configuration that is applied to the EdgeMax devices.

```yaml
routers:
  - name: router01
    connection:
      ip: 10.0.0.1
      port: 22
      username: ubnt
      password: ubnt
    # If an interface is not specifically listed here, it will be assumed disabled
    interfaces:
      eth0:
        name: WAN
        addresses:
          - 10.100.1.1/24
      eth1:
        name: LAN
        addresses:
          - 10.0.0.1/24
    dhcp:
      - name: LAN
        authoritative: true
        subnet: 10.0.0.0/24
        router: 10.0.0.1
        start: 10.0.0.150
        stop: 10.0.0.254
        dns:
          - 1.1.1.1
          - 8.8.8.8
    nat:
      - name: Masquerade for WAN
        type: masquerade
        outbound_interface: eth0
        protocol: all
        log: false
    users:
      - username: ubnt
        password: ubnt
        role: admin

#switches:
#  - name: switch01

#firewall:
#  - <rule>

#vlans:
#  -

```

## Apply

Once your configuration is written, you can apply the configuration against all devices by running `edgefix apply`
