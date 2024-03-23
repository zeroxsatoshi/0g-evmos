package cli

import (
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/coniks-sys/coniks-go/crypto/vrf"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		NewProposeCmd(),
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

func NewProposeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose",
		Short: "Propose to create a new committee",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			valAddr, err := sdk.ValAddressFromHex(hex.EncodeToString(clientCtx.GetFromAddress().Bytes()))
			if err != nil {
				return err
			}

			startingHeight, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			votingStartingHeight, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			votingEndHeight, err := strconv.ParseUint(args[2], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgPropose{
				Proposer:          valAddr.String(),
				StartingHeight:    startingHeight,
				VotingStartHeight: votingStartingHeight,
				VotingEndHeight:   votingEndHeight,
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewVoteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote proposal-id",
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

			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			msg := &types.MsgVote{
				ProposalID: proposalID,
				Voter:      valAddr.String(),
				Ballots:    nil,
				// TODO
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
