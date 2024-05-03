package cmd

import (
	"context"
	"fmt"
	"github.com/setavenger/blindbit-cli/lib"
	"github.com/setavenger/blindbit-cli/lib/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"log"
)

// balanceCmd represents the balance command
var (
	listUTXOs bool

	showUnconfirmed      bool = false
	showUnspent          bool = false
	showSpent            bool = false
	showSpentUnConfirmed bool = false

	states []pb.UTXOState

	balanceCmd = &cobra.Command{
		Use:   "balance",
		Short: "shows the balance of the wallet",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			client, conn := lib.NewClient(socketPath)
			defer func(conn *grpc.ClientConn) {
				err := conn.Close()
				if err != nil {
					panic(err)
				}
			}(conn)

			if showUnconfirmed {
				states = append(states, pb.UTXOState_UNCONFIRMED)
			}
			if showUnspent {
				states = append(states, pb.UTXOState_UNSPENT)
			}
			if showSpent {
				states = append(states, pb.UTXOState_SPENT)
			}
			if showSpentUnConfirmed {
				states = append(states, pb.UTXOState_SPENT_UNCONFIRMED)
			}
			if len(states) == 0 {
				states = append(states, pb.UTXOState_UNSPENT)
			}

			utxos, err := client.ListUTXOs(context.Background(), &pb.Empty{})
			if err != nil {
				log.Fatalf("ListUTXOs failed: %v\n", err)
			}

			var balance uint64

			for _, state := range states {
				for _, utxo := range utxos.Utxos {
					if utxo.UtxoState == state {
						if listUTXOs {
							label := ""
							if utxo.Label != nil {
								label = *utxo.Label
							}
							fmt.Printf("%x:%d - %s - %s -%s\n", utxo.Txid, utxo.Vout, lib.ConvertIntToThousandString(int(utxo.Amount)), label, utxo.UtxoState)
						} else {
							// todo if another case opens up change this
							balance += utxo.Amount
						}
					}
				}
			}

			if listUTXOs {
				return
			}

			fmt.Printf("Balance is %s\n", lib.ConvertIntToThousandString(int(balance)))
		},
	}
)

func init() {
	RootCmd.AddCommand(balanceCmd)

	balanceCmd.PersistentFlags().BoolVar(&listUTXOs, "list", false, "list utxos instead showing the balance")

	balanceCmd.PersistentFlags().BoolVar(&showUnconfirmed, "unconfirmed", false, "add unconfirmed utxos to the filter")
	balanceCmd.PersistentFlags().BoolVar(&showUnspent, "unspent", false, "add unspent utxos to the filter")
	balanceCmd.PersistentFlags().BoolVar(&showSpent, "spent", false, "add spent utxos to the filter")
	balanceCmd.PersistentFlags().BoolVar(&showSpentUnConfirmed, "spentunconf", false, "add spent utxos whose spending transaction is not confirmed to the filter")

}

// usageUTXOStatesFlag
// Deprecated:: do not use
var usageUTXOStatesFlag = `	
The following states can be set. The state has to be given as an integer not a
	UTXOState_UNKNOWN           = 0
	UTXOState_UNCONFIRMED       = 1
	UTXOState_UNSPENT           = 2
	UTXOState_SPENT             = 3
	UTXOState_SPENT_UNCONFIRMED = 4
`
