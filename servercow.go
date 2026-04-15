// Package servercow implements a Caddy DNS provider module for Servercow.
//
// This module wraps github.com/libdns/servercow and registers it with Caddy's
// module system so it can be used to solve ACME DNS-01 challenges.
//
// Caddyfile syntax:
//
//	# Inline args (username, then password)
//	tls {
//	    dns servercow <username> <password>
//	}
//
//	# Block form
//	tls {
//	    dns servercow {
//	        username <username>
//	        password <password>
//	    }
//	}
//
// Credentials may use Caddy placeholder syntax, e.g. {env.SERVERCOW_USERNAME}.
package servercow

import (
	"fmt"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	servercowdns "github.com/libdns/servercow"
)

// Provider wraps the libdns Servercow provider as a Caddy module.
type Provider struct{ *servercowdns.Provider }

func init() {
	caddy.RegisterModule(Provider{})
}

// CaddyModule returns the Caddy module information.
func (Provider) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "dns.providers.servercow",
		New: func() caddy.Module { return &Provider{new(servercowdns.Provider)} },
	}
}

// Provision sets up the module. It resolves any Caddy placeholder expressions
// (e.g. environment variable references) in the credentials.
// Implements caddy.Provisioner.
func (p *Provider) Provision(ctx caddy.Context) error {
	repl := caddy.NewReplacer()
	p.Provider.Username = repl.ReplaceAll(p.Provider.Username, "")
	p.Provider.Password = repl.ReplaceAll(p.Provider.Password, "")
	if p.Provider.Username == "" {
		return fmt.Errorf("servercow: username is required")
	}
	if p.Provider.Password == "" {
		return fmt.Errorf("servercow: password is required")
	}
	return nil
}

// UnmarshalCaddyfile sets up the DNS provider from Caddyfile tokens.
//
// Syntax:
//
//	servercow [<username> <password>] {
//	    username <username>
//	    password <password>
//	}
func (p *Provider) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		// Optional inline args: servercow <username> <password>
		if d.NextArg() {
			p.Provider.Username = d.Val()
		}
		if d.NextArg() {
			p.Provider.Password = d.Val()
		}
		if d.NextArg() {
			return d.ArgErr()
		}

		// Block options
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			switch d.Val() {
			case "username":
				if p.Provider.Username != "" {
					return d.Err("username already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.Username = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			case "password":
				if p.Provider.Password != "" {
					return d.Err("password already set")
				}
				if !d.NextArg() {
					return d.ArgErr()
				}
				p.Provider.Password = d.Val()
				if d.NextArg() {
					return d.ArgErr()
				}
			default:
				return d.Errf("unrecognized subdirective '%s'", d.Val())
			}
		}
	}
	return nil
}

// Interface guards — compile-time checks.
var (
	_ caddyfile.Unmarshaler = (*Provider)(nil)
	_ caddy.Provisioner     = (*Provider)(nil)
)
