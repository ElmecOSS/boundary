package credentialstorescmd

import (
	"github.com/hashicorp/boundary/api/credentialstores"
	"github.com/hashicorp/boundary/internal/cmd/base"
)

func init() {
	extraVaultFlagsFunc = extraVaultFlagsFuncImpl
	extraVaultActionsFlagsMapFunc = extraVaultActionsFlagsMapFuncImpl
	extraVaultFlagsHandlingFunc = extraVaultFlagHandlingFuncImpl
}

const (
	addressFlagName              = "vault-address"
	namespaceFlagName            = "vault-namespace"
	vaultCaCertFlagName          = "vault-ca-cert"
	tlsServerNameFlagName        = "vault-tls-server-name"
	tlsSkipVerifyFlagName        = "vault-tls-skip-verify"
	vaultTokenFlagName           = "vault-token"
	clientCertificateFlagName    = "vault-client-certificate"
	clientCertificateKeyFlagName = "vault-client-certificate-key"
)

type extraVaultCmdVars struct {
	flagAddress       string
	flagNamespace     string
	flagCaCert        string
	flagVaultToken    string
	flagClientCert    string
	flagClientCertKey string
	flagTlsServerName string
	flagTlsSkipVerify bool
}

func extraVaultActionsFlagsMapFuncImpl() map[string][]string {
	flags := map[string][]string{
		"create": {
			addressFlagName,
			namespaceFlagName,
			vaultCaCertFlagName,
			tlsServerNameFlagName,
			tlsSkipVerifyFlagName,
			vaultTokenFlagName,
			clientCertificateFlagName,
			clientCertificateKeyFlagName,
		},
	}
	flags["update"] = flags["create"]
	return flags
}

func extraVaultFlagsFuncImpl(c *VaultCommand, set *base.FlagSets, _ *base.FlagSet) {
	f := set.NewFlagSet("Vault Credential Store Options")

	for _, name := range flagsVaultMap[c.Func] {
		switch name {
		case addressFlagName:
			f.StringVar(&base.StringVar{
				Name:   addressFlagName,
				Target: &c.flagAddress,
				Usage:  "The vault address.",
			})
		case namespaceFlagName:
			f.StringVar(&base.StringVar{
				Name:   namespaceFlagName,
				Target: &c.flagNamespace,
				Usage:  "The vault namespace the store should use.",
			})
		case vaultCaCertFlagName:
			f.StringVar(&base.StringVar{
				Name:   vaultCaCertFlagName,
				Target: &c.flagCaCert,
				Usage:  "The CA Cert to use when connecting with vault.",
			})
		case tlsServerNameFlagName:
			f.StringVar(&base.StringVar{
				Name:   tlsServerNameFlagName,
				Target: &c.flagTlsServerName,
				Usage:  `The TLS server name used when connecting to vault.`,
			})
		case tlsSkipVerifyFlagName:
			f.BoolVar(&base.BoolVar{
				Name:   tlsSkipVerifyFlagName,
				Target: &c.flagTlsSkipVerify,
				Usage:  "Whether to skip tls verification.",
			})
		case vaultTokenFlagName:
			f.StringVar(&base.StringVar{
				Name:   vaultTokenFlagName,
				Target: &c.flagVaultToken,
				Usage:  "The vault token to use when boundary connects to vault for this store.",
			})
		case clientCertificateFlagName:
			f.StringVar(&base.StringVar{
				Name:   clientCertificateFlagName,
				Target: &c.flagClientCert,
				Usage:  "The client certificate to use when boundary connects to vault for this store.",
			})
		case clientCertificateKeyFlagName:
			f.StringVar(&base.StringVar{
				Name:   clientCertificateKeyFlagName,
				Target: &c.flagClientCertKey,
				Usage:  `The client certificate's private key to use when boundary connects to vault for this store.`,
			})
		}
	}
}

func extraVaultFlagHandlingFuncImpl(c *VaultCommand, f *base.FlagSets, opts *[]credentialstores.Option) bool {
	switch c.flagAddress {
	case "":
	case "null":
		*opts = append(*opts, credentialstores.DefaultVaultCredentialStoreAddress())
	default:
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreAddress(c.flagAddress))
	}
	switch c.flagNamespace {
	case "":
	case "null":
		*opts = append(*opts, credentialstores.DefaultVaultCredentialStoreNamespace())
	default:
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreNamespace(c.flagNamespace))
	}
	switch c.flagVaultToken {
	case "":
	case "null":
		*opts = append(*opts, credentialstores.DefaultVaultCredentialStoreToken())
	default:
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreToken(c.flagVaultToken))
	}
	switch c.flagCaCert {
	case "":
	case "null":
		*opts = append(*opts, credentialstores.DefaultVaultCredentialStoreCaCert())
	default:
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreCaCert(c.flagCaCert))
	}
	switch c.flagClientCert {
	case "":
	case "null":
		*opts = append(*opts, credentialstores.DefaultVaultCredentialStoreClientCertificate())
	default:
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreClientCertificate(c.flagClientCert))
	}
	switch c.flagClientCertKey {
	case "":
	case "null":
		*opts = append(*opts, credentialstores.DefaultVaultCredentialStoreClientCertificateKey())
	default:
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreClientCertificateKey(c.flagClientCertKey))
	}
	if c.flagTlsSkipVerify {
		*opts = append(*opts, credentialstores.WithVaultCredentialStoreTlsSkipVerify(c.flagTlsSkipVerify))
	}

	return true
}

func (c *VaultCommand) extraVaultHelpFunc(helpMap map[string]func() string) string {
	var helpStr string
	switch c.Func {
	case "create":
		helpStr = base.WrapForHelpText([]string{
			"Usage: boundary credential-stores create vault [options] [args]",
			"",
			"  Create a vault-type credential store. Example:",
			"",
			`    $ boundary credential-stores create vault -vault-address "http://localhost:55002" -vault-token"s.s0m3t0k3n"`,
			"",
			"",
		})

	case "update":
		helpStr = base.WrapForHelpText([]string{
			"Usage: boundary credential-stores update vault [options] [args]",
			"",
			"  Update a vault-type credential store given its ID. Example:",
			"",
			`    $ boundary credential-stores update vault -id csvlt_1234567890 -name devops -description "For DevOps usage"`,
			"",
			"",
		})
	}
	return helpStr + c.Flags().Help()
}
