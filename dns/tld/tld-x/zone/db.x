$TTL 86400
@   IN  SOA ns1.nic.x. hostmaster.nic.x. (
        1786207528 ; Serial
        3600       ; Refresh
        1800       ; Retry
        604800     ; Expire
        86400 )    ; Minimum TTL

@   IN  NS  ns1.nic.x.

; Delegations for hacker.x
hacker.x. IN NS ns1.hacker.x.
ns1.hacker.x. IN A  10.89.20.99

; Delegations for phantom.x
phantom.x. IN NS ns1.phantom.x.
ns1.phantom.x. IN A  10.89.20.100

; Delegations for shadow.x
shadow.x. IN NS ns1.shadow.x.
ns1.shadow.x. IN A  10.89.20.101

