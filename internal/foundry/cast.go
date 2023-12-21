// Copyright (c) Gabriel de Quadros Ligneul
// SPDX-License-Identifier: Apache-2.0 (see LICENSE)

package foundry

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var SenderAddress = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")

// Send an input using the devnet addresses.
// This should be used only for testing.
func AddInput(ctx context.Context, payload []byte) error {
	if len(payload) == 0 {
		return fmt.Errorf("cast: cannot send empty payload")
	}
	input := hexutil.Encode(payload)
	cmd := exec.CommandContext(ctx,
		"cast", "send",
		"--mnemonic", TestMnemonic,
		"--rpc-url", fmt.Sprintf("http://127.0.0.1:%v", AnvilDefaultPort),
		InputBoxAddress,                    // TO
		"addInput(address,bytes)(bytes32)", // SIG
		ApplicationAddress, input,          // ARGS
	)
	log.Printf(`calling: "%v"`, strings.Join(cmd.Args, `" "`))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("cast: %w: %v", err, string(output))
	}
	return nil
}
