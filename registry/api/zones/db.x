$TTL 86400
@   IN  SOA ns1.nic.x. hostmaster.nic.x. (
        2026080508 ; Serial
        3600       ; Refresh
        1800       ; Retry
        604800     ; Expire
        86400 )    ; Minimum TTL

@   IN  NS  ns1.nic.x.

; Delegations for example.x
example.x. IN NS ns1.example.x.
ns1.example.x. IN A  10.89.20.10

