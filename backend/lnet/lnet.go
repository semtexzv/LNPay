package lnet

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"io/ioutil"
	"lnpay/db"
	"os"
	"strings"
)

var LnClient lnrpc.LightningClient

func init() {
	godotenv.Load()

	opts := []grpc.DialOption{}

	creds, err := credentials.NewClientTLSFromFile(os.Getenv("LND_TLS_CERT"), "")
	if err != nil {
		panic("Could not get lnd cert")
	}

	lndAddr := os.Getenv("LND_ADDRESS")
	println("Loaded tls cert, connecting to lnd on: ", lndAddr)

	opts = append(opts,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(creds),
	)

	conn, err := grpc.DialContext(context.Background(), lndAddr, opts...)
	if err != nil {
		panic(err)
	}

	walletClient := lnrpc.NewWalletUnlockerClient(conn)

	mac, err := InitLndConn(walletClient)
	if err != nil {
		panic(err)
	}
	conn.Close()

	opts = append(opts, grpc.WithPerRPCCredentials(macaroons.NewMacaroonCredential(mac)))
	conn, err = grpc.DialContext(context.Background(), lndAddr, opts...)

	LnClient = lnrpc.NewLightningClient(conn)
	_, err = LnClient.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	if err != nil {
		panic(err)
	}
	println("Opened connection to lnd")
}

func InitLndWallet(cl lnrpc.WalletUnlockerClient) (*macaroon.Macaroon, error) {
	seed, err := cl.GenSeed(context.Background(), &lnrpc.GenSeedRequest{
		AezeedPassphrase: []byte("password"),
	})
	if err != nil {
		return nil, err
	}

	wallet, err := cl.InitWallet(context.Background(), &lnrpc.InitWalletRequest{
		WalletPassword:     []byte("password"),
		AezeedPassphrase:   []byte("password"),
		StatelessInit:      true,
		CipherSeedMnemonic: seed.CipherSeedMnemonic,
	})
	if err != nil {
		return nil, err
	}

	err = ioutil.WriteFile(os.Getenv("LND_ADMIN_MACAROON"), wallet.AdminMacaroon, 0644)
	if err != nil {
		return nil, err
	}
	mac := &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(wallet.AdminMacaroon)
	if err != nil {
		return nil, err
	}

	return mac, nil
}

func LndAddresses(cl lnrpc.LightningClient) (string, error) {
	res, err := cl.NewAddress(context.Background(), &lnrpc.NewAddressRequest{
		Type: lnrpc.AddressType_WITNESS_PUBKEY_HASH,
	})

	if err != nil {
		return "", err
	}
	return res.Address, nil
}

func UnlockLndWallet(cl lnrpc.WalletUnlockerClient) (*macaroon.Macaroon, error) {
	mdata, err := ioutil.ReadFile(os.Getenv("LND_ADMIN_MACAROON"))
	if err != nil {
		return nil, err
	}
	mac := &macaroon.Macaroon{}
	err = mac.UnmarshalBinary(mdata)
	if err != nil {
		return nil, err
	}

	_, err = cl.UnlockWallet(context.Background(), &lnrpc.UnlockWalletRequest{
		StatelessInit:  true,
		WalletPassword: []byte("password"),
	})
	if err != nil && strings.Contains(err.Error(), "unknown service lnrpc.WalletUnlocker") {
		return mac, nil
	} else if err != nil {
		return nil, err
	}
	return mac, nil
}

func InitLndConn(cl lnrpc.WalletUnlockerClient) (*macaroon.Macaroon, error) {
	if _, err := os.Stat(os.Getenv("LND_ADMIN_MACAROON")); os.IsNotExist(err) {
		return InitLndWallet(cl)
	} else {
		return UnlockLndWallet(cl)
	}
}

func ConnectToNodes(cl lnrpc.LightningClient) {
	cl.ConnectPeer(context.Background(), &lnrpc.ConnectPeerRequest{

	})
}
func SendPayment(invoice db.Invoice) (string, error) {
	info, err := LnClient.GetInfo(context.Background(), &lnrpc.GetInfoRequest{})
	println(info)

	payment, err := LnClient.SendPaymentSync(context.Background(), &lnrpc.SendRequest{
		PaymentRequest: invoice.Invoice,
	})
	if err != nil {
		return "", err
	}
	if payment.PaymentError != "" {
		return "", errors.New(payment.PaymentError)
	}

	pi, err := lntypes.MakePreimage(payment.PaymentPreimage)
	// Everything is complete, we can't revert now
	if err != nil {
		return "", nil
	}
	return pi.String(), nil
}
