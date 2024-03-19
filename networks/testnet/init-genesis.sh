#!/bin/bash

ROOT_DIR=${ROOT_DIR:-testnet}
CHAIN_ID=${CHAIN_ID:-zgtendermint_9000-1}

# Usage: init-genesis.sh IP1,IP2,IP3 KEYRING_PASSWORD
OS_NAME=`uname -o`
USAGE="Usage: ${BASH_SOURCE[0]} IP1,IP2,IP3"
if [[ "$OS_NAME" = "GNU/Linux" ]]; then
	USAGE="$USAGE KEYRING_PASSWORD"
fi

if [[ $# -eq 0 ]]; then
	echo "IP list not specified"
	echo $USAGE
	exit 1
fi

if [[ "$OS_NAME" = "GNU/Linux" ]]; then
	if [[ $# -eq 1 ]]; then
		echo "Keyring password not specified"
		echo $USAGE
		exit 1
	fi

	PASSWORD=$2
fi

evmosd version 2>/dev/null || export PATH=$PATH:$(go env GOPATH)/bin

set -e

IFS=","; declare -a IPS=($1); unset IFS

NUM_NODES=${#IPS[@]}
BALANCE=$((100000000/$NUM_NODES))evmos
STAKING=$((50000000/$NUM_NODES))evmos

# Init configs
for ((i=0; i<$NUM_NODES; i++)) do
	HOMEDIR="$ROOT_DIR"/node$i

	# Init
	evmosd init "node$i" --home "$HOMEDIR" --chain-id "$CHAIN_ID" >/dev/null 2>&1
	
	# Change parameter token denominations to aevmos
	GENESIS="$HOMEDIR"/config/genesis.json
	TMP_GENESIS="$HOMEDIR"/config/tmp_genesis.json
	cat "$GENESIS" | jq '.app_state["staking"]["params"]["bond_denom"]="aevmos"' >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
	cat "$GENESIS" | jq '.app_state["gov"]["params"]["min_deposit"][0]["denom"]="aevmos"' >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

	# Change app.toml
	APP_TOML="$HOMEDIR"/config/app.toml
	sed -i 's/minimum-gas-prices = "0aevmos"/minimum-gas-prices = "1000000000aevmos"/' "$APP_TOML"
	sed -i '/\[json-rpc\]/,/^\[/ s/enable = false/enable = true/' "$APP_TOML"
	sed -i '/\[json-rpc\]/,/^\[/ s/address = "127.0.0.1:8545"/address = "0.0.0.0:8545"/' "$APP_TOML"
done

# Update seeds in config.toml
SEEDS=""
for ((i=0; i<$NUM_NODES; i++)) do
	if [[ $i -gt 0 ]]; then SEEDS=$SEEDS,; fi
	NODE_ID=`evmosd tendermint show-node-id --home $ROOT_DIR/node$i`
	SEEDS=$SEEDS$NODE_ID@${IPS[$i]}:26656
done

for ((i=0; i<$NUM_NODES; i++)) do
	sed -i "/seeds = /c\seeds = \"$SEEDS\"" "$ROOT_DIR"/node$i/config/config.toml
done

# Prepare validators
#
# Note, keyring backend `file` works bad on Windows, and `add-genesis-account`
# do not supports --keyring-dir flag. As a result, we use keyring backend `os`,
# which is the default value.
#
# Where key stored:
# - Windows: Windows credentials management.
# - Linux: under `--home` specified folder.
if [[ "$OS_NAME" = "Msys" ]]; then
	for ((i=0; i<$NUM_NODES; i++)) do
		VALIDATOR="0gchain_9000_validator_$i"
		set +e
		ret=`evmosd keys list --keyring-backend os -n | grep $VALIDATOR`
		set -e
		if [[ "$ret" = "" ]]; then
			echo "Create validator key: $VALIDATOR"
			evmosd keys add $VALIDATOR --keyring-backend os
		fi
	done
elif [[ "$OS_NAME" = "GNU/Linux" ]]; then
	# Create N validators for node0
	for ((i=0; i<$NUM_NODES; i++)) do
		yes $PASSWORD | evmosd keys add "0gchain_9000_validator_$i" --keyring-backend os --home "$ROOT_DIR"/node0
	done

	# Copy validators to other nodes
	for ((i=1; i<$NUM_NODES; i++)) do
		cp "$ROOT_DIR"/node0/keyhash "$ROOT_DIR"/node$i
		cp "$ROOT_DIR"/node0/*.address "$ROOT_DIR"/node$i
		cp "$ROOT_DIR"/node0/*.info "$ROOT_DIR"/node$i
	done
else
	echo -e "\n\nOS: $OS_NAME"
	echo "Unsupported OS to generate keys for validators!!!"
	exit 1
fi

# Add all validators in genesis
for ((i=0; i<$NUM_NODES; i++)) do
	for ((j=0; j<$NUM_NODES; j++)) do
		if [[ "$OS_NAME" = "GNU/Linux" ]]; then
			yes $PASSWORD | evmosd add-genesis-account "0gchain_9000_validator_$i" $BALANCE --home "$ROOT_DIR/node$j"
		else
			evmosd add-genesis-account "0gchain_9000_validator_$i" $BALANCE --home "$ROOT_DIR/node$j"
		fi 
	done
done

# Prepare genesis txs
mkdir -p "$ROOT_DIR"/gentxs
for ((i=0; i<$NUM_NODES; i++)) do
	if [[ "$OS_NAME" = "GNU/Linux" ]]; then
		yes $PASSWORD | evmosd gentx "0gchain_9000_validator_$i" $STAKING --home "$ROOT_DIR/node$i" --output-document "$ROOT_DIR/gentxs/node$i.json"
	else
		evmosd gentx "0gchain_9000_validator_$i" $STAKING --home "$ROOT_DIR/node$i" --output-document "$ROOT_DIR/gentxs/node$i.json"
	fi 
done

# Create genesis at node0 and copy to other nodes
evmosd collect-gentxs --home "$ROOT_DIR/node0" --gentx-dir "$ROOT_DIR/gentxs" >/dev/null 2>&1
sed -i '/persistent_peers = /c\persistent_peers = ""' "$ROOT_DIR"/node0/config/config.toml
evmosd validate-genesis --home "$ROOT_DIR/node0"
for ((i=1; i<$NUM_NODES; i++)) do
	cp "$ROOT_DIR"/node0/config/genesis.json "$ROOT_DIR"/node$i/config/genesis.json
done

# For linux, backup keys for all validators
if [[ "$OS_NAME" = "GNU/Linux" ]]; then
	mkdir -p "$ROOT_DIR"/keyring-os

	cp "$ROOT_DIR"/node0/keyhash "$ROOT_DIR"/keyring-os
	cp "$ROOT_DIR"/node0/*.address "$ROOT_DIR"/keyring-os
	cp "$ROOT_DIR"/node0/*.info "$ROOT_DIR"/keyring-os

	for ((i=0; i<$NUM_NODES; i++)) do
		rm -f "$ROOT_DIR"/node$i/keyhash "$ROOT_DIR"/node$i/*.address "$ROOT_DIR"/node$i/*.info
	done
fi

echo -e "\n\nSucceeded to init genesis!\n"
