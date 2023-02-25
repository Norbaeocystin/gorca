// Code Generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package whirlpool

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// SetFeeAuthority is the `setFeeAuthority` instruction.
type SetFeeAuthority struct {

	// [0] = [WRITE] whirlpoolsConfig
	//
	// [1] = [SIGNER] feeAuthority
	//
	// [2] = [] newFeeAuthority
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewSetFeeAuthorityInstructionBuilder creates a new `SetFeeAuthority` instruction builder.
func NewSetFeeAuthorityInstructionBuilder() *SetFeeAuthority {
	nd := &SetFeeAuthority{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 3),
	}
	return nd
}

// SetWhirlpoolsConfigAccount sets the "whirlpoolsConfig" account.
func (inst *SetFeeAuthority) SetWhirlpoolsConfigAccount(whirlpoolsConfig ag_solanago.PublicKey) *SetFeeAuthority {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(whirlpoolsConfig).WRITE()
	return inst
}

// GetWhirlpoolsConfigAccount gets the "whirlpoolsConfig" account.
func (inst *SetFeeAuthority) GetWhirlpoolsConfigAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetFeeAuthorityAccount sets the "feeAuthority" account.
func (inst *SetFeeAuthority) SetFeeAuthorityAccount(feeAuthority ag_solanago.PublicKey) *SetFeeAuthority {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(feeAuthority).SIGNER()
	return inst
}

// GetFeeAuthorityAccount gets the "feeAuthority" account.
func (inst *SetFeeAuthority) GetFeeAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetNewFeeAuthorityAccount sets the "newFeeAuthority" account.
func (inst *SetFeeAuthority) SetNewFeeAuthorityAccount(newFeeAuthority ag_solanago.PublicKey) *SetFeeAuthority {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(newFeeAuthority)
	return inst
}

// GetNewFeeAuthorityAccount gets the "newFeeAuthority" account.
func (inst *SetFeeAuthority) GetNewFeeAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

func (inst SetFeeAuthority) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_SetFeeAuthority,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst SetFeeAuthority) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *SetFeeAuthority) Validate() error {
	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.WhirlpoolsConfig is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.FeeAuthority is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.NewFeeAuthority is not set")
		}
	}
	return nil
}

func (inst *SetFeeAuthority) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("SetFeeAuthority")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=0]").ParentFunc(func(paramsBranch ag_treeout.Branches) {})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=3]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("whirlpoolsConfig", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("    feeAuthority", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta(" newFeeAuthority", inst.AccountMetaSlice.Get(2)))
					})
				})
		})
}

func (obj SetFeeAuthority) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	return nil
}
func (obj *SetFeeAuthority) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	return nil
}

// NewSetFeeAuthorityInstruction declares a new SetFeeAuthority instruction with the provided parameters and accounts.
func NewSetFeeAuthorityInstruction(
	// Accounts:
	whirlpoolsConfig ag_solanago.PublicKey,
	feeAuthority ag_solanago.PublicKey,
	newFeeAuthority ag_solanago.PublicKey) *SetFeeAuthority {
	return NewSetFeeAuthorityInstructionBuilder().
		SetWhirlpoolsConfigAccount(whirlpoolsConfig).
		SetFeeAuthorityAccount(feeAuthority).
		SetNewFeeAuthorityAccount(newFeeAuthority)
}