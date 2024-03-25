package cli

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	sdkmath "cosmossdk.io/math"
	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/evmos/evmos/v16/x/committee/v1/types"
	"github.com/spf13/cobra"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(
		NewRegisterCmd(),
		NewVoteCmd(),
	)
	return cmd
}

func NewRegisterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register a voter",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromHex(hex.EncodeToString(clientCtx.GetFromAddress().Bytes()))
			if err != nil {
				return err
			}

			sk, err := vrf.GenerateKey(nil)
			if err != nil {
				return err
			}
			pk, _ := sk.Public()
			// TODO: save private key

			msg := &types.MsgRegister{
				Voter: valAddr.String(),
				Key:   pk,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote committee-id",
		Short: "Vote on a proposal",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromHex(hex.EncodeToString(clientCtx.GetFromAddress().Bytes()))
			if err != nil {
				return err
			}

			committeeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			votingStartHeight := types.DefaultVotingStartHeight + (committeeID-1)*types.DefaultVotingPeriod

			rsp, err := stakingtypes.NewQueryClient(clientCtx).HistoricalInfo(cmd.Context(), &stakingtypes.QueryHistoricalInfoRequest{Height: int64(votingStartHeight)})
			if err != nil {
				return err
			}

			// TODO: DO NOT generate a new pair of keys
			sk, err := vrf.GenerateKey(nil)
			if err != nil {
				return err
			}

			var tokens sdkmath.Int
			for _, val := range rsp.Hist.Valset {
				if val.GetOperator().Equals(valAddr) {
					tokens = val.GetTokens()
				}
			}

			// 1_000 0AGI token / vote
			numBallots := tokens.Quo(sdk.NewInt(1_000_000_000_000_000_000)).Quo(sdk.NewInt(1_000)).Uint64()
			ballots := make([]*types.Ballot, numBallots)
			for i := range ballots {
				ballotID := uint64(i)
				ballots[i] = &types.Ballot{
					ID:      ballotID,
					Content: sk.Compute(bytes.Join([][]byte{rsp.Hist.Header.LastCommitHash, types.Uint64ToBytes(ballotID)}, nil)),
				}
			}

			msg := &types.MsgVote{
				CommitteeID: committeeID,
				Voter:       valAddr.String(),
				Ballots:     ballots,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
