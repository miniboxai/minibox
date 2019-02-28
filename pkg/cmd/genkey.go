package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"path"
	"time"

	cobra "github.com/spf13/cobra"
)

var (
	secretDir  string
	hostnames  []string
	validFrom  string
	validFor   time.Duration
	isCA       bool
	rsaBits    int
	ecdsaCurve string
)

// genkeyCmd only in server
var genkeyCmd = &cobra.Command{
	Use:   "genkey",
	Short: "Generate a pair tls secret keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(hostnames) == 0 {
			log.Fatalf("Missing required --host parameter")
		}
		var priv interface{}
		var err error
		switch ecdsaCurve {
		case "":
			priv, err = rsa.GenerateKey(rand.Reader, rsaBits)
		case "P224":
			priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
		case "P256":
			priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		case "P384":
			priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
		case "P521":
			priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		default:
			fmt.Fprintf(os.Stderr, "Unrecognized elliptic curve: %q", ecdsaCurve)
			os.Exit(1)
		}
		if err != nil {
			log.Fatalf("failed to generate private key: %s", err)
		}
		var notBefore time.Time
		if len(validFrom) == 0 {
			notBefore = time.Now()
		} else {
			notBefore, err = time.Parse("Jan 2 15:04:05 2006", validFrom)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse creation date: %s\n", err)
				os.Exit(1)
			}
		}

		notAfter := notBefore.Add(validFor)

		serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
		serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
		if err != nil {
			log.Fatalf("failed to generate serial number: %s", err)
		}

		template := x509.Certificate{
			SerialNumber: serialNumber,
			Subject: pkix.Name{
				Organization: []string{"Minibox Inc"},
			},
			NotBefore: notBefore,
			NotAfter:  notAfter,

			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
		}

		// hosts := strings.Split(hostnames, ",")
		for _, h := range hostnames {
			if ip := net.ParseIP(h); ip != nil {
				template.IPAddresses = append(template.IPAddresses, ip)
			} else {
				template.DNSNames = append(template.DNSNames, h)
			}
		}

		if isCA {
			template.IsCA = true
			template.KeyUsage |= x509.KeyUsageCertSign
		}

		derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
		if err != nil {
			log.Fatalf("Failed to create certificate: %s", err)
		}

		os.MkdirAll(secretDir, 0755)

		certOut, err := os.Create(path.Join(secretDir, "server.pem"))
		if err != nil {
			log.Fatalf("failed to open server.pem for writing: %s", err)
		}
		pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
		certOut.Close()
		log.Print("written server.pem\n")

		keyOut, err := os.OpenFile(path.Join(secretDir, "server.key"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			log.Print("failed to open server.key for writing:", err)
			return
		}
		pem.Encode(keyOut, pemBlockForKey(priv))
		keyOut.Close()
		log.Print("written server.key\n")
	},
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func init() {
	genkeyCmd.PersistentFlags().StringVar(&secretDir, "secret-dir", "./certs", "secret keys output dir file (default is ./certs)")
	genkeyCmd.PersistentFlags().StringArrayVar(&hostnames, "hostname", []string{"localhost:8080"}, "Comma-separated hostnames and IPs to generate a certificate for")
	genkeyCmd.PersistentFlags().StringVar(&validFrom, "start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	genkeyCmd.PersistentFlags().DurationVar(&validFor, "duration", 365*24*time.Hour, "Duration that certificate is valid for")
	genkeyCmd.PersistentFlags().BoolVar(&isCA, "ca", false, "whether this cert should be its own Certificate Authority")
	genkeyCmd.PersistentFlags().IntVar(&rsaBits, "rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	genkeyCmd.PersistentFlags().StringVar(&ecdsaCurve, "ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")

}
