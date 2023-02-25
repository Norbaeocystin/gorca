// Code Generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package whirlpool

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// InitializeTickArray is the `initializeTickArray` instruction.
type InitializeTickArray struct {
	StartTickIndex *int32

	// [0] = [] whirlpool
	//
	// [1] = [WRITE, SIGNER] funder
	//
	// [2] = [WRITE] tickArray
	//
	// [3] = [] systemProgram
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewInitializeTickArrayInstructionBuilder creates a new `InitializeTickArray` instruction builder.
func NewInitializeTickArrayInstructionBuilder() *InitializeTickArray {
	nd := &InitializeTickArray{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 4),
	}
	return nd
}

// SetStartTickIndex sets the "startTickIndex" parameter.
func (inst *InitializeTickArray) SetStartTickIndex(startTickIndex int32) *InitializeTickArray {
	inst.StartTickIndex = &startTickIndex
	return inst
}

// SetWhirlpoolAccount sets the "whirlpool" account.
func (inst *InitializeTickArray) SetWhirlpoolAccount(whirlpool ag_solanago.PublicKey) *InitializeTickArray {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(whirlpool)
	return inst
}

// GetWhirlpoolAccount gets the "whirlpool" account.
func (inst *InitializeTickArray) GetWhirlpoolAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetFunderAccount sets the "funder" account.
func (inst *InitializeTickArray) SetFunderAccount(funder ag_solanago.PublicKey) *InitializeTickArray {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(funder).WRITE().SIGNER()
	return inst
}

// GetFunderAccount gets the "funder" account.
func (inst *InitializeTickArray) GetFunderAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetTickArrayAccount sets the "tickArray" account.
func (inst *InitializeTickArray) SetTickArrayAccount(tickArray ag_solanago.PublicKey) *InitializeTickArray {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(tickArray).WRITE()
	return inst
}

// GetTickArrayAccount gets the "tickArray" account.
func (inst *InitializeTickArray) GetTickArrayAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetSystemProgramAccount sets the "systemProgram" account.
func (inst *InitializeTickArray) SetSystemProgramAccount(systemProgram ag_solanago.PublicKey) *InitializeTickArray {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(systemProgram)
	return inst
}

// GetSystemProgramAccount gets the "systemProgram" account.
func (inst *InitializeTickArray) GetSystemProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

func (inst InitializeTickArray) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_InitializeTickArray,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst InitializeTickArray) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *InitializeTickArray) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.StartTickIndex == nil {
			return errors.New("StartTickIndex parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Whirlpool is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Funder is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.TickArray is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.SystemProgram is not set")
		}
	}
	return nil
}

func (inst *InitializeTickArray) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("InitializeTickArray")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=1]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("StartTickIndex", *inst.StartTickIndex))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=4]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("    whirlpool", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("       funder", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("    tickArray", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("systemProgram", inst.AccountMetaSlice.Get(3)))
					})
				})
		})
}

func (obj InitializeTickArray) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `StartTickIndex` param:
	err = encoder.Encode(obj.StartTickIndex)
	if err != nil {
		return err
	}
	return nil
}
func (obj *InitializeTickArray) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `StartTickIndex`:
	err = decoder.Decode(&obj.StartTickIndex)
	if err != nil {
		return err
	}
	return nil
}

// NewInitializeTickArrayInstruction declares a new InitializeTickArray instruction with the provided parameters and accounts.
func NewInitializeTickArrayInstruction(
	// Parameters:
	startTickIndex int32,
	// Accounts:
	whirlpool ag_solanago.PublicKey,
	funder ag_solanago.PublicKey,
	tickArray ag_solanago.PublicKey,
	systemProgram ag_solanago.PublicKey) *InitializeTickArray {
	return NewInitializeTickArrayInstructionBuilder().
		SetStartTickIndex(startTickIndex).
		SetWhirlpoolAccount(whirlpool).
		SetFunderAccount(funder).
		SetTickArrayAccount(tickArray).
		SetSystemProgramAccount(systemProgram)
}
