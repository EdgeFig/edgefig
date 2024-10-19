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
          - 203.0.113.2/30
          - 2001:db8:40:1::2/126
      eth1:
        name: LAN
        addresses:
          - 192.0.2.0/24
        mtu: 9000
        ipv6:
          nameserver: 2606:4700:4700::1111
          prefixes:
            - prefix: 2001:db8:40:2::/64
              autonomous: true
        vlans:
          - examplevlan
    firewall:
      groups: {}
      zones:
        - name: WAN_IN
          ip-type: ipv4
          default-action: drop
          description: Traffic inbound on WAN port to somewhere else in the network
          in:
            - eth0
          rules:
            - action: accept
              description: Allow established/related
              destination:
                prefix: 192.0.2.0/24
              log: disable
              protocol: all
              established: enable
              invalid: disable
              new: disable
              related: enable

            - action: accept
              description: Allow ICMP
              log: disable
              protocol: icmp

        - name: WAN_LOCAL
          ip-type: ipv4
          default-action: drop
          description: Traffic inbound on WAN port to the local router
          local:
            - eth0
          rules:
            - action: accept
              description: Allow established/related
              protocol: all
              established: enable
              related: enable

            - action: drop
              description: Drop invalid state
              protocol: all
              invalid: enable

            - action: accept
              description: Enable Ping
              protocol: icmp

        - name: WAN_IN_6
          ip-type: ipv6
          default-action: drop
          description: ipv6 traffic inbound on WAN port to somewhere else in the network
          in:
            - eth0
          rules:
            - action: accept
              log: disable
              protocol: icmpv6

            - action: accept
              established: enable
              related: enable

        - name: WAN_LOCAL_6
          ip-type: ipv6
          default-action: drop
          description: ipv6 traffic inbound on WAN port to the local router
          local:
            - eth0
          rules:
            - action: accept
              description: ICMPv6
              protocol: icmpv6

            - action: accept
              description: "Allow related & established"
              protocol: all
              established: enable
              related: enable
    bgp:
      - asn: 65536 # This is our ASN
        router-id: 203.0.113.2
        peers:
          # ISPv4
          - ip: 203.0.113.1
            source-ip: 203.0.113.2 # If you need the BGP session to originate from a specific IP
            asn: 65537 # This is our peer's ASN
            announce-default: false
            announcements:
              - 1.2.3.4/24 # Announce this block to the peer
            accept:
              - prefix: 0.0.0.0/0
                le: 24
          # ISPv6
          - ip: 2001:db8:40:1::1
            source-ip: 2001:db8:40:1::2
            asn: 65537 # This is our peer's ASN
            announce-default: false
            announcements:
              - 2001:db8:40::/48
            accept:
              - prefix: ::/0
                le: 64
    routes:
      - description: Default Route ipv4
        route: 0.0.0.0/0
        next-hop: 203.0.113.1
        distance: 1
      - description: Default Route ipv6
        route: ::/0
        next-hop: 2001:db8:40:1::1
        distance: 1
    dhcp:
      - name: LAN
        authoritative: true
        subnet: 192.0.2.0/24
        router: 192.0.2.1
        start: 192.0.2.150
        stop: 192.0.2.254
        lease: 86400
        dns:
          - 1.1.1.1
          - 8.8.8.8
        reservations:
          - name: example-reservation-name
            mac: 00:00:5E:00:53:00
            ip: 192.0.2.21
      - name: examplevlan
        authoritative: true
        subnet: 198.51.100.0/24
        router: 198.51.100.1
        start: 198.51.100.50
        stop: 198.51.100.254
        lease: 86400
        dns:
          - 1.1.1.1
          - 8.8.8.8
        reservations: []
    dns:
      forwarding:
        cache-size: 150
        listen-on:
          - eth1
        nameservers:
          - 1.1.1.1
          - 8.8.8.8
    nat:
      - name: Destination NAT Example
        type: destination
        inbound_interface: eth1
        protocol: all
        inside_address:
          address: 192.0.2.21
        outside_address:
          address: 203.0.113.7
      - name: Source NAT Example
        type: source
        outbound_interface: eth1
        protocol: all
        inside_address:
          address: 192.0.2.21
        outside_address:
          address: 203.0.113.7
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

# VLANs are defined here and assigned by name to routers, switch ports, etc
vlans:
  - name: examplevlan
    id: 15
    address: 198.51.100.0/24
    mtu: 9000

```

## Apply

Once your configuration is written, you can apply the configuration against all devices by running `edgefig apply`
