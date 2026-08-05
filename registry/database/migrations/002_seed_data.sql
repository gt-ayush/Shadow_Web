INSERT INTO tlds (name) VALUES ('x'), ('web'), ('shop'), ('cloud'), ('mail'), ('dev')
ON CONFLICT (name) DO NOTHING;

INSERT INTO registrars (name, api_key, status)
VALUES ('ShadowRegistrarPrimary', 'sec_key_alpha_99887766', 'active')
ON CONFLICT (name) DO NOTHING;
