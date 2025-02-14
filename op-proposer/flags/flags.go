package flags

import (
	"fmt"
	"time"

	"github.com/urfave/cli/v2"

	opservice "github.com/ethereum-optimism/optimism/op-service"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	opmetrics "github.com/ethereum-optimism/optimism/op-service/metrics"
	oppprof "github.com/ethereum-optimism/optimism/op-service/pprof"
	oprpc "github.com/ethereum-optimism/optimism/op-service/rpc"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"
)

const EnvVarPrefix = "OP_PROPOSER"

func prefixEnvVars(name string) []string {
	return opservice.PrefixEnvVar(EnvVarPrefix, name)
}

var (
	// Required Flags
	L1EthRpcFlag = &cli.StringFlag{
		Name:    "l1-eth-rpc",
		Usage:   "HTTP provider URL for L1",
		EnvVars: prefixEnvVars("L1_ETH_RPC"),
	}
	RollupRpcFlag = &cli.StringFlag{
		Name:    "rollup-rpc",
		Usage:   "HTTP provider URL for the rollup node. A comma-separated list enables the active rollup provider.",
		EnvVars: prefixEnvVars("ROLLUP_RPC"),
	}

	// Optional flags
	L2OOAddressFlag = &cli.StringFlag{
		Name:    "l2oo-address",
		Usage:   "Address of the L2OutputOracle contract",
		EnvVars: prefixEnvVars("L2OO_ADDRESS"),
	}
	PollIntervalFlag = &cli.DurationFlag{
		Name:    "poll-interval",
		Usage:   "How frequently to poll L2 for new blocks",
		Value:   6 * time.Second,
		EnvVars: prefixEnvVars("POLL_INTERVAL"),
	}
	AllowNonFinalizedFlag = &cli.BoolFlag{
		Name:    "allow-non-finalized",
		Usage:   "Allow the proposer to submit proposals for L2 blocks derived from non-finalized L1 blocks.",
		EnvVars: prefixEnvVars("ALLOW_NON_FINALIZED"),
	}
	DisputeGameFactoryAddressFlag = &cli.StringFlag{
		Name:    "dgf-address",
		Usage:   "Address of the DisputeGameFactory contract",
		EnvVars: prefixEnvVars("DGF_ADDRESS"),
		Hidden:  true,
	}
	ProposalIntervalFlag = &cli.DurationFlag{
		Name:    "proposal-interval",
		Usage:   "Interval between submitting L2 output proposals when the DGFAddress is set",
		EnvVars: prefixEnvVars("PROPOSAL_INTERVAL"),
		Hidden:  true,
	}
	DisputeGameTypeFlag = &cli.UintFlag{
		Name:    "dg-type",
		Usage:   "Dispute game type to create via the configured DisputeGameFactory",
		Value:   0,
		EnvVars: prefixEnvVars("DG_TYPE"),
		Hidden:  true,
	}
	// Legacy Flags
	L2OutputHDPathFlag = txmgr.L2OutputHDPathFlag
)

var requiredFlags = []cli.Flag{
	L1EthRpcFlag,
	RollupRpcFlag,
}

var optionalFlags = []cli.Flag{
	L2OOAddressFlag,
	PollIntervalFlag,
	AllowNonFinalizedFlag,
	L2OutputHDPathFlag,
	DisputeGameFactoryAddressFlag,
	ProposalIntervalFlag,
	DisputeGameTypeFlag,
}

func init() {
	optionalFlags = append(optionalFlags, oprpc.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, oplog.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, opmetrics.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, oppprof.CLIFlags(EnvVarPrefix)...)
	optionalFlags = append(optionalFlags, txmgr.CLIFlags(EnvVarPrefix)...)

	Flags = append(requiredFlags, optionalFlags...)
}

// Flags contains the list of configuration options available to the binary.
var Flags []cli.Flag

func CheckRequired(ctx *cli.Context) error {
	for _, f := range requiredFlags {
		if !ctx.IsSet(f.Names()[0]) {
			return fmt.Errorf("flag %s is required", f.Names()[0])
		}
	}
	return nil
}
