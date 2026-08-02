$TTL 86400
@   IN  SOA ns1.dns.x. admin.dns.x. (
            2026080301 ; serial
            3600       ; refresh
            1800       ; retry
            604800     ; expire
            86400 )    ; minimum

@   IN  NS  ns1.dns.x.
@   IN  NS  ns2.dns.x.

; Glue records for .x TLD NS itself
ns1.dns.x.   IN  A   10.89.10.20
ns2.dns.x.   IN  A   10.89.10.21

; --- Subdomain / SLD Delegations ---
example.x.  IN  NS  ns1.example.x.
ns1.example.x. IN A 10.89.20.10
