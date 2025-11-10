package registry

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	crypto2 "github.com/onflow/crypto"
	flowsdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go/fvm/systemcontracts"
	"github.com/onflow/flow-go/model/flow"
	"github.com/onflow/flow-go/module/signature"

	"github.com/onflow/flow-execution-effort-estimation/load_generator/models"
	"github.com/onflow/flow-execution-effort-estimation/load_generator/templates"
)

func simpleTemplateWithLoop(
	name string,
	label models.Label,
	initialLoopLength uint64,
	body string,
) *templates.SimpleTemplate {
	return templates.NewSimpleTemplate(
		name,
		label,
		1,
	).
		WithInitialParameters(models.Parameters{initialLoopLength}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(
					context models.Context,
					account models.Account,
				) (models.TransactionEdit, error) {
					sc := systemcontracts.SystemContractsForChain(context.ChainID)

					return models.TransactionEdit{
						Imports: map[string]flow.Address{
							sc.FlowToken.Name:     sc.FlowToken.Address,
							sc.FungibleToken.Name: sc.FungibleToken.Address,
						},
						PrepareBlock: templates.LoopTemplate(parameters[0], body),
					}, nil
				}
			},
		)
}

var simpleTemplates = []models.Template{
	simpleTemplateWithLoop(
		"empty loop",
		"EL",
		5397,
		"",
	),
	simpleTemplateWithLoop(
		"Assert True",
		"Assert",
		2657,
		"assert(true)",
	),
	simpleTemplateWithLoop(
		"get signer address",
		"GSA",
		3590,
		"signer.address",
	),
	simpleTemplateWithLoop(
		"get signer public account",
		"GSAcc",
		2627,
		"getAccount(signer.address)",
	),
	simpleTemplateWithLoop(
		"get signer account balance",
		"GSAccBal",
		27,
		"signer.balance",
	),
	simpleTemplateWithLoop(
		"get signer account available balance",
		"GSAccAwBal",
		24,
		"signer.availableBalance",
	),
	simpleTemplateWithLoop(
		"get signer account storage used",
		"GSAccSU",
		670,
		"signer.storage.used",
	),
	simpleTemplateWithLoop(
		"get signer account storage capacity",
		"GSAccSC",
		26,
		"signer.storage.capacity",
	),
	simpleTemplateWithLoop(
		"borrow signer account FlowToken.Vault",
		"BFTV",
		719,
		"let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)!",
	),
	simpleTemplateWithLoop(
		"borrow signer account FungibleToken.Receiver",
		"BFR",
		393,
		`
			let receiverRef = getAccount(signer.address)
				.capabilities.borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)!
			`,
	),
	simpleTemplateWithLoop(
		"transfer tokens to self",
		"TTS",
		26,
		`
			let vaultRef = signer.storage.borrow<auth(FungibleToken.Withdraw) &FlowToken.Vault>(from: /storage/flowTokenVault)!
			let receiverRef = getAccount(signer.address)
				.capabilities.borrow<&{FungibleToken.Receiver}>(/public/flowTokenReceiver)!
			receiverRef.deposit(from: <-vaultRef.withdraw(amount: 0.00001))
			`,
	),
	simpleTemplateWithLoop(
		"create new account",
		"CA",
		7,
		`let acct = Account(payer: signer)`,
	),
	simpleTemplateWithLoop(
		"create new account with contract",
		"CAWC",
		5,
		`
				let acct = Account(payer: signer)
				acct.contracts.add(name: "EmptyContract", code: "61636365737328616c6c2920636f6e747261637420456d707479436f6e7472616374207b7d".decodeHex())
			`,
	),
	simpleTemplateWithLoop(
		"decode hex",
		"HEX",
		811,
		`
				"f847b84000fb479cb398ab7e31d6f048c12ec5b5b679052589280cacde421af823f93fe927dfc3d1e371b172f97ceeac1bc235f60654184c83f4ea70dd3b7785ffb3c73802038203e8".decodeHex()
			`,
	),
	simpleTemplateWithLoop(
		"Revertible random",
		"RR",
		1958,
		`
				revertibleRandom<UInt64>(modulo: UInt64(100))
			`,
	),
	simpleTemplateWithLoop(
		"number to string conversion",
		"TS",
		2627,
		`
				i.toString()
			`,
	),
	simpleTemplateWithLoop(
		"concatenate string",
		"CS",
		2092,
		`
				"x".concat(i.toString())
			`,
	),
	// no unique intensities
	//templates.NewSimpleTemplate(
	//	"store string",
	//	"AStSt",
	//	1,
	//).
	//	WithInitialParameters(models.Parameters{20}).
	//	WithTransactionEdit(
	//		func(parameters models.Parameters) models.TransactionEditFunc {
	//			return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
	//				body := fmt.Sprintf(`
	//					signer.storage.save("%s", to: /storage/AStSt)
	//					signer.storage.load<String>(from: /storage/AStSt)
	//				`, stringOfLen(parameters[0]))
	//
	//				return models.TransactionEdit{
	//					PrepareBlock: templates.LoopTemplate(100, body),
	//				}, nil
	//			}
	//		},
	//	).
	//	WithAccountSetup(
	//		func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, params models.Parameters) error {
	//			// remove anything that might be in storage
	//			return models.RunTransactionBodyAsAccount(
	//				ctx,
	//				"signer.storage.load<String>(from: /storage/AStSt)",
	//				account,
	//				interaction,
	//				nil,
	//			)
	//		},
	//	),
	// no unique intensities
	//templates.NewSimpleTemplate(
	//	"create long string",
	//	"LngStr",
	//	1,
	//).
	//	WithInitialParameters(models.Parameters{100}).
	//	WithTransactionEdit(
	//		func(parameters models.Parameters) models.TransactionEditFunc {
	//			return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
	//				body := fmt.Sprintf(`
	//					var s = "%s"
	//					var S = ""
	//                    %s
	//				`, stringOfLen(parameters[0]+1), templates.LoopTemplate(100,
	//					`
	//						S = S.concat(s)
	//					`))
	//
	//				return models.TransactionEdit{
	//					PrepareBlock: body,
	//				}, nil
	//			}
	//		},
	//	),
	// no unique intensities
	//templates.NewSimpleTemplate(
	//	"create long string with templating",
	//	"LngStrTmpl",
	//	1,
	//).
	//	WithInitialParameters(models.Parameters{100}).
	//	WithTransactionEdit(
	//		func(parameters models.Parameters) models.TransactionEditFunc {
	//			return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
	//				body := fmt.Sprintf(`
	//					var s = "%s"
	//					var S = ""
	//                    %s
	//				`, stringOfLen(parameters[0]+1), templates.LoopTemplate(100,
	//					`
	//						S = "\(S)\(s)"
	//					`))
	//
	//				return models.TransactionEdit{
	//					PrepareBlock: body,
	//				}, nil
	//			}
	//		},
	//	),
	templates.NewSimpleTemplate(
		"borrow string",
		"ABrSt",
		1,
	).
		WithInitialParameters(models.Parameters{1326}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let strings = signer.storage.borrow<&[String]>(from: /storage/ABrSt)!
						var i = 0
						var lenSum = 0
						while (i < strings.length) {
							lenSum = lenSum + strings[i].length
							i = i + 1
						}
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<[String]>(from: /storage/ABrSt)
						let strings: [String] = %s
						signer.storage.save<[String]>(strings, to: /storage/ABrSt)
						`, stringArrayOfLen(20, parameters[0]))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"copy string",
		"ACpSt",
		1,
	).
		WithInitialParameters(models.Parameters{1352}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let strings = signer.storage.copy<[String]>(from: /storage/ACpSt)!
						var i = 0
						var lenSum = 0
						while (i < strings.length) {
							lenSum = lenSum + strings[i].length
							i = i + 1
						}
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<[String]>(from: /storage/ACpSt)
						let strings: [String] = %s
						signer.storage.save<[String]>(strings, to: /storage/ACpSt)
						`, stringArrayOfLen(20, parameters[0]))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"copy string and save a duplicate",
		"ACpStSv",
		1,
	).
		WithInitialParameters(models.Parameters{1223}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let strings = signer.storage.copy<[String]>(from: /storage/ACpStSv)!
						var i = 0
						var lenSum = 0
						while (i < strings.length) {
							lenSum = lenSum + strings[i].length
							i = i + 1
						}
						signer.storage.save(strings, to: /storage/ACpStSv2)
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<[String]>(from: /storage/ACpStSv)
						signer.storage.load<[String]>(from: /storage/ACpStSv2)
						let strings: [String] = %s
						signer.storage.save<[String]>(strings, to: /storage/ACpStSv)
						`, stringArrayOfLen(20, parameters[0]))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"store and load dict string",
		"AStDSt",
		1,
	).
		WithInitialParameters(models.Parameters{786}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						signer.storage.save<{String: String}>(%s, to: /storage/AStDSt)
						signer.storage.load<{String: String}>(from: /storage/AStDSt)
					`, stringDictOfLen(parameters[0], 75))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, params models.Parameters) error {
				// remove anything that might be in storage
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					"signer.storage.load<{String: String}>(from: /storage/AStDSt)",
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"store load and destroy dict string",
		"ALdDStD",
		1,
	).
		WithInitialParameters(models.Parameters{3324}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := `
						let strings = signer.storage.load<{String: String}>(from: /storage/ALdDStD)!
						for key in strings.keys {
							strings.remove(key: key)
						}
					`

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<{String: String}>(from: /storage/ALdDStD)
						let strings: {String: String} = %s
						signer.storage.save<{String: String}>(strings, to: /storage/ALdDStD)
						`, stringDictOfLen(100, parameters[0]))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"borrow dict string",
		"ABrDSt",
		1,
	).
		WithInitialParameters(models.Parameters{206}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let strings = signer.storage.borrow<&{String: String}>(from: /storage/ABrDSt)!
						var lenSum = 0
						strings.forEachKey(fun (key: String): Bool {
							lenSum = lenSum + strings[key]!.length
							return true
						})
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<{String: String}>(from: /storage/ABrDSt)
						let strings: {String: String} = %s
						signer.storage.save<{String: String}>(strings, to: /storage/ABrDSt)
						`, stringDictOfLen(parameters[0], 100))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"copy dict string",
		"ACpDSt",
		1,
	).
		WithInitialParameters(models.Parameters{813}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let strings = signer.storage.copy<{String: String}>(from: /storage/ACpDSt)!
						var lenSum = 0
						strings.forEachKey(fun (key: String): Bool {
							lenSum = lenSum + strings[key]!.length
							return true
						})
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<{String: String}>(from: /storage/ACpDSt)
						let strings: {String: String} = %s
						signer.storage.save<{String: String}>(strings, to: /storage/ACpDSt)
						`, stringDictOfLen(30, parameters[0]))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"copy string dict and save a duplicate",
		"ACpDStSv",
		1,
	).
		WithInitialParameters(models.Parameters{1179}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let strings = signer.storage.copy<{String: String}>(from: /storage/ACpDStSv)!
						var lenSum = 0
						strings.forEachKey(fun (key: String): Bool {
							lenSum = lenSum + strings[key]!.length
							return true
						})
						signer.storage.save(strings, to: /storage/ACpDStSv2)
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(`
						signer.storage.load<{String: String}>(from: /storage/ACpDStSv)
						signer.storage.load<{String: String}>(from: /storage/ACpDStSv2)
						let strings: {String: String} = %s
						signer.storage.save(strings, to: /storage/ACpDStSv)
						`, stringDictOfLen(20, parameters[0]))),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	templates.NewSimpleTemplate(
		"load dict and destroy it",
		"DestDict",
		1,
	).
		WithInitialParameters(models.Parameters{967}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					return models.TransactionEdit{
						PrepareBlock: `
						let r <- signer.storage.load<@{String: AnyResource}>(from: /storage/DestDict)!
						destroy r
						`,
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, parameters models.Parameters) error {
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					models.TransactionBody(fmt.Sprintf(
						`
						let r <- signer.storage.load<@{String: AnyResource}>(from: /storage/DestDict)
                        destroy r
						let r2: @{String: AnyResource} <- {}
						var i = 0
						while (i < %d) {
							i = i + 1
							let d: @{String: AnyResource} <- {}
							r2[i.toString()] <-! d
						}

						signer.storage.save<@{String: AnyResource}>( <- r2, to: /storage/DestDict)
						`, parameters[0])),
					account,
					interaction,
					nil,
				)
				return err
			},
		),
	simpleTemplateWithLoop(
		"add key to account",
		"KA",
		49,
		`
				let key = PublicKey(
					publicKey: "f7901e9161b9b53f2e1f27b0f1e4711fcc8f234a90f55fd2068a67b152948389c0ee1e40f74a0e194ef7c2b59666270b16d52cf585fd8e65fc00958f78af77b0".decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1
				)
		
				signer.keys.add(
					publicKey: key,
					hashAlgorithm: HashAlgorithm.SHA3_256,
					weight: 0.0
				)
			`,
	),
	simpleTemplateWithLoop(
		"add and revoke key to account",
		"KAR",
		38,
		`
				let key = PublicKey(
					publicKey: "f7901e9161b9b53f2e1f27b0f1e4711fcc8f234a90f55fd2068a67b152948389c0ee1e40f74a0e194ef7c2b59666270b16d52cf585fd8e65fc00958f78af77b0".decodeHex(),
					signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1
				)
		
				let ac = signer.keys.add(
					publicKey: key,
					hashAlgorithm: HashAlgorithm.SHA3_256,
					weight: 0.0
				)
				signer.keys.revoke(keyIndex: ac.keyIndex)
			`,
	),
	simpleTemplateWithLoop(
		"get account key",
		"KGet",
		177,
		`
				let key = signer.keys.get(keyIndex: 0)
			`,
	),
	templates.NewSimpleTemplate(
		"Get contracts",
		"GetCon",
		1,
	).
		WithInitialParameters(models.Parameters{65}).
		WithTransactionEdit(
			func(params models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := `
						signer.contracts.names
					`

					return models.TransactionEdit{
						PrepareBlock: templates.LoopTemplate(params[0], body),
					}, nil
				}
			},
		).
		WithAccountSetup(
			func(ctx context.Context, c models.Context, interaction models.ChainInteraction, account models.Account, params models.Parameters) error {
				// deploy contracts
				_, err := models.RunTransactionBodyAsAccount(
					ctx,
					`
						var c = signer.contracts.names.length
						while c < 20 {
							// deploy contract
							let contractName = "TestContract".concat(c.toString())
							let contractCode = "access(all) contract ".concat(contractName).concat(" {}")
							signer.contracts.add(name: contractName, code: contractCode.utf8)
							c = c + 1
						}`,
					account,
					interaction,
					nil,
				)
				return err
			},
		),

	templates.NewSimpleTemplate(
		"hash",
		"H",
		1,
	).
		WithInitialParameters(models.Parameters{328}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						Crypto.hash("%s".utf8, algorithm: HashAlgorithm.SHA2_256)
					`, stringOfLen(20))

					return models.TransactionEdit{
						Imports:      map[string]flow.Address{"Crypto": flow.EmptyAddress},
						PrepareBlock: templates.LoopTemplate(parameters[0], body),
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"String toLower",
		"STL",
		1,
	).
		WithInitialParameters(models.Parameters{5221}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						var s = "%s"
						s = s.toLower()
					`, stringOfLen(parameters[0]))

					return models.TransactionEdit{
						PrepareBlock: templates.LoopTemplate(500, body),
					}, nil
				}
			},
		),

	simpleTemplateWithLoop(
		"Get Current Block",
		"GCB",
		1107,
		`
				getCurrentBlock()
			`,
	),
	simpleTemplateWithLoop(
		"Get Block At",
		"GBA",
		548,
		`
				let at = getCurrentBlock().height
				getBlock(at: at)
			`,
	),
	simpleTemplateWithLoop(
		"Destroy Resource Dictionary",
		"DRD",
		1058,
		`
				let r: @{String: AnyResource} <- {}
				destroy r
			`,
	),
	simpleTemplateWithLoop(
		"Parse UFix64",
		"PUFix",
		2032,
		`
				let smol: UFix64? = UFix64.fromString("0.123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse Fix64",
		"PFix",
		2047,
		`
				let smol: Fix64? = Fix64.fromString("-0.123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse UInt64",
		"PUInt64",
		2578,
		`
				let smol: UInt64? = UInt64.fromString("123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse Int64",
		"PInt64",
		2417,
		`
				let smol: Int64? = Int64.fromString("-123456")
			`,
	),
	simpleTemplateWithLoop(
		"Parse Int",
		"PInt",
		2460,
		`
				let smol: Int? = Int.fromString("-12345")
			`,
	),
	simpleTemplateWithLoop(
		"Issue storage capability",
		"ISCap",
		182,
		`
				let cap = signer.capabilities.storage.issue<&Int>(/storage/foo)
			`,
	),
	simpleTemplateWithLoop(
		"Get key count",
		"GKC",
		677,
		`
				let count = signer.keys.count
			`,
	),
	templates.NewSimpleTemplate(
		"Create Key ECDSA_P256",
		"CrKeyP256",
		1,
	).
		WithInitialParameters(models.Parameters{112}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					seed := make([]byte, crypto.MinSeedLength)
					for i := range seed {
						seed[i] = 0
					}

					privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
					}
					key := hex.EncodeToString(privateKey.PublicKey().Encode())

					body := fmt.Sprintf(`
						  let publicKey = PublicKey(
							publicKey: "%s".decodeHex(),
							signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
						)
					`, key)

					return models.TransactionEdit{
						PrepareBlock: templates.LoopTemplate(parameters[0], body),
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Create Key ECDSA_secp256k1",
		"CrKeysecp256k1",
		1,
	).
		WithInitialParameters(models.Parameters{112}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					seed := make([]byte, crypto.MinSeedLength)
					for i := range seed {
						seed[i] = 0
					}

					privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_secp256k1, seed)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
					}
					key := hex.EncodeToString(privateKey.PublicKey().Encode())

					body := fmt.Sprintf(`
						  let publicKey = PublicKey(
							publicKey: "%s".decodeHex(),
							signatureAlgorithm: SignatureAlgorithm.ECDSA_secp256k1
						)
					`, key)

					return models.TransactionEdit{
						PrepareBlock: templates.LoopTemplate(parameters[0], body),
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Create Key BLS_BLS12_381",
		"CrKeyBLS",
		1,
	).
		WithInitialParameters(models.Parameters{81}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					seed := make([]byte, crypto.MinSeedLength)
					for i := range seed {
						seed[i] = 0
					}

					privateKey, err := crypto.GeneratePrivateKey(crypto.BLS_BLS12_381, seed)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
					}
					key := hex.EncodeToString(privateKey.PublicKey().Encode())

					body := fmt.Sprintf(`
						  let publicKey = PublicKey(
							publicKey: "%s".decodeHex(),
							signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
						)
					`, key)

					return models.TransactionEdit{
						PrepareBlock: templates.LoopTemplate(parameters[0], body),
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Insert",
		"ArrIns",
		1,
	).
		WithInitialParameters(models.Parameters{1577}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = [0]
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`x.insert(at: i, 1)`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Insert Remove",
		"ArrInsDel",
		1,
	).
		WithInitialParameters(models.Parameters{1232}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = [0]
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(at: 0, 1)
							x.remove(at: 1)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Insert Set Remove",
		"ArrInsSetDel",
		1,
	).
		WithInitialParameters(models.Parameters{994}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = [0]
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(at: 0, 1)
							x[0] = i
							x.remove(at: 1)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Insert Map",
		"ArrInsMap",
		1,
	).
		WithInitialParameters(models.Parameters{1345}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = [0]
                           %s
							let addOne =
								fun (_ v: Int): Int {
									return v+1
								}
							let y = x.map(addOne)
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(at: 0, i)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Insert Filter",
		"ArrInsFilt",
		1,
	).
		WithInitialParameters(models.Parameters{1356}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = [0]
                           %s
							let isEven =
								view fun (element: Int): Bool {
									return element %% 2 == 0
								}
							let y = x.filter(isEven)
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(at: 0, i)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Append",
		"ArrApp",
		1,
	).
		WithInitialParameters(models.Parameters{1905}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = [0]
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.append(i)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Dict Insert",
		"DictIns",
		1,
	).
		WithInitialParameters(models.Parameters{1598}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = {"0": 0}
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(key: i.toString(), i)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Dict Insert Remove",
		"DictInsDel",
		1,
	).
		WithInitialParameters(models.Parameters{713}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = {"0": 0}
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(key: i.toString(), i)
							x.remove(key: (i-1).toString())
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Dict Insert Set Remove",
		"DictInsSetDel",
		1,
	).
		WithInitialParameters(models.Parameters{565}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						   let x = {"0": 0}
                           %s
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(key: i.toString(), i)
							x[(i-1).toString()] = i
							x.remove(key: (i-1).toString())
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Dict Iter Copy",
		"DictItrCpy",
		1,
	).
		WithInitialParameters(models.Parameters{667}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					body := fmt.Sprintf(`
						let x = {"0": 0}
						let y = {"0": 0}
						%s
						x.forEachKey(fun (key: String): Bool {
							y[key] = x[key]
							return true
						})
					`, templates.LoopTemplate(
						parameters[0],
						`
							x.insert(key: i.toString(), i)
							`,
					))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Array Create Batch",
		"ArrCB",
		1,
	).
		WithInitialParameters(models.Parameters{226}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					sumStr := "0"
					for i := 0; i < int(parameters[0]); i++ {
						sumStr += fmt.Sprintf(",%d", i)
					}

					body := fmt.Sprintf(`
						var i = 0
						while i < 200 {
							i = i + 1
							let a = [%s]
						}
					`, sumStr)

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),

	templates.NewSimpleTemplate(
		"Verify Signature",
		"VerSig",
		1,
	).
		WithInitialParameters(models.Parameters{14}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					numKeys := parameters[0] + 1
					message := []byte("hello world")

					rawKeys := make([]string, numKeys)
					signers := make([]crypto.Signer, numKeys)
					signatures := make([]string, numKeys)

					for i := 0; i < int(numKeys); i++ {
						seed := make([]byte, crypto.MinSeedLength)
						_, err := rand.Read(seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate seed: %w", err)
						}

						privateKey, err := crypto.GeneratePrivateKey(crypto.ECDSA_P256, seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
						}
						rawKeys[i] = hex.EncodeToString(privateKey.PublicKey().Encode())
						sig, err := crypto.NewInMemorySigner(privateKey, crypto.SHA3_256)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate signer: %w", err)
						}
						signers[i] = sig
					}

					keyListAdd := ""
					for i := 0; i < int(numKeys); i++ {
						keyListAdd += fmt.Sprintf(`
								  keyList.add(
									PublicKey(
									  publicKey: "%s".decodeHex(),
									  signatureAlgorithm: SignatureAlgorithm.ECDSA_P256
									),
									hashAlgorithm: HashAlgorithm.SHA3_256,
									weight: 1.0/%d.0+0.000001 ,
								  )
								`, rawKeys[i], int(numKeys),
						)
					}

					for i := 0; i < int(numKeys); i++ {
						sig, err := flowsdk.SignUserMessage(signers[i], message)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to sign message: %w", err)
						}
						signatures[i] = hex.EncodeToString(sig)
					}

					signaturesAdd := ""
					for i := 0; i < int(numKeys); i++ {
						signaturesAdd += fmt.Sprintf(`
								   signatureSet.append(
									Crypto.KeyListSignature(
									  keyIndex: %d,
									  signature: "%s".decodeHex()
									)
								  )
								`, i, signatures[i],
						)
					}

					body := fmt.Sprintf(`
								let keyList = Crypto.KeyList()

								%s

								let signatureSet: [Crypto.KeyListSignature] = []
				
								%s
				
								let domainSeparationTag = "FLOW-V0.0-user"
								let message = "%s".decodeHex()
								
								
								let valid = keyList.verify(
								  signatureSet: signatureSet,
								  signedData: message,
								  domainSeparationTag: domainSeparationTag
								)
								if !valid {
									panic("invalid signature")
								}
								
							  
					`, keyListAdd, signaturesAdd, hex.EncodeToString(message))

					return models.TransactionEdit{
						PrepareBlock: body,
						Imports: map[string]flow.Address{
							"Crypto": {},
						},
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Aggregate BLS aggregate signature",
		"BLSAggSig",
		1,
	).
		WithInitialParameters(models.Parameters{186}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					numSigs := int(parameters[0]) + 1
					sigs := make([]crypto2.Signature, 0, numSigs)
					kmac := signature.NewBLSHasher("test tag")
					signatureAlgorithm := crypto2.BLSBLS12381
					input := make([]byte, 100)
					_, err := rand.Read(input)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate random data to sign: %w", err)
					}

					for i := 0; i < numSigs; i++ {
						seed := make([]byte, crypto2.KeyGenSeedMinLen)
						_, err := rand.Read(seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate seed: %w", err)
						}
						sk, err := crypto.GeneratePrivateKey(signatureAlgorithm, seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
						}
						// a valid BLS signature
						s, err := sk.Sign(input, kmac)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to sign message: %w", err)
						}
						sigs = append(sigs, s)
					}

					signatures := ""
					for i := 0; i < numSigs; i++ {
						signatures += fmt.Sprintf(`
								signatures.append("%s".decodeHex())
							`, hex.EncodeToString(sigs[i].Bytes()))
					}

					body := fmt.Sprintf(`
						var signatures: [[UInt8]] = []
                        %s
						BLS.aggregateSignatures(signatures)!
					`, signatures)

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"Aggregate BLS aggregate keys",
		"BLSAggKey",
		1,
	).
		WithInitialParameters(models.Parameters{41}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					numSigs := int(parameters[0]) + 1
					pks := make([]crypto2.PublicKey, 0, numSigs)
					signatureAlgorithm := crypto2.BLSBLS12381
					input := make([]byte, 100)
					_, err := rand.Read(input)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate random data to sign: %w", err)
					}

					for i := 0; i < numSigs; i++ {
						seed := make([]byte, crypto2.KeyGenSeedMinLen)
						_, err := rand.Read(seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate seed: %w", err)
						}
						sk, err := crypto.GeneratePrivateKey(signatureAlgorithm, seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
						}

						pks = append(pks, sk.PublicKey())
					}

					pkString := ""
					for i := 0; i < numSigs; i++ {
						pkString += fmt.Sprintf(`
								pks.append(PublicKey(
									publicKey: "%s".decodeHex(), 
									signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
								))
							`, hex.EncodeToString(pks[i].Encode()))
					}

					body := fmt.Sprintf(`
						let pks: [PublicKey] = []
                        %s
						BLS.aggregatePublicKeys(pks)!.publicKey
					`, pkString)

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"BLS verify signature",
		"BLSVer",
		1,
	).
		WithInitialParameters(models.Parameters{32}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					numSigs := int(parameters[0]) + 1
					pks := make([]crypto2.PublicKey, 0, numSigs)
					signatures := make([]crypto2.Signature, 0, numSigs)
					signatureAlgorithm := crypto2.BLSBLS12381
					input := make([]byte, 100)
					_, err := rand.Read(input)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate random data to sign: %w", err)
					}

					message := []byte("random_message")
					tag := "random_tag"
					kmac := signature.NewBLSHasher(tag)

					for i := 0; i < numSigs; i++ {
						seed := make([]byte, crypto2.KeyGenSeedMinLen)
						_, err := rand.Read(seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate seed: %w", err)
						}
						sk, err := crypto.GeneratePrivateKey(signatureAlgorithm, seed)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
						}
						pks = append(pks, sk.PublicKey())

						sig, err := sk.Sign(message, kmac)
						if err != nil {
							return models.TransactionEdit{}, fmt.Errorf("failed to sign message: %w", err)
						}
						signatures = append(signatures, sig)
					}

					signaturesString := ""
					for i := 0; i < numSigs; i++ {
						signaturesString += fmt.Sprintf(`
								signatures.append("%s".decodeHex())
							`, hex.EncodeToString(signatures[i].Bytes()))
					}

					pkString := ""
					for i := 0; i < numSigs; i++ {
						pkString += fmt.Sprintf(`
								pks.append(PublicKey(
									publicKey: "%s".decodeHex(), 
									signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
								))
							`, hex.EncodeToString(pks[i].Encode()))
					}

					body := fmt.Sprintf(`
						var pks: [PublicKey] = []
                        var signatures: [[UInt8]] = []
                        %s
						%s
						let aggPk = BLS.aggregatePublicKeys(pks)!
                        let aggSignature = BLS.aggregateSignatures(signatures)!
						let boo = aggPk.verify(
									signature: aggSignature, 
									signedData: "%s".decodeHex(),
									domainSeparationTag: "random_tag", 
									hashAlgorithm: HashAlgorithm.KMAC128_BLS_BLS12_381)
						if !boo {
							panic("invalid signature")
						}
					`, pkString, signaturesString, hex.EncodeToString(message))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
	templates.NewSimpleTemplate(
		"BLS verify proof of possession",
		"BLSVerPoP",
		1,
	).
		WithInitialParameters(models.Parameters{8}).
		WithTransactionEdit(
			func(parameters models.Parameters) models.TransactionEditFunc {
				return func(c models.Context, a models.Account) (models.TransactionEdit, error) {
					signatureAlgorithm := crypto2.BLSBLS12381
					seed := make([]byte, crypto2.KeyGenSeedMinLen)
					_, err := rand.Read(seed)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate seed: %w", err)
					}
					sk, err := crypto.GeneratePrivateKey(signatureAlgorithm, seed)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate private key: %w", err)
					}
					pk := sk.PublicKey()

					proof, err := crypto2.BLSGeneratePOP(sk)
					if err != nil {
						return models.TransactionEdit{}, fmt.Errorf("failed to generate proof of possession: %w", err)
					}

					body := fmt.Sprintf(`
							let p = PublicKey(
								publicKey: "%s".decodeHex(), 
								signatureAlgorithm: SignatureAlgorithm.BLS_BLS12_381
							)
							var proof = "%s".decodeHex()

							%s
						`,
						hex.EncodeToString(pk.Encode()),
						hex.EncodeToString(proof.Bytes()),
						templates.LoopTemplate(parameters[0], `
							var valid = p.verifyPoP(proof)
							if !valid {
								panic("invalid proof of possession")
							}
						`))

					return models.TransactionEdit{
						PrepareBlock: body,
					}, nil
				}
			},
		),
}

func stringOfLen(length uint64) string {
	someString := make([]byte, length)
	for i := 0; i < len(someString); i++ {
		someString[i] = 'x'
	}
	return string(someString)
}

func stringArrayOfLen(arrayLen uint64, stringLen uint64) string {
	builder := strings.Builder{}
	builder.WriteRune('[')
	for i := uint64(0); i < arrayLen; i++ {
		if i > 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune('"')
		builder.WriteString(stringOfLen(stringLen))
		builder.WriteRune('"')
	}
	builder.WriteRune(']')
	return builder.String()
}

func stringDictOfLen(dictLen uint64, stringLen uint64) string {
	builder := strings.Builder{}
	builder.WriteRune('{')
	for i := uint64(0); i < dictLen; i++ {
		if i > 0 {
			builder.WriteRune(',')
		}
		builder.WriteRune('"')
		someString := make([]byte, stringLen)
		for i := 0; i < len(someString); i++ {
			someString[i] = 'x'
		}
		builder.WriteString(string(someString))
		builder.WriteString(strconv.Itoa(int(i)))
		builder.WriteRune('"')
		builder.WriteRune(':')
		builder.WriteRune('"')
		builder.WriteString(string(someString))
		builder.WriteRune('"')
	}
	builder.WriteRune('}')
	return builder.String()
}
